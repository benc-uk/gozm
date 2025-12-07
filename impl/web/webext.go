//go:build js && wasm

package main

import (
	"syscall/js"
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

func (w *WebExternal) PlaySound(soundID uint16, effect uint16, volume uint16) {
	js.Global().Call("playSound", soundID, effect, volume)
}

func (w *WebExternal) Save(mem []byte) bool {
	// Convert mem to a JS Uint8Array
	uint8Array := js.Global().Get("Uint8Array").New(len(mem))
	js.CopyBytesToJS(uint8Array, mem)

	js.Global().Call("saveGame", uint8Array)
	return true
}

func (w *WebExternal) Load() []byte {
	// Call JS function to get saved game data
	savedData := js.Global().Call("loadGame")
	if savedData.IsNull() || savedData.IsUndefined() {
		return []byte{}
	}

	// Convert JS Uint8Array back to Go byte slice
	length := savedData.Get("length").Int()
	mem := make([]byte, length)
	js.CopyBytesToGo(mem, savedData)
	return mem
}
