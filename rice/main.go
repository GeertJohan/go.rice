package main

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"path"
)

func main() {
	log.Println("Jon's version")
	// parser arguments
	parseArguments()

	// find package for path
	var pkgs []*build.Package
	var boxes = make(map[string]bool)

	for _, boxPath := range flags.AppendSimple.BoxPath {
		boxPath, value := pathForBoxPath(boxPath)
		boxes[boxPath] = value
		verbosef("box path %q added and exists=%t\n", boxPath, value)
	}

	for _, importPath := range flags.ImportPaths {
		pkg := pkgForPath(importPath)
		pkgs = append(pkgs, pkg)
	}

	// switch on the operation to perform
	switch flagsParser.Active.Name {
	case "embed", "embed-go":
		for _, pkg := range pkgs {
			operationEmbedGo(pkg)
		}
	case "embed-syso":
		log.Println("WARNING: embedding .syso is experimental..")
		for _, pkg := range pkgs {
			operationEmbedSyso(pkg)
		}
	case "append-simple":
		operationAppendSimple(boxes)
	case "append":
		operationAppend(pkgs)
	case "clean":
		for _, pkg := range pkgs {
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

// helper function to get *build.Package for given path
func pathForBoxPath(boxPath string) (string,bool) {
	// get pwd for relative imports
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("error getting pwd (required for relative box paths): %s\n", err)
		os.Exit(1)
	}

	if path.IsAbs(boxPath) {
		boxPath = boxPath
	} else {
		boxPath = path.Join(pwd, boxPath)
	}

	result, err := exists(boxPath)
	if err != nil {
		fmt.Printf("error finding box path %q : %s\n", boxPath, err)
	}
	return boxPath, result
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
func verbosef(format string, stuff ...interface{}) {
	if flags.Verbose {
		log.Printf(format, stuff...)
	}
}
