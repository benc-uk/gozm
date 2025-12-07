// =======================================================================
// Package: zmachine - Core Z-machine interpreter
// debug.go - Debugging utilities for the Z-machine
//
// Copyright (c) 2025 Ben Coleman. Licensed under the MIT License
// =======================================================================

package zmachine

import (
	"fmt"

	"github.com/benc-uk/gozm/internal/decode"
)

// DumpMem dumps a section of memory for debugging
func (m *Machine) DumpMem(addr uint16, length uint16) {
	fmt.Printf("\nMemory dump at %04x:\n", addr)
	for i := uint16(0); i < length; i += 2 {
		word := decode.GetWord(m.mem, addr+i)
		fmt.Printf("%04x: %04x (%04d)\n", addr+i, word, word)
	}
}

func (m *Machine) debug(format string, a ...interface{}) {
	if m.debugLevel > DEBUG_NONE {
		fmt.Printf("\033[32m"+format+"\033[0m", a...)
	}
}

func (m *Machine) trace(format string, a ...interface{}) {
	if m.debugLevel == DEBUG_TRACE {
		fmt.Printf("  \033[31m"+format+"\033[0m", a...)
	}
}

// Internal helper to dump object properties for debugging
func (o *zObject) propDebugDump() string {
	result := ""
	for _, propData := range o.Props {
		result += fmt.Sprintf("    Prop num:%d, size:%d, data:%x\n", propData.Num, propData.Size, propData.Data)
	}
	return result
}
