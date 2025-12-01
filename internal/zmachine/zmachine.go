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
	mem          []byte
	pc           uint16
	callStack    []callFrame
	debugLevel   int
	ext          External
	propDefaults []uint16
	objects      map[byte]*zObject

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

func NewMachine(data []byte, debugLevel int, ext External) *Machine {

	m := &Machine{
		mem:          data,
		pc:           decode.GetWord(data, 0x06),
		callStack:    make([]callFrame, 0),
		debugLevel:   debugLevel,
		ext:          ext,
		propDefaults: make([]uint16, 31),
		objects:      make(map[byte]*zObject),

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
	}

	// Initialize objects, property defaults table
	m.initObjects()

	m.debug("Z-machine initialized...\nVersion: %d, Size: %d\n", data[0x00], len(data))
	m.debug(" - High Memory Address: %04x\n", m.highAddr)
	m.debug(" - Initial PC: %04x\n", m.initialPC)
	m.debug(" - Globals Address: %04x\n", m.globalsAddr)

	// Initialize the stack with the main__ call frame
	m.addCallFrame(0)

	return m
}

// Run starts the main execution loop of the Z-machine
func (m *Machine) Run() {
	m.debug("Starting the main execution loop...\n")

	// We just loop forever for now, this is our life
	for {
		m.Step()
	}
}

// Step executes a single instruction at the current program counter
func (m *Machine) Step() {
	inst := m.decodeInst()
	m.debug("\n%04X: %s\n", m.pc, inst.String())

	// Decode and execute instructions!
	switch inst.code {
	// ===================== 0OP INSTRUCTIONS =====================
	// NOP
	case 0x00:
		m.pc += inst.len

	// RET_TRUE
	case 0xB0:
		m.returnFromCall(1)

	// RET_FALSE
	case 0xB1:
		m.returnFromCall(0)

	// PRINT (literal string)
	case 0xB2:
		str, wordCount := m.readStringLiteral()
		m.ext.TextOut(str)
		m.pc += uint16(wordCount*2) + 1 // Advance PC past the string

	// PRINT_RET (literal string)
	case 0xB3:
		str, _ := m.readStringLiteral()
		m.ext.TextOut(str)
		m.returnFromCall(1)

	// NOP (Never used!)
	case 0xB4:
		m.debug("++ NOP (B4)\n")
		m.pc += inst.len

	// SAVE
	case 0xB5:
		panic("NOT_IMPLEMENTED: SAVE")

	// RESTORE
	case 0xB6:
		panic("NOT_IMPLEMENTED: RESTORE")

	// RESTART
	case 0xB7:
		panic("NOT_IMPLEMENTED: RESTART")

	// QUIT
	case 0xBA:
		m.debug("QUIT instruction encountered, exiting...\n")
		if m.debugLevel > DEBUG_NONE {
			m.DumpMem(m.globalsAddr, 24)
		}
		os.Exit(0)

	// NEW_LINE
	case 0xBB:
		m.ext.TextOut("\n")
		m.pc += inst.len

	// SHOW_STATUS
	case 0xBC:
		score := m.GetVar(17) // global variable 17 is score
		turns := m.GetVar(16) // global variable 16 is turns
		// TODO: Placeholder for scoreboard/status line handling, use ansi code to invert colors
		m.ext.TextOut(fmt.Sprintf("\n\033[32m\033[7m Unknown location                  score:%d turns:%d \033[27m\033[0m\n", score, turns))
		m.pc += inst.len

	// VERIFY
	case 0xBD:
		// Not worth implementing, just skip
		m.pc += inst.len

	// ===================== 1OP INSTRUCTIONS =====================

	// JZ
	case 0x80, 0x90, 0xA0:
		val := inst.operands[0]
		m.branchHandler(inst.len, val == 0)

	// JUMP
	case 0x8C, 0x9C, 0xAC:
		offset := decode.Convert14BitToSigned(inst.operands[0])
		m.pc = uint16(int16(m.pc) + int16(inst.len) + offset - 2)

	// GET_SIBLING
	// case 0xE1:
	// 	objNum := inst.operands[0]
	// 	sibling := m.getSibling(objNum)
	// 	m.debug(" - get_sibling of obj %d = %d\n", objNum, sibling)

	// 	// Store result in variable specified in next byte
	// 	dest := m.mem[m.pc+inst.len]
	// 	m.debug(" - store in var:%d\n", dest)
	// 	m.StoreVar(uint16(dest), sibling)

	// 	m.pc += inst.len + 1 // +1 for dest byte

	// PRINT_NUM
	case 0xE6:
		v := inst.operands[0]
		m.ext.TextOut(fmt.Sprintf("%d", v))
		m.pc += inst.len

	// ADD
	case 0x14, 0x34, 0x54, 0x74, 0xD4:
		v := inst.operands[0]
		s := inst.operands[1]
		dest := m.mem[m.pc+inst.len] // destination in next byte
		m.debug(" - add dest var:%d\n", dest)

		m.StoreVar(uint16(dest), v+s)
		m.pc += inst.len + 1 // +1 for dest byte

	// STORE
	case 0x0D, 0x2D:
		v := inst.operands[0]
		s := inst.operands[1]

		m.StoreVar(v, s)
		m.pc += inst.len

	// CALL
	case 0xE0:
		routineAddr := decode.PackedAddress(inst.operands[0])

		// Count locals from routine header
		numLocals := m.mem[routineAddr]
		m.debug(" - call to %04x with %d locals\n", routineAddr, numLocals)

		// Push new stack frame
		frame := m.addCallFrame(int(numLocals))
		frame.returnAddr = m.pc + inst.len

		// Populate locals (word sized) from the routine header
		// Note: Many compilers don't initialize locals, so this step may be unnecessary
		for i := byte(0); i < numLocals; i++ {
			localVal := decode.GetWord(m.mem, routineAddr+1+uint16(i*2))
			m.trace(" - local init %d = %d\n", i, localVal)
			frame.locals[i] = localVal
		}

		if len(inst.operands) > 1 {
			// Push arguments into local variables
			for i, argVal := range inst.operands[1:] {
				frame.locals[i] = argVal
				m.trace(" - arg %d = %d\n", i, argVal)
			}
		}

		// Set PC to start of routine after header and locals
		m.pc = routineAddr + 1 + uint16(numLocals*2)

	// RET_POPPED
	case 0xB8:
		val := m.getCallFrame().Pop()
		m.returnFromCall(val)

	// RET
	case 0x8B, 0x9B, 0xAB:
		val := inst.operands[0]
		m.returnFromCall(val)

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

func (m *Machine) branchHandler(instLen uint16, condition bool) {
	branchInfo := m.mem[m.pc+instLen]
	// Decode branch info
	branchOnTrue := (branchInfo & 0x80) != 0
	bit6Set := (branchInfo & 0x40) != 0
	branchDataLen := int16(1)

	var offset int16 // signed offset
	if bit6Set {
		// 6 bit offset 0-63
		offset = int16(branchInfo & 0x3F)
	} else {
		// 14 bit offset from next two bytes
		nextByte := m.mem[m.pc+instLen+1]
		branchDataLen = 2
		// If bit 6 is clear, then the offset is a signed 14-bit number given in bits 0 to 5 of the first byte followed by all 8 of the second.
		offset14 := (uint16(branchInfo&0x3F) << 8) | uint16(nextByte)

		offset = decode.Convert14BitToSigned(offset14)
	}

	m.debug(" - branchOnTrue: %t, condition: %t (info:%02x) offset:%d\n", branchOnTrue, condition, branchInfo, offset)

	if condition == branchOnTrue {
		m.pc = uint16(int16(m.pc) + int16(instLen) + branchDataLen + offset - 2)
		m.debug("   -> branching to %04x\n", m.pc)
	} else {
		m.pc += instLen + uint16(branchDataLen)
		m.debug("   -> no branch, next pc %04x\n", m.pc)
	}
}
