//go:build js && wasm

package main

import (
	"fmt"
	"io"
	"net/http"
	"syscall/js"

	"github.com/benc-uk/gozm/internal/zmachine"
)

func main() {
	// fmt.Println(os.Args)
	// // Load the z3 file using HTTP GET request URL in first argument.
	// url := os.Args[0]
	// fmt.Println("Loading story file from:", url)
	url := "moonglow.z3"
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

	ext := NewWebExternal()
	m := zmachine.NewMachine(data, 0, ext)

	// Set up inputSend BEFORE running the machine
	js.Global().Set("inputSend", js.FuncOf(ext.ReceiveInput))

	m.Run()

	select {} // keep running
}
