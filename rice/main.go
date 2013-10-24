package main

import (
	"fmt"
)

func main() {
	parseArguments()

	switch operation {
	case "embed":
		fmt.Printf("embedding boxes for '%s'\n", path)
		operationEmbed(path)
	case "clean":
		fmt.Printf("cleaning embedded boxes for '%s'\n", path)
		operationClean(path)
	}

	if flags.Verbose {
		fmt.Println("verbose")
	}

	fmt.Println("all done")
}
