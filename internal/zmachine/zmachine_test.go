package zmachine

import (
	"gozm/externals/terminal"
	"gozm/internal/decode"
	"testing"
)

const globalVarStart = uint16(0x10)

func makeTestStory() []byte {
	d := make([]byte, 1536)
	d[0x00] = 0x03 // Version 3

	// Set two bytes word for start of globals table at 0x02AC
	d[0x0C] = 0x02
	d[0x0D] = 0xAC

	d[0x04] = 0x04 // High memory address (0x0496)
	d[0x05] = 0x96

	return d
}

func TestStoreGlobal(t *testing.T) {
	s := makeTestStory()
	m := NewMachine(s, 0, terminal.NewTerminal())

	// Set local variable should NOT affect globals
	m.StoreVar(0x0f, 77)
	val0 := decode.GetWord(m.mem, m.globalsAddr)
	if val0 != 0 {
		t.Errorf("Expected global variable %d to be 0, got %d", 0, val0)
	}

	gVar := uint16(0)
	m.StoreVar(globalVarStart+gVar, 55)
	val := decode.GetWord(m.mem, m.globalsAddr+gVar)
	if val != 55 {
		t.Errorf("Expected global variable %d to be 55, got %d", gVar, val)
	}

	gVar = 5
	m.StoreVar(globalVarStart+gVar, 65535)
	val = decode.GetWord(m.mem, m.globalsAddr+gVar*2)
	if val != 65535 {
		t.Errorf("Expected global variable %d to be 55, got %d", gVar, val)
		m.DumpMem(m.globalsAddr, 10*2)
	}

}

func TestStoreNegatives(t *testing.T) {
	s := makeTestStory()
	m := NewMachine(s, 0, terminal.NewTerminal())

	gVar := uint16(240)

	// This should panic, we need trap and recover
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when accessing invalid global variable")
		}
	}()
	m.StoreVar(globalVarStart+gVar, 123)
	val := decode.GetWord(m.mem, m.globalsAddr+gVar*2)
	if val != 123 {
		t.Errorf("Expected global variable %d to be 55, got %d", gVar, val)
		m.DumpMem(m.globalsAddr, 10*2)
	}
}
