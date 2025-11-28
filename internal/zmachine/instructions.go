package zmachine

import "fmt"

// Enum for forms of instructions
const (
	FORM_LONG  = 0
	FORM_SHORT = 1
	FORM_VAR   = 2
)

// instruction represents a decoded Z-machine instruction
type instruction struct {
	code     byte     // opcode byte
	operands []uint16 // operand values for the instruction
	len      uint16   // total length of instruction + operands in bytes
	form     byte     // form of the instruction, not used yet
}

// Decodes the instruction at the current program counter
// Following v3 spec
func (m *Machine) decodeInst() instruction {
	inst := instruction{}

	inst.code = m.mem[m.pc]
	inst.len = 1

	// VAR FORM has $11 in the top bits
	if inst.code&0xC0 == 0xC0 {
		// TODO: Implement VAR FORM decoding
		return inst
	}

	// SHORT FORM has $10 in the top bits
	if inst.code&0xC0 == 0x80 {
		// TODO: Implement SHORT FORM decoding
		return inst
	}

	// LONG FORM otherwise, see docs on why
	// https://zspec.jaredreisinger.com/04-instructions#4_3
	inst.form = FORM_LONG
	op1 := uint16(m.mem[m.pc+1]) // Both operands are single byte in long form
	op2 := uint16(m.mem[m.pc+2])
	inst.operands = []uint16{op1, op2}
	inst.len += 2

	return inst
}

// string representation of the instruction
func (inst *instruction) String() string {
	return "Code: " + fmt.Sprintf("%02x", inst.code) +
		", Length: " + fmt.Sprintf("%d", inst.len) +
		", Operands: " + fmt.Sprintf("%v", inst.operands)
}
