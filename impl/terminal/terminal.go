package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"encoding/json"

	"github.com/benc-uk/gozm/internal/zmachine"
	"github.com/peterh/liner"
)

// Implements a simple terminal interface for Z-machine IO
type Terminal struct {
	liner         *liner.State
	pendingPrompt string // Text output that didn't end with newline (used as prompt)
}

func NewTerminal() *Terminal {
	t := &Terminal{}

	// Try to create a liner instance for better UX (arrow-key history)
	l := liner.NewLiner()
	if l != nil {
		l.SetCtrlCAborts(true)
		t.liner = l
	}
	return t
}

// TextOut outputs text to the console
// This is made more complex by the need to track text that might be a prompt
func (t *Terminal) TextOut(text string) {
	// Track text that doesn't end with newline as pending prompt
	if strings.HasSuffix(text, "\n") {
		// Has newline - print everything and clear pending prompt
		fmt.Printf("%s%s", t.pendingPrompt, text)
		t.pendingPrompt = ""
	} else {
		// No newline - find the last line to use as prompt
		lastNewline := strings.LastIndex(text, "\n")
		if lastNewline >= 0 {
			// Print everything up to and including the last newline
			fmt.Printf("%s%s", t.pendingPrompt, text[:lastNewline+1])
			// Store the remainder as pending prompt
			t.pendingPrompt = text[lastNewline+1:]
		} else {
			// No newlines at all - append to pending prompt
			t.pendingPrompt += text
		}
	}
	os.Stdout.Sync()
}

// ReadInput reads a line of input from the console
func (t *Terminal) ReadInput() string {
	// If liner is available, use it to provide history navigation
	if t.liner != nil {
		line, err := t.liner.Prompt(t.pendingPrompt)
		t.pendingPrompt = ""
		if err == nil {
			if len(line) > 0 {
				t.liner.AppendHistory(line)
			}
			return line + "\n"
		}
		// If liner errors (e.g., EOF/Ctrl-C), fall back to stdio
	}

	// Print the pending prompt before reading
	if t.pendingPrompt != "" {
		fmt.Print(t.pendingPrompt)
		t.pendingPrompt = ""
		os.Stdout.Sync()
	}

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
