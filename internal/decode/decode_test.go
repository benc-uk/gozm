// This is a Go package for handling number representations in the Z-machine.
// The Z-machine uses 2-byte numbers in big-endian format.

package decode

import (
	"testing"
)

func TestDecodeString(t *testing.T) {
	data := []byte{
		0xb1, 0x1c, 0xd7, 0x03, 0x8e, 0x1c, 0xf1, 0x28,
		0x07, 0x9b, 0xe5,
	}

	result := StringBytes(data)
	expected := string("bar wibble baz")

	if result != expected {
		t.Errorf("decode.String failed: got %v, want %v", []byte(result), []byte(expected))
	}
}
