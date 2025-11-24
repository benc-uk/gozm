package zmachine

import (
	"fmt"
	"gozm/internal/decode"
)

type Machine struct {
	mem []byte
	pc  int16
}

func NewMachine(data []byte) *Machine {
	memCopy := make([]byte, len(data))
	copy(memCopy, data)

	pc := decode.GetInt16(memCopy, 0x06)

	return &Machine{
		mem: memCopy,
		pc:  pc,
	}
}

func (m *Machine) Step() {
	opcode := m.mem[m.pc]
	fmt.Printf("PC: %04x, Opcode: %02x\n", m.pc, opcode)

	// For demonstration, just increment the PC by 1
	m.pc += 1
}

func (m *Machine) Run() {
	// debug: print initial state
	fmt.Printf("Z-machine version: %d\n", m.GetVersion())
	m.PrintState()

	//	for {
	m.Step()
	// 	// Break condition for demonstration purposes
	// 	if m.pc >= int16(len(m.mem)) {
	// 		break
	// 	}
	// }
}

func (m *Machine) GetVersion() byte {
	return m.mem[0]
}

func (m *Machine) PrintState() {
	fmt.Printf("PC: %04x\n", m.pc)
}
