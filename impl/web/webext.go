//go:build js && wasm

package main

import (
	"syscall/js"
	"time"

	"github.com/benc-uk/gozm/internal/zmachine"
)

// Implements a simple web/wasm interface for Z-machine IO
type WebExternal struct{}

func NewWebExternal() *WebExternal {
	ext := &WebExternal{}
	// js.Global().Set("textOut", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// 	text := args[0].String()
	// 	ext.TextOut(text)
	// 	return nil
	// }))

	// js.Global().Set("readInput", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// 	input := ext.ReadInput()
	// 	return input
	// }))

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
	time.Sleep(300 * time.Second) // Give the browser time to update the UI
	return "l"
}

func (c *WebExternal) SetOutputStream(streamNum byte, tableAddr uint16, m *zmachine.Machine) {

}
