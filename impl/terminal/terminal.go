package main

import (
	"bufio"
	"fmt"
	"os"
)

// Implements a simple terminal interface for Z-machine IO
type Terminal struct{}

func NewTerminal() *Terminal {
	return &Terminal{}
}

// TextOut outputs text to the console
func (c *Terminal) TextOut(text string) {
	fmt.Print(text)
}

// ReadInput reads a line of input from the console
func (c *Terminal) ReadInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	fmt.Println(text)
	return text
}
