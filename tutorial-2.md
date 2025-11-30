# Step by Step Guide - Decoding instructions, Strings & Call Stack

Next we'll tackle decoding instructions properly, handling CALLs and RETurns, and decoding & printing strings. This will be a lot more work, but the result will be a lot more functional and serve as a solid foundation for further development.

Working code for this tutorial can be found in the `tutorial-2` branch of the [gozm repository]

## Decoding Instructions

There's no avoiding needing to decode instructions properly, it's daunting AF but we need to bite the bullet.

### Guidance

The specification has all the details, but presented in it's usual impenetrable way. I'll put together some tips to help you navigate through this, you need to consider four main things when decoding an instruction:

- Each instruction has an opcode, which determines what operation to perform, pretty standard.
- Instructions come in three forms: short, long, and variable.
- Instructions can have between 0 and 4 operands.
- Operands can be of different types: large constant, small constant, variable, etc.

#### Getting Instruction Form

- If bits 7 and 6 of the opcode byte are `11` (0xC0), it's a variable form instruction.
- If bits 7 and 6 are `10` (0x80) it's a short form instruction.
- Otherwise, it's a long form instruction.

#### Getting Operand Types & Count

- Short form: bits 4 and 5 of the opcode byte indicate the operand type. These instructions can have 0 or 1 operands
- Variable form:, these have an extra byte after the opcode, this contains the operand types, with each operand type represented by 2 bits.
- Long form: bits 6 and 5 of the opcode byte indicate the types of the two operands. These instructions always have 2 operands. There's a gotcha here, as only a single bit is used, mapping to either small constant or variable (see below).

#### Operand Type Encoding

- 00 (0x0): Large constant (2 bytes holding a literal value)
- 01 (0x1): Small constant (1 byte holding a literal value)
- 10 (0x2): Variable (1 byte, refers to a variable number, remember 0 = stack, 1-15 = locals, 16-255 = globals)
- 11 (0x3): Omitted (no operand, used in short & variable form instructions)

Putting this all together you can decode any instruction. See the `decodeInst` function in `instructions.go` for a full implementation. A rough outline is shown below

```go
type instruction struct {
	code     byte     // opcode byte
	operands []uint16 // operand values for the instruction
	len      uint16   // total length of instruction + operands in bytes
}

func (m *Machine) decodeInst() instruction {
	inst := instruction{
		code: m.mem[m.pc],
		len:  1, // start with 1 for the opcode byte
	}

	// Short form
	if (inst.code & 0xC0) == 0x80 {
		opType := (inst.code >> 4) & 0x3	// Get bits 4 and 5 for operand type

		// TODO: fetch operand if not omitted

		return inst
	}

	// Variable form
	if (inst.code & 0xC0) == 0xC0 {

		opTypesByte := m.mem[m.pc+1]

		// Loop over the opTypesByte reading 2 bits at a time
		// - to determine operand types & count

		// TODO: fetch operands based on types

		return inst
	}

	// Long form
	op1Type := (inst.code>>6)&0x1 + 1 // See note about mapping below
	op2Type := (inst.code>>5)&0x1 + 1

	// TODO: fetch operands based on types

	return inst
}
```

You'll notice in the long form decoding we add 1 to the extracted bits to map them to the operand type encoding. This is because long form instructions only use a single bit to indicate either small constant (0) or variable (1), so we adjust the value to fit the same operand type scheme used elsewhere.

You'll also need a helper function to fetch operands based on their type, which reads from memory and returns the appropriate value. See the `fetchOperand` function in `instructions.go` for an example implementation.

## String Decoding

Next up is string decoding. Strings in the Z Machine are encoded in a special format called ZSCII, which is a bit tricky to decode. The characters are stored in a series of 5-bit chunks, packed into pairs of bytes. Yes, it's as weird as it sounds. TBH I don't think it ever really saved much space, but hey, it was the 80s.

If we want to print strings, we need to decode them properly. Here's a basic outline of how to do this:

- Read pairs of bytes from memory (aka words) from starting address, until you hit a word where the high bit (bit 15) is set, indicating the end of the string. This can tested by checking if `word & 0x8000 != 0`.
- For each word, extract the three 5-bit characters. The first character is in bits 10-14, the second in bits 5-9, and the third in bits 0-4.
- Map each 5-bit character to its corresponding ZSCII character.
- You need to handle special cases for shifting between alphabets (lowercase, uppercase, punctuation) and for ZSCII characters that require two 5-bit chunks.

> NOTE: The use of the term "word" here refers to a 2-byte (16-bit) unit of data, which is standard terminology in computing, not a linguistic word, which is confusing when discussing strings!

### Alphabet Mapping

The Z Machine has this concept of alphabets. There are three main alphabets:

- Alphabet 0: Lowercase letters (a-z)
- Alphabet 1: Uppercase letters (A-Z)
- Alphabet 2: Punctuation and special characters

Characters 4 and 5 are used to shift between these alphabets, in v3 specifically they are temporary shifts for the next character only. The alphabet resets back to 0 after each character.

It's all totally bonkers but we have to deal with it.

Here's a rough implementation of the string decoding function:

```go
func StringDecode(words []uint16) string {
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

		// Handle cases:
		// 0 = space
		// 1-3 = abbreviations (not implemented here)
		// 4-5 = alphabet shifts
		// 6+ = regular characters

		result += string(decodedChar)
	}
	return result
}
```

To read the string from memory, you need to read words starting from the string's start address until you hit a word with the high bit set. You can pass this is list of words to the decode function.

```go
words := make([]uint16, 0)
for i := uint16(0); int(i) < len(mem); i += 2 {
	word := decode.GetWord(mem, startOffset+i)
	words = append(words, word)

	// If the high bit is set, this is the end of the string
	if word&0x8000 != 0 {
		break
	}
}
```

## Implementing CALL and RET

One major piece missing from our interpreter is handling routine calls and returns. The Z Machine uses a call stack to manage function calls, which means we need to implement a stack structure to keep track of return addresses and local variables. It sounds scary but compared to instruction & string decoding it's relatively straightforward, especially in a high-level language like Go.

Routines are called using the CALL instruction, which jumps to the address of the routine, but we need to save the return address too, we'll create a call frame for this. The routine can have local variables, which are stored in the call frame as well. When a routine finishes executing, it uses one of the RET instructions to return to the saved return address (and calling routine), and the current call frame is discarded along with its local variables & stack.

Routines can also accept arguments, which are passed in and set as local variables in the call frame, they can also return a value, which is typically stored in a special variable or on the stack.

Routine Structure:

- Header:
  - Routines all start with a single byte that indicates the number of local variables.
  - This is followed by the starting values of local variables themselves (2 bytes each), depending on the compiler these may be initialized to 0
- After the header, the actual routine code begins.
- Routines have no explicit "end" marker, they simply return using a RET instruction.

### Call Stack & Frame Structure

We'll define a simple stack structure to hold our call frames. Each frame will store the return address and any local variables.

```go
type callFrame struct {
	returnAddr uint16    // Address to return to after the call
	locals     []uint16  // Local variables for the routine
	stack      []uint16  // Temporary stack for this frame, variables can be pushed/popped here
}

func (cf *callFrame) push(val uint16) {
	// Push value onto the frame's stack
}

func (cf *callFrame) pop() uint16 {
	// Pop & return value from the frame's stack
}
```

We'll add a slice to our `Machine` struct to hold the call stack:

```go
type Machine struct {
	// ...
	callStack []callFrame
	// ...
}
```

### Handling CALL Instruction

When we encounter a CALL instruction, we'll need to:

- Decode the routine address from the first instruction operand, see below on packed addresses.
- Determine the number of local variables from the routine header byte.
- Create a new call frame, setting the return address to the instruction after the CALL, and initializing local variables.
- Push the new call frame onto the call stack.
- Set the program counter (PC) to the routine address (past the header byte) to start executing it.

WTF is a packed address? In the Z-machine, routine addresses are often stored in a "packed" format to save space. Instead of using the full 16-bit address, they use a smaller representation that needs to be converted back to a full address when used. The exact method of unpacking depends on the version of the Z-machine, but for version 3, you typically [multiply the packed address by 2](https://zspec.jaredreisinger.com/01-memory-map#1_2_3) to get the actual address in memory. Just another weird quirk of the Z-machine!

There will also be one byte after the CALL + operands, that indicates the location to store the return value from the routine, This is a variable number (reminder: 0 = stack, 1-15 = locals, 16-255 = globals)

A rough implementation of the CALL handling might look like this:

```go
routineAddr := inst.operands[0] * 2 // Unpack the address
numLocals := m.mem[routineAddr]     // First byte is number of locals

// Push new stack frame
frame := m.addCallFrame(int(numLocals))
frame.returnAddr = m.pc + inst.len // This is the address after the CALL instruction, and will hold the variable to store return value

// Populate locals (word sized) from the initial values in the routine header
for i := 0; i < int(numLocals); i++ {
	localVal := decode.GetWord(m.mem, routineAddr+1+uint16(i*2))
	frame.locals[i] = localVal
}

// Copy arguments into local variables
if len(inst.operands) > 1 {
	for i, argVal := range inst.operands[1:] {
		frame.locals[i] = argVal
	}
}

// Set PC to the start of the routine code (after header)
m.pc = routineAddr + 1 + uint16(numLocals*2)
```

### Handling RET Instructions

When we encounter a RET instruction, we'll need to:

- Pop the current call frame off the call stack, it's discarded.
- Set the program counter (PC) to the return address stored in the popped call frame.
- All RET instructions return a value, we need to store this value in a variable (even it's ignored, dumped into a global var that isn't used). The location to store the return value is indicated by a byte immediately after the CALL instruction that initiated the routine call.

The RET instruction has a few variants, but the basic idea is the same. Here's a rough implementation of the RET handling:

```go
// RET_FALSE; val := 0
// RET_TRUE;  val := 1
// RET_VAL;   val := operand value
frame := m.currentCallFrame()
m.pc = frame.returnAddr // Set PC to return address after the CALL
m.callStack = m.callStack[:len(m.callStack)-1] // Pop the frame
resultStoreLoc := m.mem[m.pc] // Location to store return value, is the byte after the CALL instruction
m.StoreVar(uint16(resultStoreLoc), val)
m.pc += 1 // Move past the return value location byte
```

## Variable Storage & Retrieval

If you recall from earlier tutorials, we had a simple way to read and write variables, but we only implemented globals, which isn't sufficient now we have routines with local variables and a call stack. Thankfully it's not too hard to extend our existing `LoadVar` and `StoreVar` functions to handle locals and stack variables.

Here's an updated implementation of StoreVar()

```go
func (m *Machine) StoreVar(loc uint16, val uint16) {
	if loc == 0 {
		m.getCallFrame().Push(val)
	} else if loc > 0 && loc < 0x10 {
		m.getCallFrame().locals[loc-1] = val // Locals are 1-15, map to 0-14 index
	} else {
		addr := uint16(m.globalsAddr + (loc-0x10)*2)
		decode.SetWord(m.mem, addr, val)
	}
}
```

## Next Steps

We've covered a LOT of ground here, implementing instruction decoding, string decoding, and handling CALL and RET instructions with a call stack. This sets a solid foundation for further development of our Z-machine interpreter.

Implementing more instructions, is a matter of following the same patterns we've established here. Each instruction will have its own logic, but the overall structure remains consistent. Check out the `Step()` function in `zmachine.go` for how to handle different instructions.

Take a look at [stories/basic.inf](./stories/basic.inf) for an example story that uses these features. Try running the compiled stories/basic.z3 file with your interpreter to see the results.
