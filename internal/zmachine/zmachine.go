package zmachine

import (
	"fmt"
	"gozm/internal/decode"
	"os"
)

type Machine struct {
	mem   []byte
	pc    uint16
	stack []callFrame

	header header
}

type header struct {
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

type callFrame struct {
	returnAddr uint16
	locals     []uint16
}

func (sf *callFrame) Push(val uint16) {
	sf.locals = append(sf.locals, val)
}

func (sf *callFrame) Pop() uint16 {
	if len(sf.locals) == 0 {
		return 0
	}

	val := sf.locals[len(sf.locals)-1]
	sf.locals = sf.locals[:len(sf.locals)-1]
	return val
}

func NewMachine(data []byte) *Machine {
	// Copy data into machine memory, very likely this is unnecessary
	memCopy := make([]byte, len(data))
	copy(memCopy, data)

	h := header{
		version:     memCopy[0x00],
		highAddr:    decode.GetWord(memCopy, 0x04),
		initialPC:   decode.GetWord(memCopy, 0x06),
		dictAddr:    decode.GetWord(memCopy, 0x08),
		objectsAddr: decode.GetWord(memCopy, 0x0A),
		globalsAddr: decode.GetWord(memCopy, 0x0C),
		staticAddr:  decode.GetWord(memCopy, 0x0E),
		abbrvAddr:   decode.GetWord(memCopy, 0x18),
		fileLen:     decode.GetWord(memCopy, 0x1A),
		checksum:    decode.GetWord(memCopy, 0x1C),
	}

	fmt.Println("Z-machine initialized & loaded...")
	fmt.Printf("Version: %d, PC: %04x, HiMem: %04x, StaticMem: %04x\n", h.version, h.initialPC, h.highAddr, h.staticAddr)

	return &Machine{
		mem:    memCopy,
		pc:     h.initialPC,
		header: h,
		stack:  make([]callFrame, 0),
	}
}

func (m *Machine) Step() {
	inst := m.decodeInst()
	fmt.Printf("PC: %04x, Code: %02x, Len: %d\n", m.pc, inst.code, inst.byteLen)

	// Decode and execute instruction here
	// This is a ultra minimum implementation with only a few opcodes
	switch inst.code {
	case 0x00:
		fmt.Println(" ++ NOP")
		m.pc += 1
	case 0xE0:
		fmt.Print(" ++ CALL: \x1b[32m")

		inst := m.decodeInst()
		funcAddr := decode.PackedAddress(inst.operands[0])
		fmt.Printf("%04x\x1b[0m\n", funcAddr)

		// Push new stack frame
		frame := callFrame{
			returnAddr: m.pc + uint16(inst.byteLen),
			// TODO: initialize locals from operands
			locals: make([]uint16, 0),
		}
		m.stack = append(m.stack, frame)

		// TODO: moving past routine header
		m.pc = funcAddr

	case 0xB0:

		// Pop stack frame
		if len(m.stack) == 0 {
			fmt.Println(" !! Stack underflow on RET")
			os.Exit(1)
		}

		frame := m.stack[len(m.stack)-1]
		m.stack = m.stack[:len(m.stack)-1]

		m.pc = frame.returnAddr
		storeLoc := m.mem[m.pc]
		fmt.Printf(" ++ RET_TRUE into: %0x\n", storeLoc)
		m.Store(uint16(storeLoc), 1) // true
		m.pc += 1
	case 0xBA:
		fmt.Println(" ++ QUIT")
		os.Exit(0)
	case 0xB2:
		fmt.Print(" ++ PRINT: \x1b[35m")
		words := []uint16{}
		for i := uint16(0); int(i) < len(m.mem); i += 2 {
			word := decode.GetWord(m.mem, m.pc+1+i)
			words = append(words, word)

			// If the high bit is set, this is the end of the string
			if word&0x8000 != 0 {
				break
			}
		}
		str := decode.String(words)
		fmt.Printf("%s\n\x1b[0m", str)

		// Advance PC past the string
		m.pc += uint16(len(words)*2) + 1
	default:
		fmt.Printf(" !! UNKNOWN: %02x\n", inst.code)
		panic("BYE")
		//m.pc += 1
	}
}

func (m *Machine) Store(loc uint16, val uint16) {
	if loc == 0 {
		// TODO: store on stack
	} else if loc > 0 && loc < 0x10 {
		// TODO: store in local variable
	} else {
		// Global variable
		addr := uint16(m.header.globalsAddr + ((loc - 16) * 2))
		decode.SetWord(m.mem, addr, val)
	}
}

func (m *Machine) Run() {
	for {
		m.Step()
	}
}
