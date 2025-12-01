// =======================================================================
// Package: zmachine - Core Z-machine interpreter
// call-stack.go - Call stack management, for routine calls and returns
//
// Copyright (c) 2025 Ben Coleman. Licensed under the MIT License
// =======================================================================

package zmachine

// callFrame represents a single routine call in the Z-machine call stack
type callFrame struct {
	returnAddr uint16
	locals     []uint16
	stack      []uint16
}

// Push a value onto the call frame stack
func (sf *callFrame) Push(val uint16) {
	sf.stack = append(sf.stack, val)
}

// Pop a value from the call frame stack
func (sf *callFrame) Pop() uint16 {
	if len(sf.locals) == 0 {
		// TODO: Maybe panic here?
		return 0
	}

	val := sf.stack[len(sf.stack)-1]
	sf.stack = sf.stack[:len(sf.stack)-1]
	return val
}

// Helper to get the current call frame
func (m *Machine) getCallFrame() *callFrame {
	if len(m.callStack) == 0 {
		panic("Frame underflow, no current call frame!")
	}

	return &m.callStack[len(m.callStack)-1]
}

// Helper to add a new empty call frame to machine call stack
func (m *Machine) addCallFrame(localCount int) *callFrame {
	frame := callFrame{
		returnAddr: 0,
		locals:     make([]uint16, localCount),
		stack:      make([]uint16, 0),
	}

	m.callStack = append(m.callStack, frame)
	return &m.callStack[len(m.callStack)-1]
}

// Helper to return from a call with a value
func (m *Machine) returnFromCall(val uint16) {
	frame := m.getCallFrame()
	m.pc = frame.returnAddr
	m.callStack = m.callStack[:len(m.callStack)-1]

	// The next byte after a CALL is the variable to store the result in
	resultStoreLoc := m.mem[m.pc]
	m.trace("Result %04x into %d\n", val, resultStoreLoc)

	m.storeVar(uint16(resultStoreLoc), val)
	m.pc += 1
}
