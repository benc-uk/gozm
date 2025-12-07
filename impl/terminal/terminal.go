package main

import (
	"bufio"
	"fmt"
	"os"

	"encoding/json"

	"github.com/benc-uk/gozm/internal/decode"
	"github.com/benc-uk/gozm/internal/zmachine"
)

// Implements a simple terminal interface for Z-machine IO
type Terminal struct{}

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
	fmt.Printf("Playing sound ID:%d effect:%d volume:%d\n", soundID, effect, volume)
}

func (t *Terminal) Save(m *zmachine.Machine) bool {
	// Snapshot all of the machine state to a file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		return false
	}

	savePath := homeDir + "/" + m.GetName() + ".save"
	file, err := os.Create(savePath)
	if err != nil {
		fmt.Printf("Error creating save file: %v\n", err)
		return false
	}
	defer file.Close()

	memCopy := make([]byte, len(m.GetMem()))
	copy(memCopy, m.GetMem())

	// overwrite the PC in the memory copy
	decode.SetWord(memCopy, 0x06, uint16(m.GetPC()))

	// Create a JSON file with the memory and PC and call stack

	type SaveData struct {
		Memory    []byte               `json:"memory"`
		PC        uint32               `json:"pc"`
		CallStack []zmachine.CallFrame `json:"call_stack"`
	}

	fmt.Printf("!!!!! Saving PC=%08x\nCallStack=%+x\n", m.GetPC()+2, m.GetCallStack())

	saveData := SaveData{
		Memory:    memCopy,
		PC:        m.GetPC() + 2,
		CallStack: m.GetCallStack(),
	}

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(saveData); err != nil {
		fmt.Printf("Error encoding save data: %v\n", err)
		return false
	}

	fmt.Printf("Game saved to %s\n", savePath)
	return true
}

func (t *Terminal) Load(name string) *zmachine.Machine {
	// Load a snapshot of the machine state from a file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("Error getting home directory: %v\n", err))
	}

	savePath := homeDir + "/" + name + ".save"
	file, err := os.Open(savePath)
	if err != nil {
		panic(fmt.Sprintf("Error opening save file: %v\n", err))
	}
	defer file.Close()

	type SaveData struct {
		Memory    []byte               `json:"memory"`
		PC        uint32               `json:"pc"`
		CallStack []zmachine.CallFrame `json:"call_stack"`
	}

	var saveData SaveData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&saveData); err != nil {
		panic(fmt.Sprintf("Error decoding save data: %v\n", err))
	}

	// Create a new machine with the loaded memory and state
	machine := zmachine.NewMachine(saveData.Memory, name, zmachine.DEBUG_STEP, t)

	fmt.Printf("!!!!! Loaded save with PC=%08x\nCallStack=%+x\n", saveData.PC, saveData.CallStack)
	machine.SetPC(saveData.PC)
	machine.SetCallStack(saveData.CallStack)

	return machine
}
