// This is a Go package for handling number representations in the Z-machine.
// The Z-machine uses 2-byte numbers in big-endian format.

package decode

import (
	"testing"
)

func TestDecodeString(t *testing.T) {
	data := []byte{
		0x11, 0xAA, 0x46, 0x34, 0x14, 0xE4, 0x9D, 0x53,
	}

	result := String(data)
	expected := string("Hello\nBen")

	if result != expected {
		t.Errorf("decode.String failed: got %v, want %v", []byte(result), []byte(expected))
	}
}
