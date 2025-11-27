package zmachine

import "gozm/internal/decode"

// implementation of pcode 2B - Print literal string
type instruction struct {
	code     byte     // opcode byte
	operands []uint16 // operands for the instruction
	byteLen  uint16   // total length of instruction + operands in bytes
}

// Decodes the instruction at the current program counter
// Following v3 spec
func (m *Machine) decodeInst() instruction {
	inst := instruction{}

	inst.code = m.mem[m.pc]
	inst.byteLen = 1

	// 0OP is just the opcode, no operands, B0-BF
	if inst.code&0xF0 == 0xB0 {
		return inst
	}

	// TODO: handle other instruction formats
	// 1OP
	// 2OP

	// VAR
	operandTypes := m.mem[m.pc+1]

	// Decode the operands based on operand types
	inst.byteLen++
	for i := uint16(0); i < 2; i++ {
		shift := 6 - (i * 2)
		otype := (operandTypes >> shift) & 0x03

		switch otype {
		case 0x00: // large constant
			val := decode.GetWord(m.mem, m.pc+2+i*2) // TODO: Possible bug if mixed operand types
			inst.operands = append(inst.operands, val)
			inst.byteLen += 2
		case 0x01: // small constant
			val := uint16(m.mem[m.pc+2+i])
			inst.operands = append(inst.operands, val)
			inst.byteLen += 1
		case 0x02: // variable
			// For now, just read the variable number, not the value
			// TODO: implement variable reading
			varNum := m.mem[m.pc+2+i]
			inst.operands = append(inst.operands, uint16(varNum))
			inst.byteLen += 1
		case 0x03: // omitted
			// No operand
		}
	}

	return inst
}
