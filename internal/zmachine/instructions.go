package zmachine

import (
	"fmt"
	"gozm/internal/decode"
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
	code     byte     // opcode byte
	operands []uint16 // operand values for the instruction
	len      uint16   // total length of instruction + operands in bytes
	form     byte     // form of the instruction, not used yet
}

// Decodes the instruction at the current program counter
func (m *Machine) decodeInst() instruction {
	inst := instruction{}

	inst.code = m.mem[m.pc]
	inst.len = 1
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
			m.trace("DECODE VAR for (%02x) type:%02x\n", inst.code, opType)

			if opType == OPTYPE_OMITTED {
				break
			}

			val, opLen := fetchOperand(m, opType, operandPtr)
			inst.operands = append(inst.operands, val)
			inst.len += opLen
			operandPtr += opLen
		}

		return inst
	}

	// SHORT form has $10 in the top bits
	if inst.code&0xC0 == 0x80 {
		inst.form = FORM_SHORT
		// Get bits 4 and 5 for operand type
		opType := (inst.code >> 4) & 0x3
		m.trace("DECODE SHORT for (%02x) type:%02x\n", inst.code, opType)

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

	m.trace("DECODE LONG for (%02x) type1:%d type2:%d\n", inst.code, op1Type, op2Type)

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
		val := m.GetVar(uint16(m.mem[loc]))
		return val, 1
	case OPTYPE_OMITTED: // omitted, should not happen here
		return 0, 0
	default:
		panic(fmt.Sprintf("Invalid operand type: %02x", operandType))
	}
}

// String representation of the instruction
func (inst *instruction) String() string {
	return "Code: " + fmt.Sprintf("%02x", inst.code) +
		", Length: " + fmt.Sprintf("%d", inst.len) +
		", Operands: " + fmt.Sprintf("%v", inst.operands)
}
