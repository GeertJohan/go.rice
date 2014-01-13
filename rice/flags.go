package main

import (
	"fmt"
	goflags "github.com/jessevdk/go-flags" // rename import to `goflags` (file scope) so we can use `var flags` (package scope)
	"go/build"

	"os"
)

var operationInfo = `Please perform an operation such as 'embed' or 'clean'.`

// operation ("embed", "clean", etc)
var operation string

// path on which to operate
var path string

// flags
var flags struct {
	Verbose bool `long:"verbose" description:"Show verbose debug information"`
}

// leftover args (not flags)
var args []string

// initFlags parses the given flags.
// when the user asks for help (-h or --help): the application exists with status 0
// when unexpected flags is given: the application exits with status 1
func parseArguments() {
	// check if there is at least an operation given
	if len(os.Args) < 2 {
		fmt.Println("No operation given.\n" + operationInfo)
		os.Exit(-1)
	}

	// get opration argument
	operation = os.Args[1]

	// parse generic flags (operation independent)
	// save leftover args on package-wide args var (to be used for operation-depending flags)
	args = parseFlags(&flags, os.Args[2:])

	// operation base flags
	switch operation {
	case "embed":
		// fmt.Println("operation embed")
		// nothing special to do now
	case "clean":
		// fmt.Println("operation clean")
		// nothing special to do now
	default:
		fmt.Println("Invalid operation given.\n" + operationInfo)
		fmt.Println("Use: rice <operation> [go import path] [flags]")
		os.Exit(-1)
	}

	if len(args) > 1 {
		fmt.Printf("too many args, expected zero or one path.")
		os.Exit(-1)
	}

	if len(args) == 0 {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("error getting pwd: %s\n", err)
			os.Exit(-1)
		}
		// find non-absolute path for this pwd
		pkg, err := build.ImportDir(pwd, build.FindOnly)
		if err != nil {
			fmt.Printf("error using current directory as import path: %s\n", err)
			os.Exit(-1)
		}
		path = pkg.ImportPath
		verbosef("using pwd as path (%s)\n", path)
		return
	}

	if len(args) == 1 {
		path = args[0]
		args = make([]string, 0)
		verbosef("using argument as path (%s)\n", path)
		return
	}
}

func parseFlags(flags interface{}, args []string) []string {
	args, err := goflags.ParseArgs(flags, args)
	if err != nil {
		// assert the err to be a flags.Error
		flagError := err.(*goflags.Error)
		if flagError.Type == goflags.ErrHelp {
			// user asked for help on flags.
			// program can exit successfully
			os.Exit(0)
		}
		if flagError.Type == goflags.ErrUnknownFlag {
			fmt.Println("Use --help to view all available options.")
			os.Exit(1)
		}
		fmt.Printf("Error parsing flags: %s\n", err)
		os.Exit(1)
	}
	return args
}
