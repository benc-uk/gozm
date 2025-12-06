// =======================================================================
// Package: zmachine - Core Z-machine interpreter
// zmachine.go - Main code, structs and core execution loop
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
	mem           []byte      // Z-machine memory
	pc            uint16      // Actually think this probably needs to be a uint32
	callStack     []callFrame // Call stack of routines
	debugLevel    int         // Debug verbosity level
	ext           External    // External interface for I/O
	propDefaults  []uint16    // Property defaults table
	objects       []*zObject  // Objects table
	rand          *rand.Rand  // Random number generator
	outputStream  int         // Current output stream
	inputStream   int         // Current input stream
	abbr          []string    // Abbreviation table
	dict          []dictEntry // Dictionary entries
	dictSep       []string    // Dictionary separator characters
	dictStartAddr uint16      // Start address of dictionary entries

	version     byte   // Header: version number
	highAddr    uint16 // Header: high memory address
	initialPC   uint16 // Header: initial program counter
	dictAddr    uint16 // Header: dictionary table start address
	objectsAddr uint16 // Header: objects table address
	globalsAddr uint16 // Header: global variables table address
	abbrvAddr   uint16 // Header: abbreviation table address
	fileLen     uint16 // Header: file length in words
	checksum    uint16 // Header: checksum
	flagStatus  bool   // Header: Is the status line shown?
}

type dictEntry struct {
	word    string
	address uint16
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

		version:     data[0x00],
		highAddr:    decode.GetWord(data, 0x04),
		initialPC:   decode.GetWord(data, 0x06),
		dictAddr:    decode.GetWord(data, 0x08),
		objectsAddr: decode.GetWord(data, 0x0A),
		globalsAddr: decode.GetWord(data, 0x0C),
		abbrvAddr:   decode.GetWord(data, 0x18),
		fileLen:     decode.GetWord(data, 0x1A),
		checksum:    decode.GetWord(data, 0x1C),
	}

	// Decode flag byte at 0x01, and status line flag is bit 4
	m.flagStatus = (data[0x01] & 0x10) != 0

	// Initialize abbreviations from the abbreviation table
	m.abbr = make([]string, 96)
	for i := uint16(0); i < 96; i++ {
		// Abbreviation table contains word addresses, need to multiply by 2
		// See: https://zspec.jaredreisinger.com/01-memory-map#1_2_2
		abbrStringAddr := decode.GetWord(m.mem, m.abbrvAddr+i*2) * 2
		s, _ := m.readStringLiteral(uint32(abbrStringAddr))
		m.abbr[i] = s
	}

	// Initialize objects, property defaults table
	m.initObjects()

	// Dictionary initialization
	numSepBytes := m.mem[m.dictAddr]
	m.dictSep = make([]string, numSepBytes)
	for i := byte(0); i < numSepBytes; i++ {
		m.dictSep[i] = string(m.mem[m.dictAddr+1+uint16(i)])
	}
	entryLen := m.mem[m.dictAddr+1+uint16(numSepBytes)]
	numEntries := decode.GetWord(m.mem, m.dictAddr+2+uint16(numSepBytes))

	// Load dictionary entries
	m.dict = make([]dictEntry, numEntries)
	m.dictStartAddr = m.dictAddr + 2 + uint16(numSepBytes) + 2
	for i := uint16(0); i < numEntries; i++ {
		entryAddr := m.dictStartAddr + uint16(i)*uint16(entryLen)
		s, _ := m.readStringLiteral(uint32(entryAddr))
		m.dict[i] = dictEntry{
			word:    s,
			address: entryAddr,
		}
	}

	m.debug("Z-machine initialized...\nVersion: %d, Size: %d\n", data[0x00], len(data))
	m.debug(" - Checksum: %04X, valid:%t\n", m.checksum, m.validateChecksum())
	m.debug(" - Objects/rooms: %d\n", len(m.objects))
	m.debug(" - Dictionary: %d entries, %d separators\n", numEntries, numSepBytes)
	m.debug(" - Abbreviations loaded: %d\n", len(m.abbr))

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
// Note: This takes a uint32 address to allow for strings in high memory
func (m *Machine) readStringLiteral(addr uint32) (string, int) {
	words := []uint16{}
	for i := uint32(0); int(i) < len(m.mem); i += 2 {
		word := decode.GetWord32(m.mem, addr+i)
		words = append(words, word)

		// If the high bit is set, this is the end of the string
		if word&0x8000 != 0 {
			break
		}
	}

	return decode.String(words, m.abbr), len(words)
}

// This is a complex helper used by all branch instructions
// See: https://zspec.jaredreisinger.com/04-instructions#4_7
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

// lookupWordInDict searches the dictionary for a word and returns its address
// Returns a dictEntry with address 0 if the word is not found
// If multiple words match, returns the longest matching word
// See: https://zspec.jaredreisinger.com/13-dictionary
// TODO: This needs to be optimized for large dictionaries!
func (m *Machine) lookupWordInDict(word string) dictEntry {
	longestMatch := dictEntry{address: 0}

	// Normalize the input word to lowercase for comparison
	word = strings.ToLower(word)

	for _, entry := range m.dict {
		entryWord := strings.ToLower(entry.word)

		// Dictionary entries are truncated (e.g., "leafle" for "leaflet" in v3)
		// Check if the input word starts with the dictionary entry
		if strings.HasPrefix(word, entryWord) {
			// Keep the longest match
			if len(entryWord) > len(longestMatch.word) {
				longestMatch = entry
			}
		}
	}

	return longestMatch
}

// And unfinished show status line function
// Almost no Z-machine games use this, so I'm leaving it unfinished
func (m *Machine) showStatus() {
	if !m.flagStatus {
		return
	}
	m.print("\033[2J") // Clear screen

	score := m.getVar(17)  // global variable 17 is score
	turns := m.getVar(18)  // global variable 16 is turns
	objNum := m.getVar(16) // global variable 1 is the current object

	obj := m.getObject(byte(objNum))
	m.print(fmt.Sprintf("\n\033[32m\033[7m %s                             score:%d turns:%d \033[27m\033[0m\n", obj.description, score, turns))

}
