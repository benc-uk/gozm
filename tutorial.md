# Step by Step Guide - Core

Step 1 is to get a minimal Z Machine story file running that does some basic arithmetic operations. Which is in absolute minimum, but it's a start!

To see the code that starts things off, check the branch in this repository called `core`.

## Getting Started

Loading a complete game like Zork is WAY too complicated for a first step, there's test files out there but it's nice to be able to control the source code too and keep things ultra simple.

We can use Inform 6, which is a programming language and compiler specifically designed for creating interactive fiction games for the Z Machine. Use of this is WAY outside the scope of this guide, but you can find more information on the [Inform 6 website](https://inform-fiction.org/).

To get us started, we'll use a very simple Inform 6 source file that performs some basic arithmetic operations and stores the results in global variables. We won't be doing any text output or user interaction yet, as decoding strings and printing text is a whole other challenge!

```fortran
Global points = 3;
Global score = 7;
Global total = 0;

! Dumb test file to get us started with no printing or strings yet
[main;
  total = score + 14;
  score = 93;
];
```

Don't worry if you don't have the inform6 compiler installed the compiled story file `core.z3` is already included in the `stories` directory.

NOTE: The `core.dump.txt` file is a full dump & disassembly of the compiled story file, which is SUPER useful for understanding what's going on under the hood, and will go some way to explaining wht we jump to address 0x0049F if you scroll nearly to the end of the file you will see `0049F 54 11 0E 12              ADD G01 0x0E -> G2`
This is dumped [using the unz tool mentioned in the readme](https://github.com/heasm66/UnZ)

So where on earth do you start! To begin all you need is an array of bytes for the machine memory and an program counter to point to the current instruction. The Z Machine story file is loaded directly & completely into the machine memory. Yes the compiled story file on disk is essentially a snapshot of memory, so we can just load the story file into a byte array that represents the Z Machine memory.

I'm not going to explain every step in this tutorial, but cover the important parts & main gotchas, you can always refer to the code on the `core` branch for full details and working code.

Go code to represent core of the Z Machine looks something like this:

```go
type Machine struct {
	mem []byte
	pc  uint16

	version     byte
	globalsAddr uint16 // This base address is held in the header at 0x0C
}

// NewMachine creates a new Z-machine instance with the provided story data
func NewMachine(data []byte) *Machine {
	fmt.Printf("Z-machine initialized...\nVersion: %d\n", data[0x00])

	return &Machine{
		mem:         data,
		pc:          decode.GetWord(data, 0x06),  // Starting program counter
		version:     data[0x00],
		globalsAddr: decode.GetWord(data, 0x0C),  // Address of global variables
	}
}
```

_Gotcha #1_: The Z-Machine is big-endian. This means that multi-byte values (16-bit words, there's no 32 or 64 bit values) are stored with the most significant byte first. When reading or writing multi-byte values, make sure to handle the byte order correctly. For example:

```go
func GetWord(b []byte, offset uint16) uint16 {
	return uint16(b[offset])<<8 | uint16(b[offset+1])
}

func SetWord(b []byte, offset uint16, value uint16) {
	b[offset] = byte((value >> 8) & 0xFF)
	b[offset+1] = byte(value & 0xFF)
}
```

_Gotcha #2_: The PC will start off pointing at a CALL instruction which will call the routine called `main`, make sense yeah. However implementing the CALL instruction is huuuuge pain in the bum, so for now we just hardcode the PC to point to a known location that does some arithmetic operations. In our case `0x0049F` is where the arithmetic happens in our simple story file.

The Z Machine has one of the [most mindboggling and poorly documented instruction sets I've ever seen](https://zspec.jaredreisinger.com/04-instructions), and I've implemented a Z80 emulator before! The challenges being variable numbers of operands (0OP, 1OP, 2OP and VAR), multiple instruction formats (short, long, variable), also a variety of operand forms (large, small, variable, etc).

Each "CPU cycle" we fetch the next instruction, decode it, and execute it. The main execution loop looks something like this:

```go
func (m *Machine) Run() {
	fmt.Println("Starting the machine...")
	m.pc = 0x0049F // HACK: Point to known location otherwise we need to implement CALL and suffer mucho pain
	for {
		m.Step()
	}
}
```

The `Step` function is the heart of everything, it fetches the next instruction, decodes it, and executes it. For our simple arithmetic operations, we only need to implement two opcodes. Here's a simplified version of the `Step` function for the full implementation see the `zmachine.go` file:

```go
func (m *Machine) Step() {
	inst := m.decodeInst()

	switch inst.code {
    case: 0x001: // NOP
      // Do nothing
      m.pc += 1
    case someOpCode:
      // Handle the opcode
      m.pc += inst.length // Advance the program counter
    default:
      panic("!! UNKNOWN OPCODE: %02x\n", inst.code)
  }
}
```

To get started we only need to implement a tiny subset of the instruction set to handle two simple arithmetic operations. These are ADD (opcode: 0x54) and STORE (opcode: 0x0d)
Both instructions are 2OP long form instructions. So decoding them isn't too bad. See the `decodeInst` function in `instructions.go` for details.

_Gotcha #3_: Handling variables. The Z Machine has mix of stack, local and global variables. Variables are referenced by a byte value. Where 0 means the top of the stack (we can ignore this for now), 1-15 are local variables, and 16-255 are global variables. For our simple example we only need to handle global variables. See the `getVar` and `setVar` functions in `zmachine.go` for details.

A helper to set a variable:

```go
func (m *Machine) StoreVar(loc uint16, val uint16) {
	if loc == 0 {
		// TODO: store on stack
	} else if loc > 0 && loc < 0x10 {
		// TODO: store in local variable
	} else {
		// Global variable, which are all word sized
		addr := uint16(m.globalsAddr + (loc-0x10)*2)
		decode.SetWord(m.mem, addr, val)
	}
}
```

You can imagine the `GetVar` function is similar but reads the value instead.

So we have a minimal Z Machine implementation that can load our simple story file and execute some basic arithmetic operations.

If we panic to stop execution when we hit an unknown opcode & dump the global variables, we can see that the arithmetic operations were performed correctly. Hopefully!

```
Memory dump at 02ac:
02ac: 0003 (0003)
02ae: 005d (0093)
02b0: 0015 (0021)
02b2: 0000 (0000)
panic: BYE!
```

Global variable 0 (points) is 3 as expected.
Global variable 1 (score) is 93 as expected after being updated by STORE.
Global variable 2 (total) is 21 as expected after being calculated by ADD (7 + 14 = 21).

This all might seem like a lot of work just to mutate some values in an array, but it's a crucial first step. From here we can gradually implement more instructions, the key ones being CALL which requires setting up a call stack, and PRINT which requires decoding strings from the Z Machine's unique string encoding format.
