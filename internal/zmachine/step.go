// =======================================================================
// Package: zmachine - Core Z-machine interpreter
// step.go - step executes a single instruction at the current program counter
//
// Copyright (c) 2025 Ben Coleman. Licensed under the MIT License
// =======================================================================

package zmachine // step executes a single instruction at the current program counter

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strings"

	"github.com/benc-uk/gozm/internal/decode"
)

func (m *Machine) step() {
	m.debugLevel = DEBUG_NONE
	inst := m.decodeInst()

	// trap panic in case of errors
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("💥 Runtime error at %04X: %s\n", m.pc, inst.String())
			m.DumpMem(m.pc, 12)
			// Print stack trace
			fmt.Printf("Stack trace:\n")
			for i := len(m.callStack) - 1; i >= 0; i-- {
				frame := m.callStack[i]
				fmt.Printf(" - Frame %d: return to %04X\n", i, frame.returnAddr)
			}
			panic(r)
		}
	}()

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
		m.print(str)
		m.pc += uint16(wordCount*2) + 1 // Advance PC past the string

	// PRINT_RET (literal string)
	case 0xB3:
		str, _ := m.readStringLiteral(m.pc + 1)
		m.print(str + "\n")
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
		m.print("\n")
		m.pc += inst.len

	// SHOW_STATUS
	case 0xBC:
		score := m.getVar(17) // global variable 17 is score
		turns := m.getVar(16) // global variable 16 is turns
		// TODO: Placeholder for scoreboard/status line handling, use ansi code to invert colors
		m.print(fmt.Sprintf("\n\033[32m\033[7m Unknown location                  score:%d turns:%d \033[27m\033[0m\n", score, turns))
		m.pc += inst.len

	// VERIFY
	case 0xBD:
		res := true //m.validateChecksum()
		m.branchHandler(inst.len, res)

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
		m.print(str)
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
		m.print(obj.description)
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
		str, _ := m.readStringLiteral(addr)
		m.print(str)
		m.pc += inst.len

	// LOAD
	case 0x8E, 0x9E, 0xAE:
		opVal := inst.operands[0]
		var actualVal uint16
		if opVal == 0 {
			// Stack variable
			actualVal = m.getCallFrame().Peek()
		} else {
			actualVal = m.getVar(opVal)
		}

		varLoc := m.mem[m.pc+inst.len] // destination in next byte
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
		m.branchHandler(inst.len, newVal < compareVal)

	// INC_CHK
	case 0x05, 0x25, 0x45, 0x65, 0xC5:
		varLoc := inst.operands[0]
		compareVal := int16(inst.operands[1])
		newVal := m.addToVar(varLoc, 1)
		m.branchHandler(inst.len, newVal > compareVal)

	// JIN
	case 0x06, 0x26, 0x46, 0x66, 0xC6:
		childObjNum := byte(inst.operands[0])
		parentObjNum := byte(inst.operands[1])
		childObj := m.getObject(childObjNum)
		m.branchHandler(inst.len, childObj.parent == parentObjNum)

	// TEST
	case 0x07, 0x27, 0x47, 0x67, 0xC7:
		bitmap := inst.operands[0]
		flags := inst.operands[1]
		m.branchHandler(inst.len, (bitmap&flags) == flags)

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

		m.setVarInPlace(v, s)
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
		// Property address is address of property data, not header & it may not exist
		prop, exist := obj.propMap[propNum]
		addr := uint16(0)
		if exist {
			addr = prop.addr
		}
		m.storeVar(uint16(dest), uint16(addr))
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
		byteAddr := uint16(arrayAddr) + uint16(index)
		val := m.mem[byteAddr]
		dest := m.mem[m.pc+inst.len] // destination in next byte
		m.debug(" - loadb from %04x dest var:%d\n", byteAddr, dest)

		m.storeVar(uint16(dest), uint16(val))
		m.pc += inst.len + 1 // +1 for dest byte

	// ===================== VAR INSTRUCTIONS =====================

	// CALL
	case 0xE0:
		routineAddr := decode.PackedAddress(inst.operands[0])

		// When the address 0 is called as a routine, nothing happens and the return value is false
		if routineAddr == 0 {
			m.debug(" - call to NULL routine, returning false\n")
			// Store byte is at PC + inst.len, store 0 (false) there
			dest := m.mem[m.pc+inst.len]
			m.storeVar(uint16(dest), 0)
			m.pc += inst.len + 1 // +1 for store byte
			return
		}

		// Count locals from routine header
		numLocals := m.mem[routineAddr]
		if numLocals > 15 {
			panic(fmt.Sprintf("Attempted to call routine at %04x with invalid local count %d", routineAddr, numLocals))
		}

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

	// SREAD aka READ in v3
	case 0xE4:
		textAddr := inst.operands[0]
		parseAddr := inst.operands[1]
		maxLen := m.mem[textAddr]
		if maxLen == 0 {
			panic("READ called with zero max length")
		}
		maxLen-- // Weirdly, the first byte is the max length, so reduce by 1 for actual input

		// Read input from user
		input := m.readString()
		input = strings.ToLower(input)
		input = strings.Trim(input, "\r\n")

		// Copy input into memory, and null terminate, important!
		copy(m.mem[textAddr+1:textAddr+uint16(maxLen)], input)
		m.mem[textAddr+1+uint16(len(input))] = 0

		// Tokenize the input string by dictionary separators
		var tokens []string
		currentWord := ""

		for _, char := range input {
			// Check if this character is a separator
			isSep := false
			for _, sep := range m.dictSep {
				if char == rune(sep[0]) {
					isSep = true
					break
				}
			}

			if isSep {
				// Add current word if not empty
				if currentWord != "" {
					tokens = append(tokens, currentWord)
					currentWord = ""
				}
				// Add separator as a token (but not spaces)
				if char != ' ' {
					tokens = append(tokens, string(char))
				}
			} else if char == ' ' {
				// Space ends a word but is not added as a token
				if currentWord != "" {
					tokens = append(tokens, currentWord)
					currentWord = ""
				}
			} else {
				// Regular character, add to current word
				currentWord += string(char)
			}
		}

		// Don't forget the last word
		if currentWord != "" {
			tokens = append(tokens, currentWord)
		}

		// Now we have tokens, look them up in the dictionary

		dictHits := []dictEntry{}
		for _, token := range tokens {
			dictHits = append(dictHits, m.lookupWordInDict(token))
		}
		//fmt.Printf(" - Dictionary results: %v\n", dictEntries)

		// Write token count to parse table, after max tokens byte
		m.mem[parseAddr+1] = byte(len(dictHits))

		// Write each dictionary entry to parse table, each is 4 bytes
		parseOffset := parseAddr + 2 // Skip max tokens & count bytes
		for i, dictHit := range dictHits {
			entryOffset := parseOffset + uint16(i*4)

			// Write dictionary address (2 bytes)
			decode.SetWord(m.mem, entryOffset, dictHit.address)

			// Write word length (1 byte)
			word := dictHit.word
			if dictHit.address == 0 {
				word = tokens[i] // Use original token if not found
			}
			m.mem[entryOffset+2] = byte(len(word))

			// Write position in text buffer (1 byte), 1-based index
			// Need to find the position of the word in the original input
			position := byte(0)
			searchOffset := uint16(1) // Start after max length byte
			for pos, _ := range input {
				if strings.HasPrefix(input[searchOffset-1:], word) {
					position = byte(pos + 1) // 1-based index
					break
				}
				searchOffset++
			}
			m.mem[entryOffset+3] = position
		}

		m.pc += inst.len

	// PRINT_CHAR
	case 0xE5:
		charCode := byte(inst.operands[0])
		r := decode.ZSCIIChar(charCode)
		m.print(string(r))
		m.pc += inst.len

	// PRINT_NUM
	case 0xE6:
		v := inst.operands[0]
		m.print(fmt.Sprintf("%d", int16(v))) // Print as signed number
		m.pc += inst.len

	// RANDOM
	case 0xE7:
		maxVal := int16(inst.operands[0])
		var result uint16
		if maxVal <= 0 {
			// Reseed random number generator
			m.rand = rand.New(rand.NewPCG(uint64(maxVal), 0))
			result = 0
		} else {
			result = uint16(m.rand.IntN(int(maxVal)))
		}
		dest := m.mem[m.pc+inst.len] // destination in next byte
		m.storeVar(uint16(dest), result)
		m.pc += inst.len + 1 // +1 for dest byte

	// PUSH
	case 0xE8:
		val := inst.operands[0]
		m.getCallFrame().Push(val)
		m.pc += inst.len

	// PULL
	case 0xE9:
		val := m.getCallFrame().Pop()
		varLoc := inst.operands[0]
		m.setVarInPlace(varLoc, val)
		m.pc += inst.len

	// OUTPUT_STREAM
	case 0xF3:
		panic("NOT_IMPLEMENTED: OUTPUT_STREAM")

	// Unimplemented instruction!
	default:
		panic(fmt.Sprintf("\n💥 Unimplemented instruction: %02x", inst.code))
	}
}
