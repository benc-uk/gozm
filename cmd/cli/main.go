package main

import (
	"fmt"
	"gozm/internal/zmachine"
	"os"
)

func main() {
	fmt.Println("👺 GOZM: Go Z-Machine Interpreter and VM")

	if len(os.Args) < 2 {
		fmt.Println("Usage: gozm <z-machine-file>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error loading Z-machine file: %v\n", err)
		os.Exit(1)
	}

	machine := zmachine.NewMachine(data)
	machine.Run()
}
