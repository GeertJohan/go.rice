package main

import (
	"fmt"
	goflags "github.com/jessevdk/go-flags" // rename import to `goflags` (file scope) so we can use `var flags` (package scope)
	"os"
)

var operationInfo = `Please perform an operation such as 'embed' or 'clean'.`

// operation ("embed", "clean", etc)
var operation string

// leftover args (not flags)
var args []string

// flags
var flags struct {
	Verbose bool `long:"verbose" description:"Show verbose debug information"`
}

// initFlags parses the given flags.
// when the user asks for help (-h or --help): the application exists with status 0
// when unexpected flags is given: the application exits with status 1
func parseArguments() {
	// check if there is at least an operation given
	if len(os.Args) == 1 {
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
		fmt.Println("will embd stuff..")
	case "clean":
		fmt.Println("TODO: clean embed files")
	default:
		fmt.Println("Invalid operation given.\n" + operationInfo)
		os.Exit(-1)
	}

	// check for unexpected arguments (at this point we don't expect any more arguments then the defiend flags and operations)
	// when an unexpected argument is given: the application exists with status 1
	if len(args) > 0 {
		fmt.Printf("Unknown argument '%s'.\n", args[0])
		os.Exit(1)
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
