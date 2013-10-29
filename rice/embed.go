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
	// create one list
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
	// loop over files
	for _, filename := range filenames {
		fullpath := filepath.Join(pkg.Dir, filename)
		if flags.Verbose {
			fmt.Printf("scanning file (%s)\n", fullpath)
		}
		file, err := os.Open(fullpath)
		if err != nil {
			fmt.Printf("error opening file '%s': %s\n", filename, err)
			os.Exit(-1)
		}
		defer file.Close()
		fileData, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Printf("error reading file '%s': %s\n", filename, err)
			os.Exit(-1)
		}
		matches := regexpBox.FindAllStringSubmatch(string(fileData), -1)
		for _, match := range matches {
			boxMap[match[1]] = true
			if flags.Verbose {
				fmt.Printf("\tfound box (%s)\n", match[1])
			}
		}
	}

	if len(boxMap) == 0 {
		fmt.Println("no calls to rice.FindBox() found")
	}

	if flags.Verbose {
		fmt.Println("")
	}

	for boxname := range boxMap {
		boxPath := filepath.Join(pkg.Dir, boxname)
		boxFilename := boxname + `.rice-box.go`
		if flags.Verbose {
			fmt.Printf("embedding box (%s)\n", boxPath)
			fmt.Printf("\tto file (%s)\n", boxFilename)
		}

		box := &boxDataType{
			Package: pkg.Name,
			BoxName: boxname,
			UnixNow: time.Now().Unix(),
			Files:   make([]*fileDataType, 0),
		}

		filepath.Walk(boxPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("error walking box: %s\n", err)
				os.Exit(-1)
			}
			if !info.IsDir() {
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

		boxFile, err := os.Create(boxFilename)
		if err != nil {
			fmt.Printf("error creating embedded box file: %s\n", err)
			os.Exit(-1)
		}
		defer boxFile.Close()
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
