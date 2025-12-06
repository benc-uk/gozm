//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/benc-uk/gozm/internal/zmachine"
)

// Implements a simple web/wasm interface for Z-machine IO
type WebExternal struct {
	inputChan chan string
}

func NewWebExternal() *WebExternal {
	ext := &WebExternal{
		inputChan: make(chan string, 1),
	}

	return ext
}

func (c *WebExternal) TextOut(text string) {
	//fmt.Printf("|..%s..|", text)
	if text == "\n>" {
		return
	}

	js.Global().Call("textOut", text)
}

func (c *WebExternal) ReadInput() string {
	js.Global().Call("requestInput")

	// Wait for input to be sent via the inputChan
	input := <-c.inputChan
	return input
}

func (c *WebExternal) ReceiveInput(this js.Value, args []js.Value) interface{} {
	input := args[0].String()
	c.inputChan <- input
	return nil
}

func (c *WebExternal) SetOutputStream(streamNum byte, tableAddr uint16, m *zmachine.Machine) {

}
