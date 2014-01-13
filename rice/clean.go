package main

import (
	"fmt"
	"go/build"
	"os"
	"strings"
)

func operationClean(pkg *build.Package) {
	for _, filename := range pkg.GoFiles {
		if strings.HasSuffix(filename, ".rice-box.go") || strings.HasSuffix(filename, ".rice-single.go") {
			err := os.Remove(filename)
			if err != nil {
				fmt.Printf("error removing file (%s): %s\n", filename, err)
				os.Exit(-1)
			}
		}
	}
}
