//go:build js && wasm

package main

import (
	"fmt"
	"io"
	"net/http"
	"syscall/js"
	"time"

	"github.com/benc-uk/gozm/internal/zmachine"
)

func main() {
	ext := NewWebExternal()

	url := "minizork.z3"
	ext.TextOut("Loading: " + url + "\n")
	js.Global().Set("inputSend", js.FuncOf(ext.ReceiveInput))

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error loading story file: %v\n", err)
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading story file: %v\n", err)
		return
	}

	// Sleep based on file size to allow loading message to be visible
	sleepDuration := time.Duration(len(data)/10) * time.Millisecond
	if sleepDuration < 2000*time.Millisecond {
		sleepDuration = 2000 * time.Millisecond
	}
	time.Sleep(sleepDuration)

	js.Global().Call("clearScreen")

	m := zmachine.NewMachine(data, 0, ext)
	m.Run()

	// Prevent main from exiting
	select {}
}
