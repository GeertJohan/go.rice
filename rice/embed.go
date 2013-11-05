package main

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func operationEmbed(pkg *build.Package) {
	// create one list of files for this package
	filenames := make([]string, 0, len(pkg.GoFiles)+len(pkg.CgoFiles))
	filenames = append(filenames, pkg.GoFiles...)
	filenames = append(filenames, pkg.CgoFiles...)

	// prepare regex to find calls to rice.FindBox(..)
	regexpBox, err := regexp.Compile(`rice\.FindBox\(["` + "`" + `]{1}([a-zA-Z0-9\\/\.-]+)["` + "`" + `]{1}\)`)
	if err != nil {
		fmt.Printf("error compiling rice.FindBox regexp: %s\n", err)
		os.Exit(-1)
	}

	// create map of boxes to embed
	var boxMap = make(map[string]bool)

	// loop over files, search for rice.FindBox(..) calls
	for _, filename := range filenames {
		// find full filepath
		fullpath := filepath.Join(pkg.Dir, filename)
		if flags.Verbose {
			fmt.Printf("scanning file (%s)\n", fullpath)
		}

		// open source file
		file, err := os.Open(fullpath)
		if err != nil {
			fmt.Printf("error opening file '%s': %s\n", filename, err)
			os.Exit(-1)
		}
		defer file.Close()

		// slurp source code
		fileData, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Printf("error reading file '%s': %s\n", filename, err)
			os.Exit(-1)
		}

		// find rice.FindBox(..) calls
		matches := regexpBox.FindAllStringSubmatch(string(fileData), -1)
		for _, match := range matches {
			boxMap[match[1]] = true
			if flags.Verbose {
				fmt.Printf("\tfound box (%s)\n", match[1])
			}
		}
	}

	// notify user when no calls to rice.FindBox are made (is this an error and therefore os.Exit(-1) ?
	if len(boxMap) == 0 {
		fmt.Println("no calls to rice.FindBox() found")
	}

	if flags.Verbose {
		fmt.Println("")
	}

	for boxname := range boxMap {
		// find path and filename for this box
		boxPath := filepath.Join(pkg.Dir, boxname)
		boxFilename := boxname + `.rice-box.go`

		// verbose info
		if flags.Verbose {
			fmt.Printf("embedding box (%s)\n", boxPath)
			fmt.Printf("\tto file (%s)\n", boxFilename)
		}

		// create box datastructure (used by template)
		box := &boxDataType{
			Package: pkg.Name,
			BoxName: boxname,
			UnixNow: time.Now().Unix(),
			Files:   make([]*fileDataType, 0),
		}

		// fill box datastructure with file data
		filepath.Walk(boxPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("error walking box: %s\n", err)
				os.Exit(-1)
			}

			//++ add references and identifiers (a, b, c, d, e, f, g, h, i, j, k, etc.)

			if info.IsDir() {
				//++ dirDataType and stuff
			} else {
				fileData := &fileDataType{
					FileName: strings.TrimPrefix(strings.TrimPrefix(path, boxPath), "/"),
					ModTime:  info.ModTime().Unix(),
				}
				fileData.Content, err = ioutil.ReadFile(path)
				if err != nil {
					fmt.Printf("error reading file content while walking box: %s\n", err)
					os.Exit(-1)
				}
				box.Files = append(box.Files, fileData)
			}
			return nil
		})

		// create go file for box
		boxFile, err := os.Create(boxFilename)
		if err != nil {
			fmt.Printf("error creating embedded box file: %s\n", err)
			os.Exit(-1)
		}
		defer boxFile.Close()

		// execute template (write result directly to file)
		err = tmplEmbeddedBox.Execute(boxFile, box)
		if err != nil {
			fmt.Printf("error writing embedded box to file (template execute): %s\n", err)
			os.Exit(-1)
		}
	}

	if flags.Verbose {
		fmt.Println("")
	}
}
