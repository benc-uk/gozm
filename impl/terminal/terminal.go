package main

import (
	"bufio"
	"fmt"
	"os"

	"encoding/json"

	"github.com/benc-uk/gozm/internal/zmachine"
)

// Implements a simple terminal interface for Z-machine IO
type Terminal struct {
}

func NewTerminal() *Terminal {
	return &Terminal{}
}

// TextOut outputs text to the console
func (t *Terminal) TextOut(text string) {
	fmt.Printf("%s", text)
}

// ReadInput reads a line of input from the console
func (t *Terminal) ReadInput() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return text
}

func (t *Terminal) PlaySound(soundID uint16, effect uint16, volume uint16) {
	info("Playing sound ID:%d effect:%d volume:%d\n", soundID, effect, volume)
}

// Save saves a snapshot of the machine state to a JSON file on disk
func (t *Terminal) Save(state *zmachine.SaveState) bool {
	// Snapshot all of the machine state to a file
	savePath := getSaveFullPath(state.Name)
	file, err := os.Create(savePath)
	if err != nil {
		fmt.Printf("Error creating save file: %v\n", err)
		return false
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(state); err != nil {
		fmt.Printf("Error encoding save data: %v\n", err)
		return false
	}

	info("Game saved to %s\n", savePath)
	return true
}

// Load loads a saved machine state from a JSON file on disk
func (t *Terminal) Load(name string, machine *zmachine.Machine) bool {
	savePath := getSaveFullPath(name)
	file, err := os.Open(savePath)
	if err != nil {
		fmt.Printf("Error opening save file: %v\n", err)
		return false
	}
	defer file.Close()

	// Decode saved state
	var saveData zmachine.SaveState
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&saveData); err != nil {
		fmt.Printf("Error decoding save data: %v\n", err)
		return false
	}

	// Restore machine state
	machine.ReplaceState(&saveData)
	info("Game loaded from %s\n", savePath)
	return true
}

func info(format string, a ...interface{}) {
	fmt.Printf("\033[34m"+format+"\033[0m", a...)
}

func getSaveFullPath(name string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return homeDir + "/" + name + ".save"
}
