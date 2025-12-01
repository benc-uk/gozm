package main

import (
	"flag"
	"fmt"
	"gozm/externals/terminal"
	"gozm/internal/zmachine"
	"os"
)

var version = "0.0.1alpha"

func main() {
	fmt.Printf("GOZM: Go Z-Machine Runtime and VM v%s\n", version)

	debugLevel := 0
	fileName := ""
	flag.IntVar(&debugLevel, "debug", zmachine.DEBUG_NONE, "Set debug level (0=none, 1=step, 2=trace)")
	flag.StringVar(&fileName, "file", "", "Path to Z-machine story file to load")
	flag.Parse()

	if debugLevel < 0 || debugLevel > 2 {
		fmt.Printf("Invalid debug level %d, must be 0, 1, or 2\n", debugLevel)
		os.Exit(1)
	}

	if fileName == "" {
		panic("No story file specified, use -file to provide a path")
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}

	machine := zmachine.NewMachine(data, debugLevel, terminal.NewTerminal())

	machine.Run()
}
