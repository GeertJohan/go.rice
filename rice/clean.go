package main

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
)

func operationClean(pkg *build.Package, outputPath string) {
	if len(outputPath) == 0 {
		outputPath = pkg.Dir
	}

	filepath.Walk(outputPath, func(filename string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("error walking pkg dir to clean files: %v\n", err)
			os.Exit(1)
		}
		if info.IsDir() {
			return nil
		}
		verbosef("checking file '%s'\n", filename)
		if generated(filename) {
			err := os.Remove(filename)
			if err != nil {
				fmt.Printf("error removing file (%s): %s\n", filename, err)
				os.Exit(-1)
			}
			verbosef("removed file '%s'\n", filename)
		}
		return nil
	})
}
