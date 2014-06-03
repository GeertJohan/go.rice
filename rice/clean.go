package main

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"
)

func operationClean(pkg *build.Package) {
	files := make([]string, 0, len(pkg.GoFiles)+len(pkg.SysoFiles))
	files = append(files, pkg.GoFiles...)
	files = append(files, pkg.SysoFiles...)
	for _, filename := range files {
		verbosef("checking file '%s'\n", filename)
		if strings.HasSuffix(filename, ".rice-box.go") ||
			strings.HasSuffix(filename, ".rice-single.go") ||
			strings.HasSuffix(filename, ".rice-box_386.syso") ||
			strings.HasSuffix(filename, ".rice-box_amd64.syso") {
			err := os.Remove(filepath.Join(pkg.Dir, filename))
			if err != nil {
				fmt.Printf("error removing file (%s): %s\n", filename, err)
				os.Exit(-1)
			}
			verbosef("removed file '%s'\n", filename)
		}
	}
}
