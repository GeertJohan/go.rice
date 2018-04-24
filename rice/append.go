package main

import (
	"archive/zip"
	"fmt"
	"go/build"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	zipexe "github.com/daaku/go.zipexe"
)

// operationAppendPkgs scans the source files, looking for referenced boxes
func operationAppendPkgs(pkgs []*build.Package) {
	// find abs path for binary file
	binfileName, err := filepath.Abs(flags.Append.Executable)
	if err != nil {
		fmt.Printf("Error finding absolute path for executable to append: %s\n", err)
		os.Exit(1)
	}

	boxMap := make(map[string]string)
	for _, pkg := range pkgs {
		// find boxes for this command
		boxes := findBoxes(pkg)

		// notify user when no calls to rice.FindBox are made (is this an error and therefore os.Exit(1) ?
		if len(boxes) == 0 {
			fmt.Printf("no calls to rice.FindBox() or rice.MustFindBox() found in import path `%s`\n", pkg.ImportPath)
			continue
		}

		verbosef("\n")

		for boxname := range boxes {
			appendedBoxName := strings.Replace(boxname, `/`, `-`, -1)

			// walk box path's and insert files
			boxPath := filepath.Clean(filepath.Join(pkg.Dir, boxname))

			boxMap[appendedBoxName] = boxPath
		}
	}

	operationAppend(binfileName, boxMap)
}

func operationAppend(binfileName string, boxMap map[string]string) {
	// create tmp zipfile
	tmpZipfileName := filepath.Join(os.TempDir(), fmt.Sprintf("ricebox-%d-%s.zip", time.Now().Unix(), randomString(10)))
	verbosef("Will create tmp zipfile: %s\n", tmpZipfileName)
	tmpZipfile, err := os.Create(tmpZipfileName)
	if err != nil {
		fmt.Printf("Error creating tmp zipfile: %s\n", err)
		os.Exit(1)
	}
	defer func() {
		tmpZipfile.Close()
		os.Remove(tmpZipfileName)
	}()

	verbosef("Will append to file: %s\n", binfileName)

	// check that command doesn't already have zip appended
	if rd, _ := zipexe.Open(binfileName); rd != nil {
		fmt.Printf("Cannot append to already appended executable. Please remove %s and build a fresh one.\n", binfileName)
		os.Exit(1)
	}

	// open binfile
	binfile, err := os.OpenFile(binfileName, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Printf("Error: unable to open executable file: %s\n", err)
		os.Exit(1)
	}
	defer binfile.Close()

	binfileInfo, err := binfile.Stat()
	if err != nil {
		fmt.Printf("Error: unable to stat executable file: %s\n", err)
		os.Exit(1)
	}

	// create zip.Writer
	zipWriter := zip.NewWriter(tmpZipfile)

	// write the zip offset into the zip data
	zipWriter.SetOffset(binfileInfo.Size())

	for appendedBoxName, boxPath := range boxMap {
		filepath.Walk(boxPath, func(path string, info os.FileInfo, err error) error {
			if info == nil {
				fmt.Printf("Error: box \"%s\" not found on disk\n", path)
				os.Exit(1)
			}
			// create zipFilename
			zipFileName := filepath.Join(appendedBoxName, strings.TrimPrefix(path, boxPath))

			// write directories as empty file with comment "dir"
			if info.IsDir() {
				_, err := zipWriter.CreateHeader(&zip.FileHeader{
					Name:    zipFileName,
					Comment: "dir",
				})
				if err != nil {
					fmt.Printf("Error creating dir in tmp zip: %s\n", err)
					os.Exit(1)
				}
				return nil
			}

			// create zipFileWriter
			zipFileHeader, err := zip.FileInfoHeader(info)
			if err != nil {
				fmt.Printf("Error creating zip FileHeader: %v\n", err)
				os.Exit(1)
			}
			zipFileHeader.Name = zipFileName
			zipFileWriter, err := zipWriter.CreateHeader(zipFileHeader)
			if err != nil {
				fmt.Printf("Error creating file in tmp zip: %s\n", err)
				os.Exit(1)
			}
			srcFile, err := os.Open(path)
			if err != nil {
				fmt.Printf("Error opening file to append: %s\n", err)
				os.Exit(1)
			}
			_, err = io.Copy(zipFileWriter, srcFile)
			if err != nil {
				fmt.Printf("Error copying file contents to zip: %s\n", err)
				os.Exit(1)
			}
			srcFile.Close()

			return nil
		})
	}

	err = zipWriter.Close()
	if err != nil {
		fmt.Printf("Error closing tmp zipfile: %s\n", err)
		os.Exit(1)
	}

	err = tmpZipfile.Sync()
	if err != nil {
		fmt.Printf("Error syncing tmp zipfile: %s\n", err)
		os.Exit(1)
	}
	_, err = tmpZipfile.Seek(0, 0)
	if err != nil {
		fmt.Printf("Error seeking tmp zipfile: %s\n", err)
		os.Exit(1)
	}
	_, err = binfile.Seek(0, 2)
	if err != nil {
		fmt.Printf("Error seeking bin file: %s\n", err)
		os.Exit(1)
	}

	_, err = io.Copy(binfile, tmpZipfile)
	if err != nil {
		fmt.Printf("Error appending zipfile to executable: %s\n", err)
		os.Exit(1)
	}
}

// operationAppendBoxes doesn't need access to source code, so is useful if you want to ship a generic
// binary and allow others to append their files to that binary. It does however rely upon the user
// to specify all required box paths, but that doesn't seem unreasonable.
func operationAppendBoxes(boxes []string) {
	// find abs path for binary file
	binfileName, err := filepath.Abs(flags.AppendBox.Executable)
	if err != nil {
		fmt.Printf("Error finding absolute path for executable to append: %s\n", err)
		os.Exit(1)
	}

	boxMap := make(map[string]string)
	for _, box := range boxes {
		boxMap[box] = box
	}
	operationAppend(binfileName, boxMap)
}
