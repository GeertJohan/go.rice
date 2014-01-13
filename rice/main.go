package main

import (
	"fmt"
	"go/build"
	"log"
	"os"
)

func main() {
	//++ TODO: use less globals, more return values
	parseArguments()

	// switch on the operation to perform
	switch operation {
	case "embed":
		pkg := pkgForPath(path)
		operationEmbed(pkg)
	case "clean":
		pkg := pkgForPath(path)
		operationClean(pkg)
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
		os.Exit(-1)
	}

	// read full package information
	pkg, err := build.Import(path, pwd, 0)
	if err != nil {
		fmt.Printf("error reading package: %s\n", err)
		os.Exit(-1)
	}

	return pkg
}

func verbosef(format string, stuff ...interface{}) {
	if flags.Verbose {
		log.Printf(format, stuff...)
	}
}
