package main

import (
	"fmt"
)

func main() {
	parseArguments()

	switch operation {
	case "embed":
		//++ take pwd and continue with that as path
		//++ future: check len(args) to equal 0 or 1 (otherwise error). on 0 use pwd. on 1 use arg as path (abs or rel) to pkg to run embed on
	case "clean":
		//++ clean embed files
		//++ future: check len(args) to equal 0 or 1 (otherwise error). on 0 use pwd. on 1 use arg as path (abs or rel) to pkg to run embed on
	}

	if flags.Verbose {
		fmt.Println("verbose")
	}

	fmt.Println("all done")
}
