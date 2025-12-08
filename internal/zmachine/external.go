// =======================================================================
// Package: zmachine - Core Z-machine interpreter
// external.go - External interface so that host apps can provide I/O
//
// Copyright (c) 2025 Ben Coleman. Licensed under the MIT License
// =======================================================================

package zmachine

// External defines the interface for external functions provided to the Z-machine
type External interface {
	TextOut(text string)
	ReadInput() string
	PlaySound(soundID uint16, effect uint16, volume uint16)
	Save(state *SaveState) bool
	Load(name string, m *Machine) bool
}
