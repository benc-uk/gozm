package zmachine

import (
	"fmt"
	"gozm/internal/decode"
	"os"
)

// enum for debug levels
const (
	DEBUG_NONE  = 0
	DEBUG_STEP  = 1
	DEBUG_TRACE = 2
)

// Machine represents the state of a Z-machine interpreter
type Machine struct {
	mem        []byte
	pc         uint16
	callStack  []callFrame
	DebugLevel int

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
	fmt.Printf("Z-machine initialized...\nVersion: %d, Size: %d\n", data[0x00], len(data))

	m := &Machine{
		mem: data,
		pc:  decode.GetWord(data, 0x06),

		version:     data[0x00],
		highAddr:    decode.GetWord(data, 0x04),
		initialPC:   decode.GetWord(data, 0x06),
		dictAddr:    decode.GetWord(data, 0x08),
		objectsAddr: decode.GetWord(data, 0x0A),
		globalsAddr: decode.GetWord(data, 0x0C),
		staticAddr:  decode.GetWord(data, 0x0E),
		abbrvAddr:   decode.GetWord(data, 0x18),
		fileLen:     decode.GetWord(data, 0x1A),
		checksum:    decode.GetWord(data, 0x1C),

		callStack:  make([]callFrame, 0),
		DebugLevel: DEBUG_NONE,
	}

	// print high memory address and initial PC andglobals address
	fmt.Printf(" - High Memory Address: %04x\n", m.highAddr)
	fmt.Printf(" - Initial PC: %04x\n", m.initialPC)
	fmt.Printf(" - Globals Address: %04x\n", m.globalsAddr)

	// Initialize the stack with the main__ call frame
	m.addCallFrame()

	return m
}

// Run starts the main execution loop of the Z-machine
func (m *Machine) Run() {
	fmt.Printf("Starting the machine...\n\n")

	for {
		m.Step()
	}
}

// Step executes a single instruction at the current program counter
func (m *Machine) Step() {
	m.debug("\nPC: %04x\n", m.pc)
	inst := m.decodeInst()
	m.debug("%04x - %s\n", m.pc, inst.String())

	// Decode and execute instructions, draw the rest of the owl
	switch inst.code {
	// NOP
	case 0x00:
		m.debug("++ NOP")
		m.pc += inst.len

	// QUIT
	case 0xBA:
		m.debug("++ QUIT\n\n")
		m.DumpMem(m.globalsAddr, 20)
		os.Exit(0)

	// PRINT_NUM
	case 0xE6:
		v := inst.operands[0]
		fmt.Printf("%d", v)
		m.pc += inst.len

	// ADD
	case 0x14, 0x34, 0x54, 0x74:
		v := inst.operands[0]
		s := inst.operands[1]
		dest := m.mem[m.pc+inst.len] // destination in next byte
		m.debug("++ ADD %d, %d -> %d\n", v, s, dest)

		m.StoreVar(uint16(dest), v+s)
		m.pc += inst.len + 1 // +1 for dest byte

	// STORE
	case 0x0D, 0x2D:
		v := inst.operands[0]
		s := inst.operands[1]
		m.debug("++ STORE %d, %d\n", v, s)

		m.StoreVar(v, s)
		m.pc += inst.len

	// CALL
	case 0xE0:
		routineAddr := decode.PackedAddress(inst.operands[0])

		// Count locals from routine header
		numLocals := m.mem[routineAddr]
		m.debug("++ CALL to %04x with %d locals\n", routineAddr, numLocals)

		// Push new stack frame
		frame := m.addCallFrame()
		frame.returnAddr = m.pc + inst.len

		// Populate locals (word sized) from the routine header
		// Note: Many compilers don't initialize locals, so this step may be unnecessary
		for i := byte(0); i < numLocals; i++ {
			localVal := decode.GetWord(m.mem, routineAddr+1+uint16(i*2))
			m.trace("  Local init %d = %d\n", i, localVal)
			frame.locals[i] = localVal
		}

		if len(inst.operands) > 1 {
			// Push arguments into local variables
			for i, argVal := range inst.operands[1:] {
				if i < int(numLocals) {
					frame.locals[i] = argVal
					m.trace("  Arg %d = %d\n", i, argVal)
				}
			}
		}

		m.debug("++ CALL: %02x", routineAddr)

		// Set PC to start of routine after header and locals
		m.pc = routineAddr + 1 + uint16(numLocals*2)

	// RET_TRUE
	case 0xB0:
		m.debug("++ RET_TRUE")
		m.returnFromCall(1)

	// RET_FALSE
	case 0xB1:
		m.debug("++ RET_FALSE")
		m.returnFromCall(0)

	// RET_POPPED
	case 0xB8:
		m.debug("++ RET_POPPED")
		val := m.getCallFrame().Pop()
		m.returnFromCall(val)

	// RET
	case 0x8B, 0x9B, 0xAB:
		val := inst.operands[0]
		m.returnFromCall(val)

	// PRINT (literal string)
	case 0xB2:
		m.debug("++ PRINT\n")
		str, wordCount := m.readStringLiteral()
		fmt.Printf("%s", str)
		m.pc += uint16(wordCount*2) + 1 // Advance PC past the string

	// PRINT_RET (literal string)
	case 0xB3:
		m.debug("++ PRINT_RET\n")
		str, _ := m.readStringLiteral()
		fmt.Printf("%s", str)

		m.returnFromCall(1)

	// Unimplemented instruction!
	default:
		panic(fmt.Sprintf("Unimplemented instruction: %02x", inst.code))
	}
}

// StoreVar stores a value into a variable location
func (m *Machine) StoreVar(loc uint16, val uint16) {
	// We made loc uint16 for ease of use, now restrict to valid range
	if loc > 0xFF {
		panic(fmt.Sprintf("Variable location out of range: %02x", loc))
	}

	if loc == 0 {
		m.getCallFrame().Push(val)
	} else if loc > 0 && loc < 0x10 {
		// Local variable
		m.getCallFrame().locals[loc-1] = val
	} else {

		// Global variable, which are all word sized
		addr := uint16(m.globalsAddr + (loc-0x10)*2)
		decode.SetWord(m.mem, addr, val)
	}
}

// GetVar retrieves a value from a variable location
func (m *Machine) GetVar(loc uint16) uint16 {
	// We made loc uint16 for ease of use, now restrict to valid range
	if loc > 0xFF {
		panic(fmt.Sprintf("Variable location out of range: %02x", loc))
	}

	if loc == 0 {
		// Stack variable
		return m.getCallFrame().Pop()
	} else if loc > 0 && loc < 0x10 {
		// Local variable
		return m.getCallFrame().locals[loc-1]
	} else {
		// Global variable, which are all word sized
		addr := uint16(m.globalsAddr + (loc-0x10)*2)
		return decode.GetWord(m.mem, addr)
	}
}

// Helper to return from a call with a value
func (m *Machine) returnFromCall(val uint16) {
	frame := m.getCallFrame()
	m.pc = frame.returnAddr
	m.callStack = m.callStack[:len(m.callStack)-1]

	// The next byte after a CALL is the variable to store the result in
	resultStoreLoc := m.mem[m.pc]
	m.debug(" > %04x into: %0x\n", val, resultStoreLoc)

	m.StoreVar(uint16(resultStoreLoc), val)
	m.pc += 1
}

// Read a Z-machine string literal: 2byte pairs from the current PC
// Returns the decoded string and number of words read
func (m *Machine) readStringLiteral() (string, int) {
	words := []uint16{}
	for i := uint16(1); int(i) < len(m.mem); i += 2 {
		word := decode.GetWord(m.mem, m.pc+i)
		words = append(words, word)

		// If the high bit is set, this is the end of the string
		if word&0x8000 != 0 {
			break
		}
	}

	return decode.String(words), len(words)
}

// DumpMem dumps a section of memory for debugging
func (m *Machine) DumpMem(addr uint16, length uint16) {
	fmt.Printf("\n\nMemory dump at %04x:\n", addr)
	for i := uint16(0); i < length; i += 2 {
		word := decode.GetWord(m.mem, addr+i)
		fmt.Printf("%04x: %04x (%04d)\n", addr+i, word, word)
	}
}

func (m *Machine) debug(format string, a ...interface{}) {
	if m.DebugLevel > DEBUG_NONE {
		fmt.Printf("\033[32m"+format+"\033[0m", a...)
	}
}

func (m *Machine) trace(format string, a ...interface{}) {
	if m.DebugLevel == DEBUG_TRACE {
		fmt.Printf("\033[36m"+format+"\033[0m", a...)
	}
}
