// =======================================================================
// Package: zmachine - Core Z-machine interpreter
// instructions.go - Instruction decoding and representation
//
// Copyright (c) 2025 Ben Coleman. Licensed under the MIT License
// =======================================================================

package zmachine

import (
	"fmt"

	"github.com/benc-uk/gozm/internal/decode"
)

// Enum for forms of instructions
const (
	FORM_LONG  = 0
	FORM_SHORT = 1
	FORM_VAR   = 2
)

const MAX_OPERANDS = 4

const OPTYPE_LARGE_CONST = 0x00
const OPTYPE_SMALL_CONST = 0x01
const OPTYPE_VARIABLE = 0x02
const OPTYPE_OMITTED = 0x03

// instruction represents a decoded Z-machine instruction
type instruction struct {
	code byte // opcode byte
	//number   byte     // instruction number, not used yet
	operands []uint16 // operand values for the instruction
	len      uint16   // total length of instruction + operands in bytes
	form     byte     // form of the instruction, not used yet
}

// Decodes the instruction at the current program counter
func (m *Machine) decodeInst() instruction {
	inst := instruction{
		code: m.mem[m.pc],
		len:  1, // start with 1 for the opcode byte
	}

	// Debug ops if they are being traced
	for _, op := range m.TracedOps {
		if m.mem[m.pc] == op {
			m.debugLevel = DEBUG_TRACE
		}
	}

	if inst.code == 0 {
		return inst
	}

	// VAR form has $11 in the top bits, and a following operand types byte
	if inst.code&0xC0 == 0xC0 {
		inst.form = FORM_VAR
		opTypesByte := m.mem[m.pc+1]
		inst.len++ // for operand types byte

		// Decode the operands, there's a max of 4 operand types held in 1 byte
		// Each operand type is represented by 2 bits
		operandPtr := m.pc + 2 // start of operands, after op types byte
		shift := uint8(6)
		for i := 0; i < MAX_OPERANDS; i++ {
			opType := (opTypesByte >> shift) & 0x3
			shift -= 2

			if opType == OPTYPE_OMITTED {
				break
			}

			val, opLen := fetchOperand(m, opType, operandPtr)
			inst.operands = append(inst.operands, val)
			inst.len += opLen
			operandPtr += opLen
		}

		m.trace("Decode var: %02x typeByte:%02x\n", inst.code, opTypesByte)

		return inst
	}

	// SHORT form has $10 in the top bits
	if inst.code&0xC0 == 0x80 {
		inst.form = FORM_SHORT
		// Get bits 4 and 5 for operand type
		opType := (inst.code >> 4) & 0x3
		m.trace("Decode short: %02x opType:%02x\n", inst.code, opType)

		if opType == OPTYPE_OMITTED {
			return inst // No operands, this is a 0OP instruction
		}

		val, len := fetchOperand(m, opType, m.pc+1)
		inst.operands = []uint16{val}
		inst.len += len

		return inst
	}

	// LONG form otherwise, this form is always 2OP
	// https://zspec.jaredreisinger.com/04-instructions#4_3
	inst.form = FORM_LONG
	// Value of bits 6 & 5 indicates types of 2 operands
	// GOTCHA: Horrible - https://zspec.jaredreisinger.com/04-instructions#4_4_2
	op1Type := (inst.code>>6)&0x1 + 1 // +1 to map 0->1, 1->2
	op2Type := (inst.code>>5)&0x1 + 1

	m.trace("Decode long: %02x opType1:%d opType2:%d\n", inst.code, op1Type, op2Type)

	op1, _ := fetchOperand(m, op1Type, m.pc+1)
	op2, _ := fetchOperand(m, op2Type, m.pc+2)

	inst.len += 2 // for the two operands
	inst.operands = []uint16{op1, op2}

	return inst
}

// Helper to fetch an operand based on its type, returning the value and length in bytes
func fetchOperand(m *Machine, operandType byte, loc uint16) (uint16, uint16) {
	switch operandType {
	case OPTYPE_LARGE_CONST: // large constant
		val := decode.GetWord(m.mem, loc)
		return val, 2
	case OPTYPE_SMALL_CONST: // small constant
		val := uint16(m.mem[loc])
		return val, 1
	case OPTYPE_VARIABLE: // variable
		val := m.getVar(uint16(m.mem[loc]))
		return val, 1
	case OPTYPE_OMITTED: // omitted, should not happen here
		return 0, 0
	default:
		panic(fmt.Sprintf("Invalid operand type: %02x", operandType))
	}
}

// String representation of the instruction
func (inst *instruction) String() string {
	return fmt.Sprintf("%s (code=%02x, operands=%v, len=%d, form=%d)", opcodeNames[inst.code], inst.code, inst.operands, inst.len, inst.form)
}

// Map opcodes to names
// opcodeNames maps opcode byte values to instruction names for Z-machine versions 1-3.
// Only instructions valid in versions 1,2,3 are included. Operand-type variants
// of 1OP and 2OP instructions (small const, large const, variable) are all mapped
// individually to the same mnemonic.
var opcodeNames = map[byte]string{
	// 0OP (short form with operand type = omitted) B0-BD, excluding extended (BE) and piracy (BF v5+)
	0xB0: "rtrue",
	0xB1: "rfalse",
	0xB2: "print",
	0xB3: "print_ret",
	0xB4: "nop", // version 1 (ignored later)
	0xB5: "save",
	0xB6: "restore",
	0xB7: "restart",
	0xB8: "ret_popped",
	0xB9: "pop",
	0xBA: "quit",
	0xBB: "new_line",
	0xBC: "show_status", // scoreboard/status line handling
	0xBD: "verify",

	// 1OP (short form with one operand) Large const (80-8F), Small const (90-9F), Variable (A0-AF)
	// Exclude call_1s (8, v4+) and call_1n (15 in v5+) variants.
	0x80: "jz", 0x90: "jz", 0xA0: "jz",
	0x81: "get_sibling", 0x91: "get_sibling", 0xA1: "get_sibling",
	0x82: "get_child", 0x92: "get_child", 0xA2: "get_child",
	0x83: "get_parent", 0x93: "get_parent", 0xA3: "get_parent",
	0x84: "get_prop_len", 0x94: "get_prop_len", 0xA4: "get_prop_len",
	0x85: "inc", 0x95: "inc", 0xA5: "inc",
	0x86: "dec", 0x96: "dec", 0xA6: "dec",
	0x87: "print_addr", 0x97: "print_addr", 0xA7: "print_addr",
	// 0x88/0x98/0xA8 call_1s (v4) skipped
	0x89: "remove_obj", 0x99: "remove_obj", 0xA9: "remove_obj",
	0x8A: "print_obj", 0x9A: "print_obj", 0xAA: "print_obj",
	0x8B: "ret", 0x9B: "ret", 0xAB: "ret",
	0x8C: "jump", 0x9C: "jump", 0xAC: "jump",
	0x8D: "print_paddr", 0x9D: "print_paddr", 0xAD: "print_paddr",
	0x8E: "load", 0x9E: "load", 0xAE: "load",
	// 0x8F/0x9F/0xAF: 'not' kept (present in early versions); later repurposed but still valid name here
	0x8F: "not", 0x9F: "not", 0xAF: "not",

	// 2OP (long form) operand-type variants: base (00-1F), +0x20, +0x40, +0x60
	// Include only instruction numbers 1-24 (je..mod) valid in v1-3; skip >=25 which are v4+ (call_2*, set_colour, throw)
	// je (1)
	0x01: "je", 0x21: "je", 0x41: "je", 0x61: "je",
	// jl (2)
	0x02: "jl", 0x22: "jl", 0x42: "jl", 0x62: "jl",
	// jg (3)
	0x03: "jg", 0x23: "jg", 0x43: "jg", 0x63: "jg",
	// dec_chk (4)
	0x04: "dec_chk", 0x24: "dec_chk", 0x44: "dec_chk", 0x64: "dec_chk",
	// inc_chk (5)
	0x05: "inc_chk", 0x25: "inc_chk", 0x45: "inc_chk", 0x65: "inc_chk",
	// jin (6)
	0x06: "jin", 0x26: "jin", 0x46: "jin", 0x66: "jin",
	// test (bitmap flags) (7)
	0x07: "test", 0x27: "test", 0x47: "test", 0x67: "test",
	// or (8)
	0x08: "or", 0x28: "or", 0x48: "or", 0x68: "or",
	// and (9)
	0x09: "and", 0x29: "and", 0x49: "and", 0x69: "and",
	// test_attr (10)
	0x0A: "test_attr", 0x2A: "test_attr", 0x4A: "test_attr", 0x6A: "test_attr",
	// set_attr (11)
	0x0B: "set_attr", 0x2B: "set_attr", 0x4B: "set_attr", 0x6B: "set_attr",
	// clear_attr (12)
	0x0C: "clear_attr", 0x2C: "clear_attr", 0x4C: "clear_attr", 0x6C: "clear_attr",
	// store (13)
	0x0D: "store", 0x2D: "store", 0x4D: "store", 0x6D: "store",
	// insert_obj (14)
	0x0E: "insert_obj", 0x2E: "insert_obj", 0x4E: "insert_obj", 0x6E: "insert_obj",
	// loadw (15)
	0x0F: "loadw", 0x2F: "loadw", 0x4F: "loadw", 0x6F: "loadw",
	// loadb (16)
	0x10: "loadb", 0x30: "loadb", 0x50: "loadb", 0x70: "loadb",
	// get_prop (17)
	0x11: "get_prop", 0x31: "get_prop", 0x51: "get_prop", 0x71: "get_prop",
	// get_prop_addr (18)
	0x12: "get_prop_addr", 0x32: "get_prop_addr", 0x52: "get_prop_addr", 0x72: "get_prop_addr",
	// get_next_prop (19)
	0x13: "get_next_prop", 0x33: "get_next_prop", 0x53: "get_next_prop", 0x73: "get_next_prop",
	// add (20)
	0x14: "add", 0x34: "add", 0x54: "add", 0x74: "add",
	// sub (21)
	0x15: "sub", 0x35: "sub", 0x55: "sub", 0x75: "sub",
	// mul (22)
	0x16: "mul", 0x36: "mul", 0x56: "mul", 0x76: "mul",
	// div (23)
	0x17: "div", 0x37: "div", 0x57: "div", 0x77: "div",
	// mod (24)
	0x18: "mod", 0x38: "mod", 0x58: "mod", 0x78: "mod",

	// VAR form (E0-EB) for versions 1-3
	0xE0: "call", // call with up to 3 args returning result
	0xE1: "storew",
	0xE2: "storeb",
	0xE3: "put_prop",
	0xE4: "sread",
	0xE5: "print_char",
	0xE6: "print_num",
	0xE7: "random",
	0xE8: "push",
	0xE9: "pull",
	0xEA: "split_window",
	0xEB: "set_window",
}
