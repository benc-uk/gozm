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
func (cf *callFrame) Push(val uint16) {
	cf.stack = append(cf.stack, val)
}

// Pop a value from the call frame stack
func (cf *callFrame) Pop() uint16 {
	if len(cf.locals) == 0 {
		// TODO: Maybe panic here?
		return 0
	}

	if len(cf.stack) == 0 {
		return 0
	}

	val := cf.stack[len(cf.stack)-1]
	cf.stack = cf.stack[:len(cf.stack)-1]

	return val
}

func (cf *callFrame) Peek() uint16 {
	if len(cf.stack) == 0 {
		return 0
	}

	return cf.stack[len(cf.stack)-1]
}

func (cf *callFrame) SetTop(val uint16) {
	if len(cf.stack) == 0 {
		return
	}

	cf.stack[len(cf.stack)-1] = val
}

// Helper to get the current call frame
func (m *Machine) getCallFrame() *callFrame {
	if len(m.callStack) == 0 {
		panic("Frame underflow, no current call frame!")
	}

	cf := &m.callStack[len(m.callStack)-1]
	m.trace("Get call frame, depth=%d retaddr=%04x %+v\n", len(m.callStack), cf.returnAddr, cf)

	return cf
}

// Helper to add a new empty call frame to machine call stack
func (m *Machine) addCallFrame(localCount int) *callFrame {
	frame := callFrame{
		returnAddr: 0,
		locals:     make([]uint16, 15), //localCount),
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
	m.trace("Return: PC restored to %04X, store byte=%02X, will advance to %04X\n", frame.returnAddr, resultStoreLoc, frame.returnAddr+1)

	m.storeVar(uint16(resultStoreLoc), val)
	m.pc += 1
}
