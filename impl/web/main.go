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

// Global variable to hold file data passed from JavaScript
var uploadedFileData []byte
var fileDataReady chan bool

// Function called from JavaScript to pass file data
func receiveFileData(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		fmt.Println("Error: no file data received")
		return nil
	}

	// Get the Uint8Array from JavaScript
	jsArray := args[0]
	length := jsArray.Get("length").Int()

	// Allocate Go slice and copy data
	uploadedFileData = make([]byte, length)
	js.CopyBytesToGo(uploadedFileData, jsArray)

	fmt.Printf("Received file data: %d bytes\n", len(uploadedFileData))

	// Signal that data is ready
	if fileDataReady != nil {
		fileDataReady <- true
	}

	return nil
}

func main() {
	file := os.Args[0]
	fmt.Printf("Starting GOZM WebAssembly: %s\n", file)

	ext := NewWebExternal()
	js.Global().Set("inputSend", js.FuncOf(ext.ReceiveInput))
	js.Global().Set("receiveFileData", js.FuncOf(receiveFileData))

	var data []byte
	if file == "tempFile" {
		// Initialize channel and wait for file data
		fileDataReady = make(chan bool, 1)
		fmt.Println("Waiting for file data...")
		<-fileDataReady
		fmt.Printf("Loading temporary file: %d bytes\n", len(uploadedFileData))
		data = uploadedFileData
	} else {
		url := "stories/" + file
		ext.TextOut("Loading: DF1:/games/" + url + "\n")

		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
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
