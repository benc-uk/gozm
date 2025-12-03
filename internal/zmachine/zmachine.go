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
	"os"

	"github.com/benc-uk/gozm/internal/decode"
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
	pc           uint16 // Actually think this probably needs to be a uint32
	callStack    []callFrame
	debugLevel   int
	ext          External
	propDefaults []uint16
	objects      []*zObject

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
		m.step()
	}
}

// step executes a single instruction at the current program counter
func (m *Machine) step() {
	inst := m.decodeInst()
	m.debug("\n%04X: %s\n", m.pc, inst.String())

	// Decode and execute instructions!
	switch inst.code {
	// ===================== 0OP INSTRUCTIONS =====================
	// NOP
	case 0x00, 0x20, 0x40, 0x60, 0xC0:
		m.pc += inst.len

	// RET_TRUE
	case 0xB0:
		m.returnFromCall(1)

	// RET_FALSE
	case 0xB1:
		m.returnFromCall(0)

	// PRINT (literal string)
	case 0xB2:
		str, wordCount := m.readStringLiteral(m.pc + 1)
		m.ext.TextOut(str)
		m.pc += uint16(wordCount*2) + 1 // Advance PC past the string

	// PRINT_RET (literal string)
	case 0xB3:
		str, _ := m.readStringLiteral(m.pc + 1)
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
		score := m.getVar(17) // global variable 17 is score
		turns := m.getVar(16) // global variable 16 is turns
		// TODO: Placeholder for scoreboard/status line handling, use ansi code to invert colors
		m.ext.TextOut(fmt.Sprintf("\n\033[32m\033[7m Unknown location                  score:%d turns:%d \033[27m\033[0m\n", score, turns))
		m.pc += inst.len

	// VERIFY
	case 0xBD:
		m.pc += inst.len // Not worth implementing, just skip

	// RET_POPPED
	case 0xB8:
		val := m.getCallFrame().Pop()
		m.returnFromCall(val)

	// POP
	case 0xB9:
		m.getCallFrame().Pop()
		m.pc += inst.len

	// ===================== 1OP INSTRUCTIONS =====================

	// JZ
	case 0x80, 0x90, 0xA0:
		val := inst.operands[0]
		m.branchHandler(inst.len, val == 0)

	// GET_SIBLING
	case 0x81, 0x91, 0xA1:
		objNum := byte(inst.operands[0])
		sibling := m.getObject(objNum).sibling
		dest := m.mem[m.pc+inst.len] // destination in next byte
		m.storeVar(uint16(dest), uint16(sibling))
		m.branchHandler(inst.len+1, sibling != NULL_OBJECT)

	// GET_CHILD
	case 0x82, 0x92, 0xA2:
		objNum := byte(inst.operands[0])
		sibling := m.getObject(objNum).child
		dest := m.mem[m.pc+inst.len] // destination in next byte
		m.storeVar(uint16(dest), uint16(sibling))
		m.branchHandler(inst.len+1, sibling != NULL_OBJECT)

	// GET_PARENT
	case 0x83, 0x93, 0xA3:
		objNum := byte(inst.operands[0])
		sibling := m.getObject(objNum).parent
		dest := m.mem[m.pc+inst.len] // destination in next byte
		m.storeVar(uint16(dest), uint16(sibling))
		m.pc += inst.len + 1 // +1 for dest byte

	// GET_PROP_LEN
	case 0x84, 0x94, 0xA4:
		propAddr := inst.operands[0]
		var length byte
		if propAddr == 0 {
			length = 0
		} else {
			// Gotcha: The property address points to the property data, not the size byte
			// The size byte is immediately before the property data
			sizeByte := m.mem[propAddr-1]
			_, length = decode.PropSizeNumber(sizeByte)
		}
		dest := m.mem[m.pc+inst.len] // destination in next byte
		m.storeVar(uint16(dest), uint16(length))
		m.pc += inst.len + 1 // +1 for dest byte

	// INC
	case 0x85, 0x95, 0xA5:
		varLoc := inst.operands[0]
		m.addToVar(varLoc, 1)
		m.pc += inst.len

	// DEC
	case 0x86, 0x96, 0xA6:
		varLoc := inst.operands[0]
		m.addToVar(varLoc, -1)
		m.pc += inst.len

	// PRINT_ADDR
	case 0x87, 0x97, 0xA7:
		addr := inst.operands[0]
		m.trace(" - print_addr from %04x\n", addr)
		str, _ := m.readStringLiteral(addr)
		m.ext.TextOut(str)
		m.pc += inst.len

	// REMOVE_OBJ
	case 0x89, 0x99, 0xA9:
		objNum := byte(inst.operands[0])
		m.getObject(objNum).removeObjectFromParent(m)
		m.pc += inst.len

	// PRINT_OBJ
	case 0x8A, 0x9A, 0xAA:
		objNum := byte(inst.operands[0])
		obj := m.getObject(objNum)
		m.ext.TextOut(obj.description)
		m.pc += inst.len

	// RET
	case 0x8B, 0x9B, 0xAB:
		val := inst.operands[0]
		m.returnFromCall(val)

	// JUMP
	case 0x8C, 0x9C, 0xAC:
		offset := decode.Convert14BitToSigned(inst.operands[0])
		m.pc = uint16(int16(m.pc) + int16(inst.len) + offset - 2)

	// PRINT_PADDR
	case 0x8D, 0x9D, 0xAD:
		packedAddr := inst.operands[0]
		addr := decode.PackedAddress(packedAddr)
		m.trace(" - print_paddr from %04x\n", addr)
		str, _ := m.readStringLiteral(addr)
		m.ext.TextOut(str)
		m.pc += inst.len

	// LOAD
	case 0x8E, 0x9E, 0xAE:
		opVal := inst.operands[0]
		actualVal := m.getVar(uint16(opVal)) // This seems wrong!
		varLoc := m.mem[m.pc+inst.len]       // destination in next byte
		m.storeVar(uint16(varLoc), actualVal)
		m.pc += inst.len + 1 // +1 for dest byte

	// NOT (BITWISE)
	case 0x8F, 0x9F, 0xAF:
		v := inst.operands[0]
		varLoc := m.mem[m.pc+inst.len]
		m.storeVar(uint16(varLoc), ^v)
		m.pc += inst.len + 1 // +1 for dest byte

	// ===================== 2OP INSTRUCTIONS =====================

	// JE
	case 0x01, 0x21, 0x41, 0x61, 0xC1:
		condition := false
		firstVal := inst.operands[0]
		for _, val := range inst.operands[1:] {
			if firstVal == val {
				condition = true
				break
			}
		}
		m.branchHandler(inst.len, condition)

	// JL
	case 0x02, 0x22, 0x42, 0x62, 0xC2:
		v1 := int16(inst.operands[0])
		v2 := int16(inst.operands[1])
		m.branchHandler(inst.len, v1 < v2)

	// JG
	case 0x03, 0x23, 0x43, 0x63, 0xC3:
		v1 := int16(inst.operands[0])
		v2 := int16(inst.operands[1])
		m.branchHandler(inst.len, v1 > v2)

	// DEC_CHK
	case 0x04, 0x24, 0x44, 0x64, 0xC4:
		varLoc := inst.operands[0]
		compareVal := int16(inst.operands[1])
		newVal := m.addToVar(varLoc, -1)
		// newVal := int16(m.getVar(varLoc))
		m.branchHandler(inst.len, newVal < compareVal)

	// INC_CHK
	case 0x05, 0x25, 0x45, 0x65, 0xC5:
		varLoc := inst.operands[0]
		compareVal := int16(inst.operands[1])
		newVal := m.addToVar(varLoc, 1)
		//newVal := int16(m.getVar(varLoc))
		m.branchHandler(inst.len, newVal > compareVal)

	// JIN
	case 0x06, 0x26, 0x46, 0x66, 0xC6:
		childObjNum := byte(inst.operands[0])
		parentObjNum := byte(inst.operands[1])
		childObj := m.getObject(childObjNum)
		m.branchHandler(inst.len, childObj.parent == parentObjNum)

	// OR (BITWISE)
	case 0x08, 0x28, 0x48, 0x68, 0xC8:
		v1 := inst.operands[0]
		v2 := inst.operands[1]

		dest := m.mem[m.pc+inst.len] // destination in next byte
		m.debug(" - or dest var:%d\n", dest)
		m.storeVar(uint16(dest), v1|v2)
		m.pc += inst.len + 1 // +1 for dest byte

	// AND (BITWISE)
	case 0x09, 0x29, 0x49, 0x69, 0xC9:
		v1 := inst.operands[0]
		v2 := inst.operands[1]
		dest := m.mem[m.pc+inst.len] // destination in next byte
		m.debug(" - and dest var:%d\n", dest)
		m.storeVar(uint16(dest), v1&v2)
		m.pc += inst.len + 1 // +1 for dest byte

	// TEST_ATTR
	case 0x0A, 0x2A, 0x4A, 0x6A, 0xCA:
		objNum := byte(inst.operands[0])
		attrNum := byte(inst.operands[1])
		obj := m.getObject(objNum)
		m.branchHandler(inst.len, obj.hasAttribute(attrNum))

	// SET_ATTR
	case 0x0B, 0x2B, 0x4B, 0x6B, 0xCB:
		objNum := byte(inst.operands[0])
		attrNum := byte(inst.operands[1])
		obj := m.getObject(objNum)
		obj.setAttribute(attrNum, true)
		m.pc += inst.len

	// CLEAR_ATTR
	case 0x0C, 0x2C, 0x4C, 0x6C, 0xCC:
		objNum := byte(inst.operands[0])
		attrNum := byte(inst.operands[1])
		obj := m.getObject(objNum)
		obj.setAttribute(attrNum, false)
		m.pc += inst.len

	// STORE
	case 0x0D, 0x2D, 0x4D, 0x6D, 0xCD:
		v := inst.operands[0]
		s := inst.operands[1]

		m.storeVar(v, s)
		m.pc += inst.len

	// INSERT_OBJ
	case 0x0E, 0x2E, 0x4E, 0x6E, 0xCE:
		objNum := byte(inst.operands[0])
		destParentNum := byte(inst.operands[1])
		obj := m.getObject(objNum)
		obj.insertIntoParent(m, destParentNum)
		m.pc += inst.len

	// GET_PROP
	case 0x11, 0x31, 0x51, 0x71, 0xD1:
		objNum := byte(inst.operands[0])
		propNum := byte(inst.operands[1])
		obj := m.getObject(objNum)
		val := obj.getPropertyValue(propNum, m.propDefaults)
		dest := m.mem[m.pc+inst.len] // destination in next byte
		m.storeVar(uint16(dest), val)
		m.pc += inst.len + 1 // +1 for dest byte

	// GET_PROP_ADDR
	case 0x12, 0x32, 0x52, 0x72, 0xD2:
		objNum := byte(inst.operands[0])
		propNum := byte(inst.operands[1])
		dest := m.mem[m.pc+inst.len] // destination in next byte
		obj := m.getObject(objNum)
		m.storeVar(uint16(dest), uint16(obj.propMap[propNum].addr))
		m.pc += inst.len + 1 // +1 for dest byte

	// GET_NEXT_PROP
	case 0x13, 0x33, 0x53, 0x73, 0xD3:
		objNum := byte(inst.operands[0])
		propNum := byte(inst.operands[1])
		obj := m.getObject(objNum)
		var nextPropNum byte

		if propNum == 0 {
			// Return first property number
			if len(obj.properties) > 0 {
				nextPropNum = obj.properties[0].num
			}
		} else {
			for i, prop := range obj.properties {
				if prop.num == propNum {
					if i+1 < len(obj.properties) {
						nextPropNum = obj.properties[i+1].num
					} else {
						nextPropNum = 0
					}
					break
				}
			}
		}

		dest := m.mem[m.pc+inst.len] // destination in next byte
		m.storeVar(uint16(dest), uint16(nextPropNum))
		m.pc += inst.len + 1 // +1 for dest byte

	// ADD
	case 0x14, 0x34, 0x54, 0x74, 0xD4:
		v := inst.operands[0]
		s := inst.operands[1]
		dest := m.mem[m.pc+inst.len] // destination in next byte
		m.debug(" - add dest var:%d\n", dest)

		m.storeVar(uint16(dest), v+s)
		m.pc += inst.len + 1 // +1 for dest byte

	// SUB
	case 0x15, 0x35, 0x55, 0x75, 0xD5:
		v := inst.operands[0]
		s := inst.operands[1]
		dest := m.mem[m.pc+inst.len] // destination in next byte
		m.debug(" - sub dest var:%d\n", dest)

		m.storeVar(uint16(dest), v-s)
		m.pc += inst.len + 1 // +1 for dest byte

	// MUL
	case 0x16, 0x36, 0x56, 0x76, 0xD6:
		v := inst.operands[0]
		s := inst.operands[1]
		dest := m.mem[m.pc+inst.len] // destination in next byte
		m.debug(" - mul dest var:%d\n", dest)

		m.storeVar(uint16(dest), v*s)
		m.pc += inst.len + 1 // +1 for dest byte

	// DIV
	case 0x17, 0x37, 0x57, 0x77, 0xD7:
		v := inst.operands[0]
		s := inst.operands[1]
		dest := m.mem[m.pc+inst.len] // destination in next byte
		m.debug(" - div dest var:%d\n", dest)

		if s == 0 {
			panic("Division by zero!")
		}

		// NOTE: div should be signed division
		m.storeVar(uint16(dest), uint16(int16(v)/int16(s)))
		m.pc += inst.len + 1 // +1 for dest byte

	// MOD
	case 0x18, 0x38, 0x58, 0x78, 0xD8:
		v := inst.operands[0]
		s := inst.operands[1]
		dest := m.mem[m.pc+inst.len] // destination in next byte
		m.debug(" - mod dest var:%d\n", dest)

		if s == 0 {
			panic("Division by zero!")
		}

		// NOTE: mod should be signed modulus
		m.storeVar(uint16(dest), uint16(int16(v)%int16(s)))
		m.pc += inst.len + 1 // +1 for dest byte

	// LOADW (read word from array)
	case 0x0F, 0x2F, 0x4F, 0x6F, 0xCF:
		arrayAddr := inst.operands[0]
		index := inst.operands[1]
		wordAddr := arrayAddr + uint16(index*2)
		val := decode.GetWord(m.mem, wordAddr)
		dest := m.mem[m.pc+inst.len] // destination in next byte
		m.debug(" - loadw from %04x dest var:%d\n", wordAddr, dest)

		m.storeVar(uint16(dest), val)
		m.pc += inst.len + 1 // +1 for dest byte

	// LOADB (read byte from array)
	case 0x10, 0x30, 0x50, 0x70, 0xD0:
		arrayAddr := inst.operands[0]
		index := inst.operands[1]
		byteAddr := arrayAddr + uint16(index)
		val := m.mem[byteAddr]
		dest := m.mem[m.pc+inst.len] // destination in next byte
		m.debug(" - loadb from %04x dest var:%d\n", byteAddr, dest)

		m.storeVar(uint16(dest), uint16(val))
		m.pc += inst.len + 1 // +1 for dest byte

	// ===================== VAR INSTRUCTIONS =====================

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

	// STOREW
	case 0xE1:
		arrayAddr := inst.operands[0]
		index := inst.operands[1]
		val := inst.operands[2]
		wordAddr := arrayAddr + uint16(index*2)
		m.debug(" - storew to %04x value:%04x\n", wordAddr, val)

		decode.SetWord(m.mem, wordAddr, val)
		m.pc += inst.len

	// STOREB
	case 0xE2:
		arrayAddr := inst.operands[0]
		index := inst.operands[1]
		val := byte(inst.operands[2])
		byteAddr := arrayAddr + uint16(index)
		m.debug(" - storeb to %04x value:%02x\n", byteAddr, val)

		m.mem[byteAddr] = val
		m.pc += inst.len

	// PUT_PROP
	case 0xE3:
		objNum := byte(inst.operands[0])
		propNum := byte(inst.operands[1])
		val := inst.operands[2]
		obj := m.getObject(objNum)
		obj.setPropertyValue(propNum, val)
		m.pc += inst.len

	// PRINT_NUM
	case 0xE6:
		v := inst.operands[0]
		m.ext.TextOut(fmt.Sprintf("%d", int16(v))) // Print as signed number
		m.pc += inst.len

	// PUSH
	case 0xE8:
		val := inst.operands[0]
		m.getCallFrame().Push(val)
		m.pc += inst.len

	// PULL
	case 0xE9:
		val := m.getCallFrame().Pop()
		varLoc := inst.operands[0]
		m.storeVar(varLoc, val)
		m.pc += inst.len

	// Unimplemented instruction!
	default:
		panic(fmt.Sprintf("\n💥 Unimplemented instruction: %02x", inst.code))
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

func (m *Machine) addToVar(loc uint16, val int16) int16 {
	// We made loc uint16 for ease of use, now restrict to valid range
	if loc > 0xFF {
		panic(fmt.Sprintf("Variable location out of range: %02x", loc))
	}

	cf := m.getCallFrame()
	if loc == 0 {
		// Stack variable
		curr := int16(cf.Pop())
		cf.Push(uint16(curr + val))
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

	// If offset is 0 or 1, this is a special case meaning return false or true
	switch offset {
	case 0:
		m.debug(" - branch offset is 0, returning false\n")
		m.returnFromCall(0)
		return
	case 1:
		m.debug(" - branch offset is 1, returning true\n")
		m.returnFromCall(1)
		return
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
