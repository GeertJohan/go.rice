package main

import (
	"fmt"
	"go/build"
	"log"
	"os"
)

func getPkgs() []*build.Package {
	// find package for path
	var pkgs []*build.Package
	for _, importPath := range flags.ImportPaths {
		pkg := pkgForPath(importPath)
		pkgs = append(pkgs, pkg)
	}
	return pkgs
}

func main() {
	// parser arguments
	parseArguments()

	// switch on the operation to perform
	switch flagsParser.Active.Name {
	case "embed", "embed-go":
		for _, pkg := range getPkgs() {
			operationEmbedGo(pkg)
		}
	case "embed-syso":
		log.Println("WARNING: embedding .syso is experimental..")
		for _, pkg := range getPkgs() {
			operationEmbedSyso(pkg)
		}
	case "append":
		operationAppendPkgs(getPkgs())
	case "appendbox":
		operationAppendBoxes(flags.AppendBox.Boxes)
	case "clean":
		for _, pkg := range getPkgs() {
			operationClean(pkg)
		}
	}

	// all done
	verbosef("\n")
	verbosef("rice finished successfully\n")
}

// helper function to get *build.Package for given path
func pkgForPath(path string) *build.Package {
	// get pwd for relative imports
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("error getting pwd (required for relative imports): %s\n", err)
		os.Exit(1)
	}

	// read full package information
	pkg, err := build.Import(path, pwd, 0)
	if err != nil {
		fmt.Printf("error reading package: %s\n", err)
		os.Exit(1)
	}

	return pkg
}

func verbosef(format string, stuff ...interface{}) {
	if flags.Verbose {
		log.Printf(format, stuff...)
	}
}
