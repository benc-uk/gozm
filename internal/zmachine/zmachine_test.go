package zmachine

import (
	"os"
	"testing"
)

func TestStoreGlobalTrue(t *testing.T) {
	d, err := os.ReadFile("../../test/hello.z3")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	m := NewMachine(d)

	// Store 'true' (1) into global variable 0 (at address 0x0020)
	m.Store(0xff, 55)

}
