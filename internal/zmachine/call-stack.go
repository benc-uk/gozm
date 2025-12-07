// =======================================================================
// Package: zmachine - Core Z-machine interpreter
// call-stack.go - Call stack management, for routine calls and returns
//
// Copyright (c) 2025 Ben Coleman. Licensed under the MIT License
// =======================================================================

package zmachine

// CallFrame represents a single routine call in the Z-machine call stack
type CallFrame struct {
	ReturnAddr uint32   `json:"return_addr"`
	Locals     []uint16 `json:"locals"`
	Stack      []uint16 `json:"stack"`
}

// Push a value onto the call frame stack
func (cf *CallFrame) Push(val uint16) {
	cf.Stack = append(cf.Stack, val)
}

// Pop a value from the call frame stack
func (cf *CallFrame) Pop() uint16 {
	if len(cf.Locals) == 0 {
		return 0 // Absolutely should never happen
	}

	if len(cf.Stack) == 0 {
		return 0
	}

	val := cf.Stack[len(cf.Stack)-1]
	cf.Stack = cf.Stack[:len(cf.Stack)-1]

	return val
}

// Some ops need to just peek at the top value of the stack
func (cf *CallFrame) Peek() uint16 {
	if len(cf.Stack) == 0 {
		return 0
	}

	return cf.Stack[len(cf.Stack)-1]
}

// Some ops need to just set the top value of the stack
func (cf *CallFrame) SetTop(val uint16) {
	if len(cf.Stack) == 0 {
		return
	}

	cf.Stack[len(cf.Stack)-1] = val
}

// Helper to get the current call frame
func (m *Machine) getCallFrame() *CallFrame {
	if len(m.callStack) == 0 {
		panic("Frame underflow, no current call frame!")
	}

	cf := &m.callStack[len(m.callStack)-1]
	m.trace("Get call frame, depth=%d retaddr=%08x %+v\n", len(m.callStack), cf.ReturnAddr, cf)

	return cf
}

// Helper to add a new empty call frame to machine call stack
func (m *Machine) addCallFrame() *CallFrame {
	frame := CallFrame{
		ReturnAddr: 0,
		Locals:     make([]uint16, 15), //localCount),
		Stack:      make([]uint16, 0),
	}

	m.callStack = append(m.callStack, frame)
	return &m.callStack[len(m.callStack)-1]
}

// Helper to return from a call with a value
func (m *Machine) returnFromCall(val uint16) {
	frame := m.getCallFrame()
	m.pc = frame.ReturnAddr
	m.callStack = m.callStack[:len(m.callStack)-1]

	// The next byte after a CALL is the variable to store the result in
	resultStoreLoc := m.mem[m.pc]
	m.trace("Return: PC restored to %08X, store byte=%02X, will advance to %08X\n", frame.ReturnAddr, resultStoreLoc, frame.ReturnAddr+1)

	m.storeVar(uint16(resultStoreLoc), val)
	m.pc += 1
}
