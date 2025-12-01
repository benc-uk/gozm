package zmachine

// External defines the interface for external functions provided to the Z-machine
type External interface {
	TextOut(text string)
	ReadInput() string
}
