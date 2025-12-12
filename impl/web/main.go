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
var uploadedFileData []byte // Slice to store uploaded file data
var fileDataReady chan bool // Channel to signal when file data is ready
var machine *zmachine.Machine
var ext *WebExternal
var bridge js.Value // JavaScript WASM bridge

// If user selects a file to upload & run, this function is called from JS to pass the data
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

	// Signal that data is ready, we have to synchronously unblock the main thread
	if fileDataReady != nil {
		fileDataReady <- true
	}

	return nil
}

func main() {
	file := os.Args[0]
	fmt.Printf("Starting GOZM WebAssembly: %s\n", file)

	ext = NewWebExternal()
	bridge = js.Global().Get("bridge")
	bridge.Set("inputSend", js.FuncOf(ext.receiveInput))
	bridge.Set("receiveFileData", js.FuncOf(receiveFileData))
	bridge.Set("save", js.FuncOf(save))
	bridge.Set("load", js.FuncOf(load))
	bridge.Set("printInfo", js.FuncOf(printInfo))

	var data []byte

	// We either load from uploaded temp file or from URL
	if file == "tempFile" {
		// Initialize channel and wait for file data
		fileDataReady = make(chan bool, 1)
		fmt.Println("Waiting for tempFile data...")
		<-fileDataReady
		ext.TextOut("Loading: RAM:/upload/temp.z3")
		data = uploadedFileData
	} else {
		url := "stories/" + file
		ext.TextOut("Loading: DF0:/games/" + url + "\n")

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
	sleepDuration := time.Duration(len(data)/30) * time.Millisecond
	ticker := time.NewTicker(80 * time.Millisecond)
	defer ticker.Stop()
	elapsed := time.Duration(0)
	for elapsed < sleepDuration {
		<-ticker.C
		ext.TextOut(".")
		elapsed += 80 * time.Millisecond
	}
	ext.TextOut("\n")

	bridge.Call("loadedFile", file)

	filenameOnly := path.Base(file)
	filenameOnly = filenameOnly[:len(filenameOnly)-len(path.Ext(filenameOnly))]

	machine = zmachine.NewMachine(data, filenameOnly, zmachine.DEBUG_NONE, ext)

	// Everything is about this one line
	exitCode := machine.Run()

	if exitCode == zmachine.EXIT_RESTART {
		ext.TextOut("Restarting the game...\n")
		// For web, we just reload the page
		js.Global().Get("location").Call("reload")
		return
	}

	fmt.Printf("Game exited with code: %d\n", exitCode)
	os.Exit(exitCode - 1)
}

func save(this js.Value, args []js.Value) interface{} {
	if machine == nil || ext == nil {
		fmt.Println("No machine or external interface available for saving")
		return nil
	}

	success := ext.Save(machine.GetSaveState())
	if success {
		ext.TextOut("Game saved successfully.\n")
	} else {
		ext.TextOut("Error saving game.\n")
	}

	return nil
}

func load(this js.Value, args []js.Value) interface{} {
	if machine == nil || ext == nil {
		fmt.Println("No machine or external interface available for loading")
		return nil
	}

	loadOK := ext.Load(machine.GetName(), machine)
	if loadOK {
		ext.TextOut("Game loaded successfully.\n")
	} else {
		ext.TextOut("Error loading game.\n")
	}

	return nil
}

func printInfo(this js.Value, args []js.Value) interface{} {
	if machine == nil {
		ext.TextOut("No machine available to print info.\n")
		return nil
	}

	ext.TextOut(machine.GetInfo())

	return nil
}
