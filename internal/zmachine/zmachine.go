package zmachine

import (
	"fmt"
	"gozm/internal/decode"
	"os"
)

type Machine struct {
	mem []byte
	pc  uint16

	version     byte
	globalsAddr uint16
}

func NewMachine(data []byte) *Machine {
	fmt.Printf("Z-machine initialized...\nVersion: %d, Size: %d\n", data[0x00], len(data))

	return &Machine{
		mem:         data,
		pc:          decode.GetWord(data, 0x06),
		version:     data[0x00],
		globalsAddr: decode.GetWord(data, 0x0C),
	}
}

// Step executes a single instruction at the current program counter
func (m *Machine) Step() {
	inst := m.decodeInst()
	fmt.Printf("\n%04x - %s\n", m.pc, inst.String())

	// Decode and execute instruction here
	// This is a ultra minimum implementation with only a few opcodes
	switch inst.code {
	case 0x00:
		fmt.Println("++ NOP")
		m.pc += 1

	case 0xBA:
		fmt.Println("++ QUIT")
		os.Exit(0)

	case 0x54: // ADD V,S -> result
		v := m.GetVar(inst.operands[0])
		s := inst.operands[1]
		dest := m.mem[m.pc+inst.len] // destination in next byte
		fmt.Printf("++ ADD %d, %d -> %d\n", v, s, dest)

		m.StoreVar(uint16(dest), v+s)
		m.pc += inst.len + 1 // +1 for dest byte
	case 0x0d: // STORE V,S
		v := inst.operands[0]
		s := inst.operands[1]
		fmt.Printf("++ STORE %d, %d\n", v, s)

		m.StoreVar(v, s)
		m.pc += inst.len
	default:
		fmt.Printf("!! UNKNOWN: %02x\n\n", inst.code)

		// Dump the first four globals, to prove we're working
		m.DumpMem(m.globalsAddr, 2*4)
		panic("BYE!")
	}
}

// StoreVar stores a value into a variable location
func (m *Machine) StoreVar(loc uint16, val uint16) {
	if loc == 0 {
		// TODO: store on stack
	} else if loc > 0 && loc < 0x10 {
		// TODO: store in local variable
	} else {
		// Global variable, which are all word sized
		addr := uint16(m.globalsAddr + (loc-0x10)*2)
		decode.SetWord(m.mem, addr, val)
	}
}

// GetVar retrieves a value from a variable location
func (m *Machine) GetVar(loc uint16) uint16 {
	if loc == 0 {
		// TODO: store on stack
	} else if loc > 0 && loc < 0x10 {
		// TODO: store in local variable
	} else {
		// Global variable, which are all word sized
		addr := uint16(m.globalsAddr + (loc-0x10)*2)
		return decode.GetWord(m.mem, addr)
	}
	return 0
}

// Run starts the main execution loop of the Z-machine
func (m *Machine) Run() {
	fmt.Println("Starting the machine...")
	// HACK: Point to known location otherwise we need to implement more opcodes
	m.pc = 0x0049F
	for {
		m.Step()
	}
}

// DumpMem dumps a section of memory for debugging
func (m *Machine) DumpMem(addr uint16, length uint16) {
	fmt.Printf("Memory dump at %04x:\n", addr)
	for i := uint16(0); i < length; i += 2 {
		word := decode.GetWord(m.mem, addr+i)
		fmt.Printf("%04x: %04x (%04d)\n", addr+i, word, word)
	}
}
