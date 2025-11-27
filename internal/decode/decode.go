// ============================================================================
// GoZm - Z-Machine interpreter written in Go
// Copyright (c) 2025 - Ben Coleman
// Decoding utilities for Z-machine number and string representations
// ============================================================================

package decode

// GetWord reads a 2-byte big-endian integer from the given byte slice at the specified offset.
func GetWord(b []byte, offset uint16) uint16 {
	return uint16(b[offset])<<8 | uint16(b[offset+1])
}

// SetWord writes a 2-byte big-endian integer to the given byte slice at the specified offset.
func SetWord(b []byte, offset uint16, value uint16) {
	b[offset] = byte((value >> 8) & 0xFF)
	b[offset+1] = byte(value & 0xFF)
}
