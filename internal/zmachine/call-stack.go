package zmachine

type callFrame struct {
	returnAddr uint16
	locals     []uint16
	stack      []uint16
}

func (sf *callFrame) Push(val uint16) {
	sf.stack = append(sf.stack, val)
}

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
