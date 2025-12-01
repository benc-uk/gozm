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
}
