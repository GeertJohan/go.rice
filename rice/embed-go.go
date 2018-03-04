package main

import (
	"bytes"
	"fmt"
	"go/build"
	"go/format"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const boxFilename = "rice-box.go"

func writeBoxesGo(pkg *build.Package, out io.Writer, boxMap map[string]bool) error {
	verbosef("\n")

	var boxes []*boxDataType

	for boxname := range boxMap {
		// find path and filename for this box
		boxPath := filepath.Join(pkg.Dir, boxname)

		// Check to see if the path for the box is a symbolic link.  If so, simply
		// box what the symbolic link points to.  Note: the filepath.Walk function
		// will NOT follow any nested symbolic links.  This only handles the case
		// where the root of the box is a symbolic link.
		symPath, serr := os.Readlink(boxPath)
		if serr == nil {
			boxPath = symPath
		}

		// verbose info
		verbosef("embedding box '%s' to '%s'\n", boxname, boxFilename)

		// read box metadata
		boxInfo, ierr := os.Stat(boxPath)
		if ierr != nil {
			return fmt.Errorf("Error: unable to access box at %s\n", boxPath)
		}

		// create box datastructure (used by template)
		box := &boxDataType{
			BoxName: boxname,
			UnixNow: boxInfo.ModTime().Unix(),
			Files:   make([]*fileDataType, 0),
			Dirs:    make(map[string]*dirDataType),
		}

		if !boxInfo.IsDir() {
			return fmt.Errorf("Error: Box %s must point to a directory but points to %s instead\n",
				boxname, boxPath)
		}

		// fill box datastructure with file data
		err := filepath.Walk(boxPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("error walking box: %s\n", err)
			}

			filename := strings.TrimPrefix(path, boxPath)
			filename = strings.Replace(filename, "\\", "/", -1)
			filename = strings.TrimPrefix(filename, "/")
			if info.IsDir() {
				dirData := &dirDataType{
					Identifier: "dir" + nextIdentifier(),
					FileName:   filename,
					ModTime:    info.ModTime().Unix(),
					ChildFiles: make([]*fileDataType, 0),
					ChildDirs:  make([]*dirDataType, 0),
				}
				verbosef("\tincludes dir: '%s'\n", dirData.FileName)
				box.Dirs[dirData.FileName] = dirData

				// add tree entry (skip for root, it'll create a recursion)
				if dirData.FileName != "" {
					pathParts := strings.Split(dirData.FileName, "/")
					parentDir := box.Dirs[strings.Join(pathParts[:len(pathParts)-1], "/")]
					parentDir.ChildDirs = append(parentDir.ChildDirs, dirData)
				}
			} else {
				fileData := &fileDataType{
					Identifier: "file" + nextIdentifier(),
					FileName:   filename,
					ModTime:    info.ModTime().Unix(),
				}
				verbosef("\tincludes file: '%s'\n", fileData.FileName)
				fileData.Content, err = ioutil.ReadFile(path)
				if err != nil {
					return fmt.Errorf("error reading file content while walking box: %s\n", err)
				}
				box.Files = append(box.Files, fileData)

				// add tree entry
				pathParts := strings.Split(fileData.FileName, "/")
				parentDir := box.Dirs[strings.Join(pathParts[:len(pathParts)-1], "/")]
				if parentDir == nil {
					return fmt.Errorf("Error: parent of %s is not within the box\n", path)
				}
				parentDir.ChildFiles = append(parentDir.ChildFiles, fileData)
			}
			return nil
		})
		if err != nil {
			return err
		}
		boxes = append(boxes, box)

	}

	embedSourceUnformated := bytes.NewBuffer(make([]byte, 0))

	// execute template to buffer
	err := tmplEmbeddedBox.Execute(
		embedSourceUnformated,
		embedFileDataType{pkg.Name, boxes},
	)
	if err != nil {
		return fmt.Errorf("error writing embedded box to file (template execute): %s\n", err)
	}

	// format the source code
	embedSource, err := format.Source(embedSourceUnformated.Bytes())
	if err != nil {
		return fmt.Errorf("error formatting embedSource: %s\n", err)
	}

	// write source to file
	_, err = io.Copy(out, bytes.NewBuffer(embedSource))
	if err != nil {
		return fmt.Errorf("error writing embedSource to file: %s\n", err)
	}
	return nil
}

func operationEmbedGo(pkg *build.Package) {
	boxMap := findBoxes(pkg)
	boxFilepath := filepath.Join(pkg.Dir, boxFilename)

	// notify user when no calls to rice.FindBox are made
	// this is an error and therefore os.Exit(1)
	if len(boxMap) == 0 {
		log.Printf("no calls to rice.FindBox() found: %s\n", pkg.ImportPath)
		log.Printf("deleting embedded box file %s\n", boxFilepath)
		err := os.Remove(boxFilepath)
		if err != nil && !os.IsNotExist(err) {
			log.Printf("error deleting embedded box file: %s\n", err)
		}
		os.Exit(1)
	}

	var newContentBuffer bytes.Buffer
	err := writeBoxesGo(pkg, &newContentBuffer, boxMap)
	if err != nil {
		log.Printf("error getting embedded box file content: %s\n", err)
		os.Exit(1)
	}

	boxFile, err := os.OpenFile(boxFilepath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("error opening embedded box file %s: %s\n", boxFilepath, err)
		os.Exit(1)
	}
	defer boxFile.Close()

	currentContent, err := ioutil.ReadAll(boxFile)
	if bytes.Compare(currentContent, newContentBuffer.Bytes()) == 0 {
		log.Printf("skipping since content did not changed: %s\n", boxFilepath)
		os.Exit(0)
	}

	err = boxFile.Truncate(0)
	if err != nil {
		log.Printf("error truncating embedded box file %s: %s\n", boxFilepath, err)
		os.Exit(1)
	}

	_, err = boxFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Printf("error seeking embedded box file %s: %s\n", boxFilepath, err)
		os.Exit(1)
	}

	_, err = io.Copy(boxFile, &newContentBuffer)
	if err != nil {
		log.Printf("error writing embedSource to file: %s\n", err)
	}
}
