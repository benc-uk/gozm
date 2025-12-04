// =======================================================================
// Package: zmachine - Core Z-machine interpreter
// zmachine.go - Main code, structs and core execution loop
// Note: This file probably needs to be split up more later!
//
// Copyright (c) 2025 Ben Coleman. Licensed under the MIT License
// =======================================================================

package zmachine

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/benc-uk/gozm/internal/decode"
)

// enum for debug levels
const (
	DEBUG_NONE               = 0
	DEBUG_STEP               = 1
	DEBUG_TRACE              = 2
	OUTPUT_STREAM_SCREEN     = 1
	OUTPUT_STREAM_FILE       = 2
	OUTPUT_STREAM_MEMORY     = 3
	OUTPUT_STREAM_MEMORY_MAX = 16
	INPUT_STREAM_KEYBOARD    = 1
	INPUT_STREAM_FILE        = 2
)

// Machine represents the state of a Z-machine interpreter
type Machine struct {
	mem          []byte
	pc           uint16 // Actually think this probably needs to be a uint32
	callStack    []callFrame
	debugLevel   int
	ext          External
	propDefaults []uint16
	objects      []*zObject
	rand         *rand.Rand
	outputStream int
	inputStream  int
	TracedOps    []byte
	TracedObjs   []byte
	Breakpoint   uint16
	//memStreamAddrStack []uint16

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
		objects:      make([]*zObject, 0),
		rand:         rand.New(rand.NewPCG(123, 456)),
		outputStream: OUTPUT_STREAM_SCREEN,
		inputStream:  INPUT_STREAM_KEYBOARD,
		TracedOps:    make([]byte, 0),
		TracedObjs:   make([]byte, 0),
		//memStreamAddrStack: make([]uint16, 0),

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

	// Initialize abbreviations from the abbreviation table
	abbr := make([]string, 96)
	for i := uint16(0); i < 96; i++ {
		// Abbreviation table contains word addresses, need to multiply by 2
		// See: https://zspec.jaredreisinger.com/01-memory-map#1_2_2
		abbrStringAddr := decode.GetWord(m.mem, m.abbrvAddr+i*2) * 2
		s, _ := m.readStringLiteral(abbrStringAddr)
		abbr[i] = s
	}
	decode.InitAbbreviations(abbr)

	// Initialize objects, property defaults table
	m.initObjects()

	for _, o := range m.TracedObjs {
		obj := m.getObject(o)
		fmt.Printf("Traced Object %d '%s':\n%s\n", obj.num, obj.description, obj.propDebugDump())
	}

	m.debug("Z-machine initialized...\nVersion: %d, Size: %d\n", data[0x00], len(data))
	m.debug(" - High Memory Address: %04x\n", m.highAddr)
	m.debug(" - Initial PC: %04x\n", m.initialPC)
	m.debug(" - Globals Address: %04x\n", m.globalsAddr)
	m.debug(" - Checksum: %04X, valid:%t\n", m.checksum, m.validateChecksum())

	// Initialize the stack with the main__ call frame
	m.addCallFrame(0)

	return m
}

// Run starts the main execution loop of the Z-machine
func (m *Machine) Run() {
	m.debug("Starting the main execution loop...\n")

	// We just loop forever for now, this is our life
	for {
		m.step()
	}
}

// storeVar stores a value into a variable location
func (m *Machine) storeVar(loc uint16, val uint16) {
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

// getVar retrieves a value from a variable location
func (m *Machine) getVar(loc uint16) uint16 {
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

func (m *Machine) setVarInPlace(loc uint16, val uint16) {
	if loc == 0 {
		m.getCallFrame().SetTop(val)
		return
	}

	m.storeVar(loc, val)
}

func (m *Machine) addToVar(loc uint16, val int16) int16 {
	// We made loc uint16 for ease of use, now restrict to valid range
	if loc > 0xFF {
		panic(fmt.Sprintf("Variable location out of range: %02x", loc))
	}

	cf := m.getCallFrame()
	if loc == 0 {
		// Top of stack update in place
		curr := int16(cf.Peek())
		cf.SetTop(uint16(int16(curr) + val))
		return int16(curr + val)
	} else if loc > 0 && loc < 0x10 {
		// check local variable exists
		if int(loc-1) >= len(cf.locals) {
			panic(fmt.Sprintf("Local variable %d does not exist in routine, has %d", loc, len(cf.locals)))
		}

		// Local variable
		curr := int16(cf.locals[loc-1])
		cf.locals[loc-1] = uint16(curr + val)
		return int16(curr + val)
	} else {
		// Global variable, which are all word sized
		addr := uint16(m.globalsAddr + (loc-0x10)*2)
		curr := decode.GetWordSigned(m.mem, addr)
		decode.SetWord(m.mem, addr, uint16(curr+val))
		return int16(curr + val)
	}
}

// Read a Z-machine string literal: 2byte pairs from the current PC
// Returns the decoded string and number of words read
func (m *Machine) readStringLiteral(addr uint16) (string, int) {
	words := []uint16{}
	for i := uint16(0); int(i) < len(m.mem); i += 2 {
		word := decode.GetWord(m.mem, addr+i)
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

	// Branch is taken
	if condition == branchOnTrue {
		// If offset is 0 or 1, this is a special case meaning return false or true
		// GOTCHA: It only applies if the branch would be taken!
		switch offset {
		case 0:
			m.debug("   -> branch offset is 0, returning false\n")
			m.returnFromCall(0)
			return
		case 1:
			m.debug("   -> branch offset is 1, returning true\n")
			m.returnFromCall(1)
			return
		}

		m.pc = uint16(int16(m.pc) + int16(instLen) + branchDataLen + offset - 2)
		m.debug("   -> branching to %04x\n", m.pc)
	} else {
		// Branch not taken, continue to next instruction
		m.pc += instLen + uint16(branchDataLen)
		m.debug("   -> no branch, next pc %04x\n", m.pc)
	}
}

// validateChecksum computes and validates the checksum of the loaded Z-machine file
func (m *Machine) validateChecksum() bool {
	checksum := uint16(0)
	// Calculate checksum from byte 0x40 (64) onwards to end of file
	for i := 0x40; i < len(m.mem); i++ {
		checksum = (checksum + uint16(m.mem[i])&0xFFFF)
	}
	return checksum == m.checksum
}

func (m *Machine) print(s string) {
	if m.outputStream == OUTPUT_STREAM_SCREEN {
		m.ext.TextOut(s)
	} else if m.outputStream == OUTPUT_STREAM_MEMORY {
		// NOT IMPLEMENTED
	}
}

func (m *Machine) readString() string {
	if m.inputStream == INPUT_STREAM_KEYBOARD {
		return m.ext.ReadInput()
	} else if m.inputStream == INPUT_STREAM_FILE {
		panic("NOT_IMPLEMENTED: input stream from file")
	}
	return ""
}

func (m *Machine) tokenizeInput(input string, parseAddr uint16) {
	// Simple tokenizer: split on spaces, no dictionary lookup
	words := strings.Fields(input)

	// First byte of parse buffer is max words
	maxWords := m.mem[parseAddr]
	numWords := byte(len(words))
	if numWords > maxWords {
		numWords = maxWords
	}
	m.mem[parseAddr+1] = numWords

	// Each word entry is 4 bytes: 1 byte length, 1 byte position, 2 bytes dictionary index (we set to 0)
	currPos := byte(0)
	for i := byte(0); i < numWords; i++ {
		word := words[i]
		wordLen := byte(len(word))
		if wordLen > 15 {
			wordLen = 15 // truncate to 15 chars
		}

		// Write length
		m.mem[parseAddr+2+uint16(i*4)] = wordLen
		// Write position
		m.mem[parseAddr+2+uint16(i*4)+1] = currPos
		// Write dictionary index (2 bytes) as 0 for now
		decode.SetWord(m.mem, parseAddr+2+uint16(i*4)+2, 0)

		currPos += wordLen + 1 // +1 for space
	}
}
