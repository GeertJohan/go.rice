package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/GeertJohan/go.rice/embedded"
	"github.com/akavel/rsrc/coff"
	"go/build"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type sizedBytes []byte

func (s sizedBytes) Size() int64 {
	return int64(len(s))
}

func operationEmbedSyso(pkg *build.Package) {

	regexpSynameReplacer := regexp.MustCompile(`[^a-z0-9_]`)

	boxMap := findBoxes(pkg)

	// notify user when no calls to rice.FindBox are made (is this an error and therefore os.Exit(1) ?
	if len(boxMap) == 0 {
		fmt.Println("no calls to rice.FindBox() found")
		return
	}

	verbosef("\n")

	for boxname := range boxMap {
		// find path and filename for this box
		boxPath := filepath.Join(pkg.Dir, boxname)
		boxFilename := strings.Replace(boxname, "/", "-", -1)
		boxFilename = strings.Replace(boxFilename, "..", "back", -1)
		boxFilename = boxFilename + `.rice-box` // append with .go and .syso

		// verbose info
		verbosef("embedding box '%s'\n", boxname)
		verbosef("\tto file %s\n", boxFilename)

		// create box datastructure (used by template)
		box := &embedded.EmbeddedBox{
			Name:      boxname,
			Time:      time.Now(),
			EmbedType: embedded.EmbedTypeSyso,
			Files:     make(map[string]*embedded.EmbeddedFile),
			Dirs:      make(map[string]*embedded.EmbeddedDir),
		}

		// fill box datastructure with file data
		filepath.Walk(boxPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("error walking box: %s\n", err)
				os.Exit(1)
			}

			filename := strings.TrimPrefix(path, boxPath)
			filename = strings.Replace(filename, "\\", "/", -1)
			filename = strings.TrimPrefix(filename, "/")
			if info.IsDir() {
				embeddedDir := &embedded.EmbeddedDir{
					Filename:   filename,
					DirModTime: info.ModTime(),
				}
				verbosef("\tincludes dir: '%s'\n", embeddedDir.Filename)
				box.Dirs[embeddedDir.Filename] = embeddedDir

				// add tree entry (skip for root, it'll create a recursion)
				if embeddedDir.Filename != "" {
					pathParts := strings.Split(embeddedDir.Filename, "/")
					parentDir := box.Dirs[strings.Join(pathParts[:len(pathParts)-1], "/")]
					parentDir.ChildDirs = append(parentDir.ChildDirs, embeddedDir)
				}
			} else {
				embeddedFile := &embedded.EmbeddedFile{
					Filename:    filename,
					FileModTime: info.ModTime(),
					Content:     "",
				}
				verbosef("\tincludes file: '%s'\n", embeddedFile.Filename)
				contentBytes, err := ioutil.ReadFile(path)
				if err != nil {
					fmt.Printf("error reading file content while walking box: %s\n", err)
					os.Exit(1)
				}
				embeddedFile.Content = string(contentBytes)
				box.Files[embeddedFile.Filename] = embeddedFile
			}
			return nil
		})

		// encode embedded box to gob file
		boxGobBuf := &bytes.Buffer{}
		err := gob.NewEncoder(boxGobBuf).Encode(box)
		if err != nil {
			fmt.Printf("error encoding box to gob: %v\n", err)
			os.Exit(1)
		}

		// write coff
		symname := regexpSynameReplacer.ReplaceAllString(boxname, "_")
		boxCoff := coff.NewRDATA()
		boxCoff.AddData("_bricebox_"+symname, sizedBytes(boxGobBuf.Bytes()))
		boxCoff.AddData("_ericebox_"+symname, io.NewSectionReader(strings.NewReader("\000\000"), 0, 2)) // TODO: why? copied from as-generated
		boxCoff.Freeze()
		err = writeCoff(boxCoff, boxFilename+"_386.syso")
		if err != nil {
			fmt.Printf("error writing coff/.syso: %v\n", err)
			os.Exit(1)
		}
	}
}
