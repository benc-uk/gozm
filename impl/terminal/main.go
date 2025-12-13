package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/benc-uk/gozm/internal/zmachine"
)

var version = "0.0.0"

func main() {
	info("GOZM: Go Z-Machine Runtime and VM v%s\n", version)

	debugLevel := 0
	fileName := ""
	flag.IntVar(&debugLevel, "debug", zmachine.DEBUG_NONE, "Set debug level (0=none, 1=step, 2=trace)")
	flag.StringVar(&fileName, "file", "", "Path to Z-machine story file to load")
	flag.StringVar(&fileName, "f", "", "Path to Z-machine story file to load")
	flag.Parse()

	if debugLevel < 0 || debugLevel > 2 {
		fmt.Printf("Invalid debug level %d, must be 0, 1, or 2\n", debugLevel)
		os.Exit(1)
	}

	if fileName == "" {
		// try to get from positional args
		if flag.NArg() > 0 {
			fileName = flag.Arg(0)
		} else {
			fmt.Printf("No story file specified\n")
			os.Exit(1)
		}
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}

	ext := NewTerminal()
	filenameOnly := path.Base(fileName)
	filenameOnly = filenameOnly[:len(filenameOnly)-len(path.Ext(filenameOnly))]
	machine := zmachine.NewMachine(data, filenameOnly, debugLevel, ext)

	exitCode := machine.Run()
	fmt.Printf("Program exited with code %d\n", exitCode)
	os.Exit(exitCode - 1)
}
