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
		' ', '\n', '0', '1', '2', '3', '4', '5',
		'6', '7', '8', '9', '.', ',', '!', '?',
		'_', '#', '\'', '"', '/', '\\', '-', ':',
		'(', ')',
	},
}

// GetWord reads a 2-byte big-endian integer from the given byte slice at the specified offset.
func GetWord(b []byte, offset uint16) uint16 {
	return uint16(b[offset])<<8 | uint16(b[offset+1])
}

// SetWord writes a 2-byte big-endian integer to the given byte slice at the specified offset.
func SetWord(b []byte, offset uint16, value uint16) {
	b[offset] = byte((value >> 8) & 0xFF)
	b[offset+1] = byte(value & 0xFF)
}

func PackedAddress(addr uint16) uint16 {
	return uint16(addr) * 2
}

// Helper to decode a Z-machine encoded string from a byte slice
func StringBytes(data []byte) string {
	words := []uint16{}
	for i := uint16(0); int(i) < len(data); i += 2 {
		word := GetWord(data, i)
		words = append(words, word)

		// If the high bit is set, this is the end of the string
		if word&0x8000 != 0 {
			break
		}
	}

	return String(words)
}

// String decodes a Z-machine encoded string from the given slice of 16-bit words
// each containing three 5-bit Z-characters. It's weird AF
// https://zspec.jaredreisinger.com/03-text
func String(words []uint16) string {
	result := ""
	zchars := make([]byte, len(words)*3)

	// Convert each 2-byte word into 3 Z-chars
	for i := 0; i < len(words); i++ {
		word := words[i]
		zchars[i*3] = byte((word >> 10) & 0x1F)
		zchars[i*3+1] = byte((word >> 5) & 0x1F)
		zchars[i*3+2] = byte(word & 0x1F)
	}

	// Decode Z-chars into a string
	alphabet := 0
	for i := 0; i < len(zchars); i++ {
		zchar := zchars[i]

		switch zchar {
		case 0:
			result += " " // Z-char 0 is space
			continue
		case 1:
			// Abbreviation handling would go here
		case 2:
			// Abbreviation handling would go here
		case 3:
			// Abbreviation handling would go here
		case 4:
			alphabet = 1 // Switch to upper case
			continue
		case 5:
			alphabet = 2 // Switch to punctuation
			continue
		case 6:
			// See https://zspec.jaredreisinger.com/03-text#3_4
			if alphabet == 2 {
				zc10 := (zchars[i+1] << 5) | zchars[i+2]
				result += getZSCIIChar(zc10)
				i += 2 // Skip next two zchars
				alphabet = 0
				continue
			}

			fallthrough
		default:
			result += string(alphabets[alphabet][zchar-6])
		}

		alphabet = 0 // Reset to default alphabet after use, this is v3 behaviour
	}

	return string(result)
}

func getZSCIIChar(zchar byte) string {
	if zchar >= 32 && zchar <= 126 {
		return string(rune(zchar))
	}

	// Handle special ZSCII characters here if needed, but for now return 0
	return string(rune(0))
}
