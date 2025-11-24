// ============================================================================
// GoZm - Z-Machine interpreter written in Go
// Copyright (c) 2025 - Ben Coleman
// Decoding utilities for Z-machine number and string representations
// ============================================================================

package decode

// Lookups used for decoding z-chars, see: https://zspec.jaredreisinger.com/03-text#3_5_3
var alphabets = [][]rune{
	{
		'a', 'b', 'c', 'd', 'e', 'f', 'g',
		'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
		'p', 'q', 'r', 's', 't', 'u', 'v', 'w',
		'x', 'y', 'z',
	},
	{
		'A', 'B', 'C', 'D', 'E', 'F', 'G',
		'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O',
		'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W',
		'X', 'Y', 'Z',
	},
	{
		// Note value 0 would normally never be output, it means escape
		'%', '\n', '0', '1', '2', '3', '4', '5',
		'6', '7', '8', '9', '.', ',', '!', '?',
		'_', '#', '\'', '"', '/', '\\', '-', ':',
		'(', ')',
	},
}

// GetInt16 reads a 2-byte big-endian integer from the given byte slice at the specified offset.
func GetInt16(b []byte, offset int16) int16 {
	return int16(b[offset])<<8 | int16(b[offset+1])
}

// String decodes a Z-machine encoded string from the given byte slice.
// These are encoded as a series of 2-byte words,
// each containing three 5-bit Z-characters. It's weird AF
// https://zspec.jaredreisinger.com/03-text
func String(data []byte) string {
	result := ""
	zchars := make([]byte, len(data)*3/2)

	// Convert each 2-byte word into 3 Z-chars
	for i := 0; i < len(data); i += 2 {
		word := uint16(data[i])<<8 | uint16(data[i+1])
		zchars[i/2*3] = byte((word >> 10) & 0x1F)
		zchars[i/2*3+1] = byte((word >> 5) & 0x1F)
		zchars[i/2*3+2] = byte(word & 0x1F)

		// If the high bit is set, this is the end of the string
		if word&0x8000 != 0 {
			break
		}
	}

	// Decode Z-chars into a string
	alphabet := 0
	for _, zchar := range zchars {
		switch zchar {
		case 0:
			result += " " // Z-char 0 is space
		case 1:
			// Abbreviation handling would go here
		case 2:
			// Abbreviation handling would go here
		case 3:
			// Abbreviation handling would go here
		case 4:
			alphabet = 1 // Switch to upper case
		case 5:
			alphabet = 2 // Switch to punctuation
		default:
			result += string(alphabets[alphabet][zchar-6])
			alphabet = 0 // Reset to default alphabet after use, this is v3 behaviour
		}
	}

	return string(result)
}
