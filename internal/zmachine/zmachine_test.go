package zmachine

import (
	"fmt"
	"testing"

	"github.com/benc-uk/gozm/internal/decode"
)

// Implements a simple terminal interface for Z-machine IO
type TestTerminal struct{}

// TextOut outputs text to the console
func (c *TestTerminal) TextOut(text string) {
	fmt.Print(text)
}

func (c *TestTerminal) ReadInput() string {
	return ""
}

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
	m := NewMachine(s, 0, &TestTerminal{})

	// Add a call frame so local variables exist
	m.addCallFrame(15)

	// Set local variable should NOT affect globals
	m.storeVar(0x0f, 77)
	val0 := decode.GetWord(m.mem, m.globalsAddr)
	if val0 != 0 {
		t.Errorf("Expected global variable %d to be 0, got %d", 0, val0)
	}

	gVar := uint16(0)
	m.storeVar(globalVarStart+gVar, 55)
	val := decode.GetWord(m.mem, m.globalsAddr+gVar)
	if val != 55 {
		t.Errorf("Expected global variable %d to be 55, got %d", gVar, val)
	}

	gVar = 5
	m.storeVar(globalVarStart+gVar, 65535)
	val = decode.GetWord(m.mem, m.globalsAddr+gVar*2)
	if val != 65535 {
		t.Errorf("Expected global variable %d to be 55, got %d", gVar, val)
		m.DumpMem(m.globalsAddr, 10*2)
	}

}

func TestStoreNegatives(t *testing.T) {
	s := makeTestStory()
	m := NewMachine(s, 0, &TestTerminal{})

	gVar := uint16(240)

	// This should panic, we need trap and recover
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when accessing invalid global variable")
		}
	}()
	m.storeVar(globalVarStart+gVar, 123)
	val := decode.GetWord(m.mem, m.globalsAddr+gVar*2)
	if val != 123 {
		t.Errorf("Expected global variable %d to be 55, got %d", gVar, val)
		m.DumpMem(m.globalsAddr, 10*2)
	}
}

func TestIncPositive(t *testing.T) {
	s := makeTestStory()
	m := NewMachine(s, 0, &TestTerminal{})

	// Add a call frame so local variables exist
	m.addCallFrame(15)

	// Test incrementing a local variable with positive value
	m.storeVar(0x01, 10)
	m.addToVar(0x01, 1)
	val := m.getVar(0x01)
	if val != 11 {
		t.Errorf("Expected local variable to be 11 after increment, got %d", val)
	}

	// Test incrementing a global variable with positive value
	gVar := uint16(0)
	m.storeVar(globalVarStart+gVar, 100)
	m.addToVar(globalVarStart+gVar, 1)
	val = m.getVar(globalVarStart + gVar)
	if val != 101 {
		t.Errorf("Expected global variable to be 101 after increment, got %d", val)
	}
}

func TestIncNegative(t *testing.T) {
	s := makeTestStory()
	m := NewMachine(s, 0, &TestTerminal{})

	// Add a call frame so local variables exist
	m.addCallFrame(15)

	// Test incrementing a negative local variable
	negVal := int16(-10)
	m.storeVar(0x01, uint16(negVal))
	m.addToVar(0x01, 1)
	val := int16(m.getVar(0x01))
	if val != -9 {
		t.Errorf("Expected local variable to be -9 after increment, got %d", val)
	}
	// Test incrementing a negative global variable
	gVar := uint16(0)
	negGlobalVal := int16(-50)
	m.storeVar(globalVarStart+gVar, uint16(negGlobalVal))
	m.addToVar(globalVarStart+gVar, 1)
	val = int16(m.getVar(globalVarStart + gVar))
	if val != -49 {
		t.Errorf("Expected global variable to be -49 after increment, got %d", val)
	}
	// Test incrementing from negative to positive
	negOne := int16(-1)
	m.storeVar(0x02, uint16(negOne))
	m.addToVar(0x02, 1)
	if m.getVar(0x02) != 0 {
		t.Errorf("Expected local variable to be 0 after incrementing -1, got %d", m.getVar(0x02))
	}
}

func TestIncOverflow(t *testing.T) {
	s := makeTestStory()
	m := NewMachine(s, 0, &TestTerminal{})

	// Add a call frame so local variables exist
	m.addCallFrame(15)

	// Test incrementing a local variable to overflow
	maxVal := int16(32767)
	m.storeVar(0x01, uint16(maxVal))
	m.addToVar(0x01, 1)
	val := int16(m.getVar(0x01))
	if val != -32768 {
		t.Errorf("Expected local variable to overflow to -32768, got %d", val)
	}

	// Test incrementing a global variable to overflow
	gVar := uint16(0)
	m.storeVar(globalVarStart+gVar, uint16(maxVal))
	m.addToVar(globalVarStart+gVar, 1)
	val = int16(m.getVar(globalVarStart + gVar))
	if val != -32768 {
		t.Errorf("Expected global variable to overflow to -32768, got %d", val)
	}
}

func TestDecPositive(t *testing.T) {
	s := makeTestStory()
	m := NewMachine(s, 0, &TestTerminal{})

	// Add a call frame so local variables exist
	m.addCallFrame(15)

	// Test decrementing a local variable with positive value
	m.storeVar(0x01, 10)
	m.addToVar(0x01, -1)
	val := m.getVar(0x01)
	if val != 9 {
		t.Errorf("Expected local variable to be 9 after decrement, got %d", val)
	}

	// Test decrementing a global variable with positive value
	gVar := uint16(0)
	m.storeVar(globalVarStart+gVar, 100)
	m.addToVar(globalVarStart+gVar, -1)
	val = m.getVar(globalVarStart + gVar)
	if val != 99 {
		t.Errorf("Expected global variable to be 99 after decrement, got %d", val)
	}
}

func TestDecNegative(t *testing.T) {
	s := makeTestStory()
	m := NewMachine(s, 0, &TestTerminal{})

	// Add a call frame so local variables exist
	m.addCallFrame(15)

	// Test decrementing a negative local variable
	negVal := int16(-10)
	m.storeVar(0x01, uint16(negVal))
	m.addToVar(0x01, -1)
	val := int16(m.getVar(0x01))
	if val != -11 {
		t.Errorf("Expected local variable to be -11 after decrement, got %d", val)
	}

	// Test decrementing a negative global variable
	gVar := uint16(0)
	negGlobalVal := int16(-50)
	m.storeVar(globalVarStart+gVar, uint16(negGlobalVal))
	m.addToVar(globalVarStart+gVar, -1)
	val = int16(m.getVar(globalVarStart + gVar))
	if val != -51 {
		t.Errorf("Expected global variable to be -51 after decrement, got %d", val)
	}

	// Test decrementing from positive to negative
	m.storeVar(0x02, 1)
	m.addToVar(0x02, -2)
	val = int16(m.getVar(0x02))
	if val != -1 {
		t.Errorf("Expected local variable to be -1 after decrementing 1 by 2, got %d", val)
	}
}

func TestDecUnderflow(t *testing.T) {
	s := makeTestStory()
	m := NewMachine(s, 0, &TestTerminal{})

	// Add a call frame so local variables exist
	m.addCallFrame(15)

	// Test decrementing a local variable to underflow
	minVal := int16(-32768)
	m.storeVar(0x01, uint16(minVal))
	m.addToVar(0x01, -1)
	val := int16(m.getVar(0x01))
	if val != 32767 {
		t.Errorf("Expected local variable to underflow to 32767, got %d", val)
	}

	// Test decrementing a global variable to underflow
	gVar := uint16(0)
	m.storeVar(globalVarStart+gVar, uint16(minVal))
	m.addToVar(globalVarStart+gVar, -1)
	val = int16(m.getVar(globalVarStart + gVar))
	if val != 32767 {
		t.Errorf("Expected global variable to underflow to 32767, got %d", val)
	}
}
