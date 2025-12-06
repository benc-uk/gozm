//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/benc-uk/gozm/internal/zmachine"
)

// Implements a simple web/wasm interface for Z-machine IO
type WebExternal struct {
	inputChan    chan string
	inputWaiting bool
}

func NewWebExternal() *WebExternal {
	ext := &WebExternal{
		inputChan: make(chan string, 1),
	}

	return ext
}

func (w *WebExternal) TextOut(text string) {
	js.Global().Call("textOut", text)
}

func (w *WebExternal) ReadInput() string {
	w.inputWaiting = true
	js.Global().Call("requestInput")

	// Wait for input to be sent via the inputChan
	input := <-w.inputChan
	return input
}

func (w *WebExternal) ReceiveInput(this js.Value, args []js.Value) interface{} {
	if !w.inputWaiting {
		return nil
	}

	input := args[0].String()

	// Echo input back to output
	w.TextOut(input + "\n")

	w.inputChan <- input
	w.inputWaiting = false
	return nil
}

func (w *WebExternal) SetOutputStream(streamNum byte, tableAddr uint16, m *zmachine.Machine) {

}
