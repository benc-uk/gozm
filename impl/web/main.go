//go:build js && wasm

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
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

	filenameOnly := path.Base(file)

	filenameOnly = filenameOnly[:len(filenameOnly)-len(path.Ext(filenameOnly))]
	machine := zmachine.NewMachine(data, filenameOnly, zmachine.DEBUG_NONE, ext)

	// We wrap the main run loop to handle restarts and loads
	for {
		exitCode := machine.Run()

		js.Global().Call("clearScreen")

		switch exitCode {
		case zmachine.EXIT_LOAD:
			ext.info("Loading saved game...\n")
			machine = ext.Load(filenameOnly)
		case zmachine.EXIT_QUIT:
			ext.info("Quitting game...\n")
			return
		case zmachine.EXIT_RESTART:
			ext.info("Restarting game...\n")
			machine = zmachine.NewMachine(data, filenameOnly, zmachine.DEBUG_NONE, ext)
		default:
			return
		}
	}
}
