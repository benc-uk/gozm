# Opcodes revisited

I personally found the [opcode tables](14-opcode-table.md) very hard to interpret, despite the [section that tries to explain things](14-opcode-table.md#reading_the_opcode_tables). After looking at some output from Txd, however, the tables make a little more sense… but they are still (to me) awkwardly constructed, not easy to interpret, and not exhaustive.

Further, [Section 4](04-instructions.md) delves into how the bitfields in an opcode should be interpreted, but then the tables in [Section 14](14-opcode-table.md) don't list the instructions by opcode number, instead (kind of?) falling back to the raw opcode value. This raises the question whether the opcode number really has any meaning or not. Through some trial-and-error, I've determined that ["opcode number" as defined in S4.3](04-instructions.md#4_3) is not sufficient in and of itself to define a Z-Machine instruction. Why? Because instructions often require a particular number of operands, and more than one of the operand-count values (0OP, 1OP, 2OP, VAR) can have the same "opcode number" value. Further, because the two high-bits of the opcode define the form and not the operand count—and in fact the operand count might be implied by form and operand type!—it is not at all trivial to describe "what does instruction xxx look like?" both concisely and completely.

When trying to author an interpreter, however, you very much need to know how to answer that question. To that end, what follows is a re-creation of the opcode tables, keyed off the operand count (or form for extended opcodes!) and opcode number, but with every variant listed explicitly. Although the table below uses **Count** as a header, it's also used for "EXT" for extended opcode instructions. The original tables started with 2OP, but I believe that was an unexplicit side-effect of `$$00`-prefixed opcodes being long-form 2OP ones. Since there will necessarily be an intermingling of opcode value (to keep like opcode numbers together), this list uses the order 0OP, 1OP, 2OP, VAR, EXT. In this table, **Num** is the opcode number, and should represent an instruction in the context of a particular Count (with the caveat that it might change across Z-machine versions). The **Dec.**, **Hex**, and **Binary** columns are the complete opcode value, to aid in looking at hex dumps and writing an interpreter. (For the original table's TYPE:Decimal identifiers, look at Count:Decimal here.) The **Form** and **Args** (operands with a shorter label!) columns are abbreviated to S,L,V,E for Short, Long, Variable, and Extended for the form, and also to S,L,V for each operand for Small constant, Large constant, and Variable. Also, note that per-version variations might look like they are specific to particular decimal/hex opcodes, but they typically aren't; Rather than duplicate the entire set of opcodes for multi-variant instructions, they are only listed once (see 1OP:15 as an example). Truly, both the decimal/hex/binary/form/operand-types and the per-version instruction variations should be thought of as orthogonal nested tables, but the layout is already rather complex.

## Complete Opcode Tables

This comprehensive table includes all opcodes organized by operand count (0OP, 1OP, 2OP, VAR, EXT), showing decimal, hexadecimal, and binary representations along with forms, arguments, versions, store/branch indicators, and instruction syntax.

**Column Key:**

- **Count**: Operand count category (0OP, 1OP, 2OP, VAR, EXT)
- **Num**: Opcode number within category
- **Dec.**: Decimal opcode value
- **Hex**: Hexadecimal opcode value
- **Binary**: Binary opcode value
- **Form**: S=Short, L=Long, V=Variable, E=Extended
- **Args**: Operand types (S=Small constant, L=Large constant, V=Variable, ?=Variable count)
- **V**: Version number when opcode is available (or changes behavior)
- **St**: Store result indicator (\*)
- **Br**: Branch indicator (\*)
- **Instruction and syntax**: The opcode name and parameter syntax

Due to the extensive nature of the complete tables (covering all 0OP, 1OP, 2OP, VAR, and EXT opcodes with all their variants), please refer to:

- [Section 14: Opcode Table](14-opcode-table.md) for the traditional presentation
- [Section 15: Opcodes](15-opcodes.md) for detailed descriptions of each instruction

The key insight from this analysis is that a complete opcode identification requires both the operand count category AND the opcode number, as the same opcode number can exist across different categories with entirely different meanings.

---

_Source: [The Z-Machine Standards Document](https://zspec.jaredreisinger.com/ZZ03-opcodes)_
