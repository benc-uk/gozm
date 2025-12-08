//go:build js && wasm

package main

import (
	"encoding/json"
	"fmt"
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

func (w *WebExternal) PlaySound(soundID uint16, effect uint16, volume uint16) {
	js.Global().Call("playSound", soundID, effect, volume)
}

func (w *WebExternal) Load(name string, machine *zmachine.Machine) bool {
	// Access localStorage to get saved game data via js
	savedData := js.Global().Get("localStorage").Call("getItem", name+"_save")
	if savedData.IsNull() || savedData.IsUndefined() || savedData.String() == "" {
		w.info("No saved game file found: " + name + "_save\n")
		return false
	}

	var saveData zmachine.SaveState

	err := json.Unmarshal([]byte(savedData.String()), &saveData)
	if err != nil {
		w.info("Error decoding saved game data: " + err.Error() + "\n")
		return false
	}

	w.info("Game loaded from browser storage: " + name + "_save\n")

	// Restore machine state
	machine.ReplaceState(&saveData)
	return true
}

func (w *WebExternal) Save(state *zmachine.SaveState) bool {
	// Snapshot all of the machine state to localStorage
	data, err := json.Marshal(state)
	if err != nil {
		fmt.Printf("Error encoding save data: %v\n", err)
		return false
	}

	js.Global().Get("localStorage").Call("setItem", state.Name+"_save", string(data))

	w.info("Game saved to DF0:/saves/" + state.Name + "_save\n")
	return true
}

func (w *WebExternal) info(format string, a ...interface{}) {
	w.TextOut(fmt.Sprintf("+++ "+format, a...))
}
