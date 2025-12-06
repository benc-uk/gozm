//go:build js && wasm

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"syscall/js"
	"time"

	"github.com/benc-uk/gozm/internal/zmachine"
)

func main() {
	file := os.Args[0]
	fmt.Printf("Starting GOZM WebAssembly: %s\n", file)

	ext := NewWebExternal()

	url := "stories/" + file
	ext.TextOut("Loading: " + url + "\n")
	js.Global().Set("inputSend", js.FuncOf(ext.ReceiveInput))

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Sleep based on file size to allow loading message to be visible
	sleepDuration := time.Duration(len(data)/10) * time.Millisecond
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	elapsed := time.Duration(0)
	for elapsed < sleepDuration {
		<-ticker.C
		ext.TextOut(".")
		elapsed += 100 * time.Millisecond
	}
	ext.TextOut("\n")

	js.Global().Call("clearScreen")
	js.Global().Call("loadedFile")

	m := zmachine.NewMachine(data, 0, ext)
	m.Run() // Note this will block until the Z-machine exits
}
