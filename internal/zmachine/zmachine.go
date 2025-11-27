package zmachine

import (
	"fmt"
	"gozm/internal/decode"
	"os"
)

type Machine struct {
	mem []byte
	pc  uint16

	header header
}

type header struct {
	version     byte
	highAddr    uint16
	initialPC   uint16
	dictAddr    uint16
	objectsAddr uint16
	globalsAddr uint16
	staticAddr  uint16
	abbrvAddr   uint16
	fileLen     uint16
	checksum    uint16
}

func NewMachine(data []byte) *Machine {
	// Copy data into machine memory, very likely this is unnecessary
	memCopy := make([]byte, len(data))
	copy(memCopy, data)

	h := header{
		version:     memCopy[0x00],
		initialPC:   decode.GetWord(memCopy, 0x06),
		globalsAddr: decode.GetWord(memCopy, 0x0C),
	}
	fmt.Println("Z-machine initialized & loaded...")
	fmt.Printf("Version: %d, pc: %04x, globalsAddr:\n", h.version, h.globalsAddr)

	return &Machine{
		mem:    memCopy,
		pc:     h.initialPC,
		header: h,
	}
}

// Step executes a single instruction at the current program counter
func (m *Machine) Step() {
	inst := m.decodeInst()
	fmt.Printf("%04x - %s\n", m.pc, inst.String())

	// Decode and execute instruction here
	// This is a ultra minimum implementation with only a few opcodes
	switch inst.code {
	case 0x00:
		fmt.Println(" ++ NOP")
		m.pc += 1
	case 0xBA:
		fmt.Println(" ++ QUIT")
		os.Exit(0)
	case 0x54: // ADD V,S -> result
		v := m.GetVar(byte(inst.operands[0]))
		s := inst.operands[1]
		dest := m.mem[m.pc+inst.len] // destination in next byte
		fmt.Printf("++ ADD %d, %d -> %d\n", v, s, dest)

		m.StoreVar(dest, v+s)
		m.pc += inst.len + 1 // +1 for dest byte
	case 0x0d: // STORE V,S
		v := inst.operands[0]
		s := inst.operands[1]
		fmt.Printf("++ STORE %d, %d\n", v, s)

		m.StoreVar(byte(v), s)
		m.pc += inst.len
	default:
		fmt.Printf(" !! UNKNOWN OPCODE: %02x\n", inst.code)

		// Dump the first four globals
		m.DumpMem(m.header.globalsAddr, 2*4)
		panic("BYE")
	}
}

// StoreVar stores a value into a variable location
func (m *Machine) StoreVar(loc byte, val uint16) {
	if loc == 0 {
		// TODO: store on stack
	} else if loc > 0 && loc < 0x10 {
		// TODO: store in local variable
	} else {
		// Global variable
		addr := uint16(m.header.globalsAddr + ((uint16(loc) - 0x10) * 2))
		decode.SetWord(m.mem, addr, val)
	}
}

// GetVar retrieves a value from a variable location
func (m *Machine) GetVar(loc byte) uint16 {
	if loc == 0 {
		// TODO: store on stack
	} else if loc > 0 && loc < 0x10 {
		// TODO: store in local variable
	} else {
		// Global variable
		addr := uint16(m.header.globalsAddr + ((uint16(loc) - 0x10) * 2))
		return decode.GetWord(m.mem, addr)
	}
	return 0
}

// Run starts the main execution loop of the Z-machine
func (m *Machine) Run() {
	// HACK: Point to known location otherwise we need to implement more opcodes
	m.pc = 0x0049F
	for {
		m.Step()
	}
}

func (m *Machine) DumpMem(addr uint16, length uint16) {
	for i := uint16(0); i < length; i += 2 {
		word := decode.GetWord(m.mem, addr+i)
		fmt.Printf("%04x: %04x (%04d)\n", addr+i, word, word)
	}
}
