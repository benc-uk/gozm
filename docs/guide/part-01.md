# Step by Step Guide - Core

First stage is to get a minimal Z Machine story file running and just something happening, we'll create a very simple story file that does some arithmetic operations and stores the results in global variables.

To see the code for this tutorial in context, [check out the `core` branch of the repository.](https://github.com/benc-uk/gozm/tree/core)

## Getting Started

Loading a complete game like Zork or anything else is WAY too complicated for a first step, there are test files out there but it's nice to be able to control & constrain what the story file does to keep things ultra simple.

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

NOTE: The `core.dump.txt` file has a full dump & disassembly of the compiled story file, which is SUPER useful for understanding what's going on under the hood, and will go some way to explaining wht we jump to address 0x0049F if you scroll nearly to the end of the file you will see `0049F 54 11 0E 12              ADD G01 0x0E -> G2`
This is dumped [using the unz tool mentioned in the readme](https://github.com/heasm66/UnZ)

## Implementing the Core

So where on earth do you start! To begin you'll need is an array for the machine memory and an program counter to point to the current instruction. The Z Machine story file is loaded directly & completely into the machine memory. It might be surprising, the compiled story file on disk is essentially a snapshot of memory, so we can just load the story file into a byte array that represents the Z Machine memory.

The first 64 bytes of the story file & memory is a special header, which contains important info such as the starting PC, and a range of offset memory addresses. The [spec doc has full details on the header](https://zspec.jaredreisinger.com/b-conventional-header) in suprisingly well written form.

I'm not going to walk through every step in this tutorial, but cover the important parts & main gotchas, you can always refer to the code on the `core` branch for full details and working code.

Code to represent core of the Z Machine looks something like this:

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

## Endianess and the Program Counter

**Gotcha #1** - The Z-Machine is big-endian. This means that multi-byte values (16-bit words, there's no 32 or 64 bit values) are stored with the most significant byte first. When reading or writing multi-byte values, make sure to handle the byte order correctly. For example:

```go
func GetWord(b []byte, offset uint16) uint16 {
	return uint16(b[offset])<<8 | uint16(b[offset+1])
}

func SetWord(b []byte, offset uint16, value uint16) {
	b[offset] = byte((value >> 8) & 0xFF)
	b[offset+1] = byte(value & 0xFF)
}
```

**Gotcha #2** - The PC will start off pointing at a CALL instruction which will call the routine called `main`, which makes sense. However implementing the CALL instruction is huuuuge pain in the bum. We'll do a cheat for now and hardcode the PC to point a little bit further on to a known location that does some arithmetic operations. In our case `0x0049F` is where the actual operations are in memory.

## The Execution Loop

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
	inst := m.decodeInst() // Fetch and decode the next instruction

	switch inst.code {
    case: 0x001: // NOP
      // Do nothing
      m.pc += 1
    case someOpCode:
      // Handle the opcode here!
      m.pc += inst.length // Advance the program counter
    default:
      panic("!! UNKNOWN OPCODE: %02x\n", inst.code)
  }
}
```

## Decoding and Implementing Instructions

OK so how do decode and implement the instructions we need?

The Z Machine has one of the [most mindboggling and poorly documented instruction sets I've ever seen](https://zspec.jaredreisinger.com/04-instructions), and I've implemented the full Z80 instruction set in an emulator. The challenges being various; variable numbers of operands (0OP, 1OP, 2OP and VAR), multiple instruction formats (short, long, variable), also a variety of operand forms (large, small, variable, etc), and various special cases. It's a nightmare.

To get started we only need to implement a tiny subset of the instruction set to handle two simple arithmetic operations. These are ADD (opcode: 0x54) and STORE (opcode: 0x0d)
Both instructions are 2OP & long form instructions. So decoding them isn't too bad. See the `decodeInst` function in `instructions.go` for details.

## Storing and Retrieving Variables

**Gotcha #3** - Handling variables. The Z Machine has mix of stack, local and global variables. Variables are referenced by a byte value. Where 0 means the top of the stack (we can ignore this for now), 1-15 are local variables, and 16-255 are global variables. For our simple example we only need to handle global variables. See the `getVar` and `setVar` functions in `zmachine.go` for details.

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

As you might imagine the `GetVar` function, is similar but reads the value instead 😀

## Ready To Run

At this point we should have a minimal Z Machine implementation that can load our simple story file and execute some basic arithmetic operations.

If we panic to stop execution when we hit an unknown opcode & dump the global variables, we can see that the arithmetic operations were performed correctly. Hopefully!

```
Memory dump at 02ac:
02ac: 0003 (0003)
02ae: 005d (0093)
02b0: 0015 (0021)
02b2: 0000 (0000)
panic: BYE!
```

Global variable 0 (points) is 3 as expected it was never changed.
Global variable 1 (score) is 93 as expected after being updated by STORE.
Global variable 2 (total) is 21 as expected after being calculated by ADD (7 + 14 = 21).

This all might seem like a lot of work just to mutate some values in an array! but it's a crucial starting point. From here we can gradually implement more instructions & functionality, the key ones being CALL which requires setting up a call stack, and PRINT which requires decoding strings from the Z Machine's unique string encoding format.
