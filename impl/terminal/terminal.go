package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/benc-uk/gozm/internal/zmachine"
)

// Implements a simple terminal interface for Z-machine IO
type Terminal struct{}

func NewTerminal() *Terminal {
	return &Terminal{}
}

// TextOut outputs text to the console
func (c *Terminal) TextOut(text string) {
	fmt.Printf("%s", text)
}

// ReadInput reads a line of input from the console
func (c *Terminal) ReadInput() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return text
}

func (c *Terminal) SetOutputStream(streamNum byte, tableAddr uint16, m *zmachine.Machine) {
	fmt.Printf("SetOutputStream called with streamNum=%d, tableAddr=%04x\n", streamNum, tableAddr)
}
