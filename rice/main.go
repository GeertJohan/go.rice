package main

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"runtime/pprof"
)

func main() {
	// parser arguments
	parseArguments()

	if flags.CpuProfile != "" {
		f, err := os.Create(flags.CpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

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
		operationAppend(getPkgs())
	case "clean":
		for _, pkg := range getPkgs() {
			operationClean(pkg)
		}
	}

	// all done
	verbosef("\n")
	verbosef("rice finished successfully\n")

	if flags.MemProfile != "" {
		f, err := os.Create(flags.MemProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
	}
}

// get all packages from import path
func getPkgs() []*build.Package {
	// find package for path
	var pkgs []*build.Package
	for _, importPath := range flags.ImportPaths {
		pkg := pkgForPath(importPath)
		pkg.AllTags = flags.Tags
		pkgs = append(pkgs, pkg)
	}
	return pkgs
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
