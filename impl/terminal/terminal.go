package main

import (
	"bufio"
	"fmt"
	"os"

	"encoding/json"

	"github.com/benc-uk/gozm/internal/zmachine"
	"github.com/chzyer/readline"
)

// Implements a simple terminal interface for Z-machine IO
type Terminal struct {
	history []string
	rl      *readline.Instance
}

func NewTerminal() *Terminal {
	t := &Terminal{
		history: make([]string, 0),
	}

	// Try to create a readline instance for better UX (arrow-key history)
	// Fall back to stdio if readline fails (e.g., unsupported terminal)
	cfg := readline.Config{
		Prompt:                 ">",
		HistoryLimit:           20,
		InterruptPrompt:        "",
		DisableAutoSaveHistory: true,
	}
	if inst, err := readline.NewEx(&cfg); err == nil {
		t.rl = inst
	}
	return t
}

// TextOut outputs text to the console
func (t *Terminal) TextOut(text string) {
	fmt.Printf("%s", text)
}

// ReadInput reads a line of input from the console
func (t *Terminal) ReadInput() string {
	// If readline is available, use it to provide history navigation
	if t.rl != nil {
		line, err := t.rl.Readline()
		if err == nil {
			trimmed := line
			if len(trimmed) > 0 && (len(t.history) == 0 || t.history[len(t.history)-1] != trimmed) {
				t.history = append(t.history, trimmed)
				// Keep a bounded history
				if len(t.history) >= 100 {
					t.history = t.history[1:]
				}
			}
			// Update readline's internal history so Up/Down work immediately
			if t.rl != nil && trimmed != "" {
				_ = t.rl.SaveHistory(trimmed)
			}
			return trimmed + "\n"
		}
		// If readline errors (e.g., EOF), fall back to stdio
	}

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	// Store in history if non-blank and not duplicate of last entry
	trimmed := string(text[:len(text)-1])
	if len(trimmed) > 0 && (len(t.history) == 0 || t.history[len(t.history)-1] != trimmed) {
		t.history = append(t.history, trimmed)
		if len(t.history) >= 20 {
			t.history = t.history[1:]
		}
	}

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
