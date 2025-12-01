# Opcodes revisited

I personally found the [opcode tables](https://zspec.jaredreisinger.com/14-opcode-table) very hard to interpret, despite the [section that tries to explain things](https://zspec.jaredreisinger.com/14-opcode-table#reading_the_opcode_tables). After looking at some output from Txd, however, the tables make a little more sense… but they are still (to me) awkwardly constructed, not easy to interpret, and not exhaustive.

Further, [S4](https://zspec.jaredreisinger.com/04-instructions) delves into how the bitfields in an opcode should be interpreted, but then the tables in [S14](https://zspec.jaredreisinger.com/14-opcode-table) don't list the instructions by opcode number, instead (kind of?) falling back to the raw opcode value. This raises the question whether the opcode number really has any meaning or not. Through some trial-and-error, I've determined that ["opcode number" as defined in S4.3](https://zspec.jaredreisinger.com/04-instructions#4_3) is not sufficient in and of itself to define a Z-Machine instruction. Why? Because instructions often require a particular number of operands, and more than one of the operand-count values (0OP, 1OP, 2OP, VAR) can have the same "opcode number" value. Further, because the two high-bits of the opcode define the form and not the operand count—and in fact the operand count might be implied by form and operand type! it is not at all trivial to describe "what does instruction xxx look like?" both concisely and completely.

When trying to author an interpreter, however, you very much need to know how to answer that question. To that end, what follows is a re-creation of the opcode tables, keyed off the operand count (or form for extended opcodes!) and opcode number, but with every variant listed explicitly. Although the table below uses Count as a header, it's also used for "EXT" for extended opcode instructions. The original tables started with 2OP, but I believe that was an unexplicit side-effect of `$$00`-prefixed opcodes being long-form 2OP ones. Since there will necessarily be an intermingling of opcode value (to keep like opcode numbers together), this list uses the order 0OP, 1OP, 2OP, VAR, EXT. In this table, Num is the opcode number, and should represent an instruction in the context of a particular Count (with the caveat that it might change across Z-machine versions). The Dec., Hex, and Binary columns are the complete opcode value, to aid in looking at hex dumps and writing an interpreter. (For the original table's TYPE:Decimal identifiers, look at Count:Decimal here.) The Form and Args (operands with a shorter label!) columns are abbreviated to S,L,V,E for Short, Long, Variable, and Extended for the form, and als to S,L,V for each operand for Small constant, Large constant, and Variable. Also, note that per-version variations might look like they are specific to particular decimal/hex opcodes, but they typically aren't; Rather than duplicate the entire set of opcodes for multi-variant instructions, they are only listed once (see 1OP:15 as an example). Truly, both the decimal/hex/binary/form/operand-types and the per-version instruction variations should be thought of as orthogonal nested tables, but the layout is already rather complex.

## 0OP opcodes

| Count | Num | Dec. | Hex | Binary   | Form | Args | V   | St  | Br  | Instruction and syntax     |
| ----- | --- | ---- | --- | -------- | ---- | ---- | --- | --- | --- | -------------------------- |
| 0OP   | 0   | 176  | b0  | 10110000 | S    |      |     |     |     | rtrue                      |
| 0OP   | 1   | 177  | b1  | 10110001 | S    |      |     |     |     | rfalse                     |
| 0OP   | 2   | 178  | b2  | 10110010 | S    |      |     |     |     | print (literal string)     |
| 0OP   | 3   | 179  | b3  | 10110011 | S    |      |     |     |     | print_ret (literal string) |
| 0OP   | 4   | 180  | b4  | 10110100 | S    |      | 1/- |     |     | nop                        |
| 0OP   | 5   | 181  | b5  | 10110101 | S    |      | 1   |     | \*  | save ?(label)              |
|       |     |      |     |          |      |      | 4   |     |     | save → (result)            |
|       |     |      |     |          |      |      | 5   |     |     | illegal                    |
| 0OP   | 6   | 182  | b6  | 10110110 | S    |      | 1   |     | \*  | restore ?(label)           |
|       |     |      |     |          |      |      | 4   |     |     | restore → (result)         |
|       |     |      |     |          |      |      | 5   |     |     | illegal                    |
| 0OP   | 7   | 183  | b7  | 10110111 | S    |      |     |     |     | restart                    |
| 0OP   | 8   | 184  | b8  | 10111000 | S    |      |     |     |     | ret_popped                 |
| 0OP   | 9   | 185  | b9  | 10111001 | S    |      | 1   |     |     | pop                        |
|       |     |      |     |          |      |      | 5/6 | \*  |     | catch → (result)           |
| 0OP   | 10  | 186  | ba  | 10111010 | S    |      |     |     |     | quit                       |
| 0OP   | 11  | 187  | bb  | 10111011 | S    |      |     |     |     | new_line                   |
| 0OP   | 12  | 188  | bc  | 10111100 | S    |      | 3   |     |     | show_status                |
|       |     |      |     |          |      |      | 4   |     |     | illegal                    |
| 0OP   | 13  | 189  | bd  | 10111101 | S    |      | 3   |     |     | verify ?(label)            |

**Note:** 0OP:14 is 190 (0xbe), the first byte for an extended opcode (see EXT)

| Count | Num | Dec. | Hex | Binary   | Form | Args | V   | St  | Br  | Instruction and syntax |
| ----- | --- | ---- | --- | -------- | ---- | ---- | --- | --- | --- | ---------------------- |
| 0OP   | 15  | 191  | bf  | 10111111 | S    |      | 5/- |     |     | piracy                 |

## 1OP opcodes

| Count | Num | Dec. | Hex | Binary   | Form | Args | V   | St  | Br  | Instruction and syntax                   |
| ----- | --- | ---- | --- | -------- | ---- | ---- | --- | --- | --- | ---------------------------------------- |
| 1OP   | 0   | 128  | 80  | 10000000 | S    | L    |     |     | \*  | jz a ?(label)                            |
|       |     | 144  | 90  | 10010000 | S    | S    |     |     |     |                                          |
|       |     | 160  | a0  | 10100000 | S    | V    |     |     |     |                                          |
| 1OP   | 1   | 129  | 81  | 10000001 | S    | L    |     | \*  | \*  | get_sibling object → (result) ?(label)   |
|       |     | 145  | 91  | 10010001 | S    | S    |     |     |     |                                          |
|       |     | 161  | a1  | 10100001 | S    | V    |     |     |     |                                          |
| 1OP   | 2   | 130  | 82  | 10000010 | S    | L    |     | \*  | \*  | get_child object → (result) ?(label)     |
|       |     | 146  | 92  | 10010010 | S    | S    |     |     |     |                                          |
|       |     | 162  | a2  | 10100010 | S    | V    |     |     |     |                                          |
| 1OP   | 3   | 131  | 83  | 10000011 | S    | L    |     | \*  |     | get_parent object → (result)             |
|       |     | 147  | 93  | 10010011 | S    | S    |     |     |     |                                          |
|       |     | 163  | a3  | 10100011 | S    | V    |     |     |     |                                          |
| 1OP   | 4   | 132  | 84  | 10000100 | S    | L    |     | \*  |     | get_prop_len property-address → (result) |
|       |     | 148  | 94  | 10010100 | S    | S    |     |     |     |                                          |
|       |     | 164  | a4  | 10100100 | S    | V    |     |     |     |                                          |
| 1OP   | 5   | 133  | 85  | 10000101 | S    | L    |     |     |     | inc (variable)                           |
|       |     | 149  | 95  | 10010101 | S    | S    |     |     |     |                                          |
|       |     | 165  | a5  | 10100101 | S    | V    |     |     |     |                                          |
| 1OP   | 6   | 134  | 86  | 10000110 | S    | L    |     |     |     | dec (variable)                           |
|       |     | 150  | 96  | 10010110 | S    | S    |     |     |     |                                          |
|       |     | 166  | a6  | 10100110 | S    | V    |     |     |     |                                          |
| 1OP   | 7   | 135  | 87  | 10000111 | S    | L    |     |     |     | print_addr byte-address-of-string        |
|       |     | 151  | 97  | 10010111 | S    | S    |     |     |     |                                          |
|       |     | 167  | a7  | 10100111 | S    | V    |     |     |     |                                          |
| 1OP   | 8   | 136  | 88  | 10001000 | S    | L    | 4   | \*  |     | call_1s routine → (result)               |
|       |     | 152  | 98  | 10011000 | S    | S    |     |     |     |                                          |
|       |     | 168  | a8  | 10101000 | S    | V    |     |     |     |                                          |
| 1OP   | 9   | 137  | 89  | 10001001 | S    | L    |     |     |     | remove_obj object                        |
|       |     | 153  | 99  | 10011001 | S    | S    |     |     |     |                                          |
|       |     | 169  | a9  | 10101001 | S    | V    |     |     |     |                                          |
| 1OP   | 10  | 138  | 8a  | 10001010 | S    | L    |     |     |     | print_obj object                         |
|       |     | 154  | 9a  | 10011010 | S    | S    |     |     |     |                                          |
|       |     | 170  | aa  | 10101010 | S    | V    |     |     |     |                                          |
| 1OP   | 11  | 139  | 8b  | 10001011 | S    | L    |     |     |     | ret value                                |
|       |     | 155  | 9b  | 10011011 | S    | S    |     |     |     |                                          |
|       |     | 171  | ab  | 10101011 | S    | V    |     |     |     |                                          |
| 1OP   | 12  | 140  | 8c  | 10001100 | S    | L    |     |     |     | jump ?(label)                            |
|       |     | 156  | 9c  | 10011100 | S    | S    |     |     |     |                                          |
|       |     | 172  | ac  | 10101100 | S    | V    |     |     |     |                                          |
| 1OP   | 13  | 141  | 8d  | 10001101 | S    | L    |     |     |     | print_paddr packed-address-of-string     |
|       |     | 157  | 9d  | 10011101 | S    | S    |     |     |     |                                          |
|       |     | 173  | ad  | 10101101 | S    | V    |     |     |     |                                          |
| 1OP   | 14  | 142  | 8e  | 10001110 | S    | L    |     | \*  |     | load (variable) → (result)               |
|       |     | 158  | 9e  | 10011110 | S    | S    |     |     |     |                                          |
|       |     | 174  | ae  | 10101110 | S    | V    |     |     |     |                                          |
| 1OP   | 15  | 143  | 8f  | 10001111 | S    | L    | 1/4 | \*  |     | not value → (result)                     |
|       |     | 159  | 9f  | 10011111 | S    | S    | 5   |     |     | call_1n routine                          |
|       |     | 175  | af  | 10101111 | S    | V    |     |     |     |                                          |

## 2OP opcodes

| Count | Num | Dec. | Hex | Binary   | Form | Args | V   | St  | Br  | Instruction and syntax                   |
| ----- | --- | ---- | --- | -------- | ---- | ---- | --- | --- | --- | ---------------------------------------- |
| 2OP   | 0   | 0    | 00  | 00000000 | L    | S,S  | ―   | ―   | ―   | ―                                        |
|       |     | 32   | 20  | 00100000 | L    | S,V  |     |     |     |                                          |
|       |     | 64   | 40  | 01000000 | L    | V,S  |     |     |     |                                          |
|       |     | 96   | 60  | 01100000 | L    | V,V  |     |     |     |                                          |
|       |     | 192  | c0  | 11000000 | V    | ?    |     |     |     |                                          |
| 2OP   | 1   | 1    | 01  | 00000001 | L    | S,S  |     |     | \*  | je a b ?(label)                          |
|       |     | 33   | 21  | 00100001 | L    | S,V  |     |     |     |                                          |
|       |     | 65   | 41  | 01000001 | L    | V,S  |     |     |     |                                          |
|       |     | 97   | 61  | 01100001 | L    | V,V  |     |     |     |                                          |
|       |     | 193  | c1  | 11000001 | V    | ?    |     |     |     |                                          |
| 2OP   | 2   | 2    | 02  | 00000010 | L    | S,S  |     |     | \*  | jl a b ?(label)                          |
|       |     | 34   | 22  | 00100010 | L    | S,V  |     |     |     |                                          |
|       |     | 66   | 42  | 01000010 | L    | V,S  |     |     |     |                                          |
|       |     | 98   | 62  | 01100010 | L    | V,V  |     |     |     |                                          |
|       |     | 194  | c2  | 11000010 | V    | ?    |     |     |     |                                          |
| 2OP   | 3   | 3    | 03  | 00000011 | L    | S,S  |     |     | \*  | jg a b ?(label)                          |
|       |     | 35   | 23  | 00100011 | L    | S,V  |     |     |     |                                          |
|       |     | 67   | 43  | 01000011 | L    | V,S  |     |     |     |                                          |
|       |     | 99   | 63  | 01100011 | L    | V,V  |     |     |     |                                          |
|       |     | 195  | c3  | 11000011 | V    | ?    |     |     |     |                                          |
| 2OP   | 4   | 4    | 04  | 00000100 | L    | S,S  |     |     | \*  | dec_chk (variable) value ?(label)        |
|       |     | 36   | 24  | 00100100 | L    | S,V  |     |     |     |                                          |
|       |     | 68   | 44  | 01000100 | L    | V,S  |     |     |     |                                          |
|       |     | 100  | 64  | 01100100 | L    | V,V  |     |     |     |                                          |
|       |     | 196  | c4  | 11000100 | V    | ?    |     |     |     |                                          |
| 2OP   | 5   | 5    | 05  | 00000101 | L    | S,S  |     |     | \*  | inc_chk (variable) value ?(label)        |
|       |     | 37   | 25  | 00100101 | L    | S,V  |     |     |     |                                          |
|       |     | 69   | 45  | 01000101 | L    | V,S  |     |     |     |                                          |
|       |     | 101  | 65  | 01100101 | L    | V,V  |     |     |     |                                          |
|       |     | 197  | c5  | 11000101 | V    | ?    |     |     |     |                                          |
| 2OP   | 6   | 6    | 06  | 00000110 | L    | S,S  |     |     | \*  | jin obj1 obj2 ?(label)                   |
|       |     | 38   | 26  | 00100110 | L    | S,V  |     |     |     |                                          |
|       |     | 70   | 46  | 01000110 | L    | V,S  |     |     |     |                                          |
|       |     | 102  | 66  | 01100110 | L    | V,V  |     |     |     |                                          |
|       |     | 198  | c6  | 11000110 | V    | ?    |     |     |     |                                          |
| 2OP   | 7   | 7    | 07  | 00000111 | L    | S,S  |     |     | \*  | test bitmap flags ?(label)               |
|       |     | 39   | 27  | 00100111 | L    | S,V  |     |     |     |                                          |
|       |     | 71   | 47  | 01000111 | L    | V,S  |     |     |     |                                          |
|       |     | 103  | 67  | 01100111 | L    | V,V  |     |     |     |                                          |
|       |     | 199  | c7  | 11000111 | V    | ?    |     |     |     |                                          |
| 2OP   | 8   | 8    | 08  | 00001000 | L    | S,S  |     | \*  |     | or a b → (result)                        |
|       |     | 40   | 28  | 00101000 | L    | S,V  |     |     |     |                                          |
|       |     | 72   | 48  | 01001000 | L    | V,S  |     |     |     |                                          |
|       |     | 104  | 68  | 01101000 | L    | V,V  |     |     |     |                                          |
|       |     | 200  | c8  | 11001000 | V    | ?    |     |     |     |                                          |
| 2OP   | 9   | 9    | 09  | 00001001 | L    | S,S  |     | \*  |     | and a b → (result)                       |
|       |     | 41   | 29  | 00101001 | L    | S,V  |     |     |     |                                          |
|       |     | 73   | 49  | 01001001 | L    | V,S  |     |     |     |                                          |
|       |     | 105  | 69  | 01101001 | L    | V,V  |     |     |     |                                          |
|       |     | 201  | c9  | 11001001 | V    | ?    |     |     |     |                                          |
| 2OP   | 10  | 10   | 0a  | 00001010 | L    | S,S  |     |     | \*  | test_attr object attribute ?(label)      |
|       |     | 42   | 2a  | 00101010 | L    | S,V  |     |     |     |                                          |
|       |     | 74   | 4a  | 01001010 | L    | V,S  |     |     |     |                                          |
|       |     | 106  | 6a  | 01101010 | L    | V,V  |     |     |     |                                          |
|       |     | 202  | ca  | 11001010 | V    | ?    |     |     |     |                                          |
| 2OP   | 11  | 11   | 0b  | 00001011 | L    | S,S  |     |     |     | set_attr object attribute                |
|       |     | 43   | 2b  | 00101011 | L    | S,V  |     |     |     |                                          |
|       |     | 75   | 4b  | 01001011 | L    | V,S  |     |     |     |                                          |
|       |     | 107  | 6b  | 01101011 | L    | V,V  |     |     |     |                                          |
|       |     | 203  | cb  | 11001011 | V    | ?    |     |     |     |                                          |
| 2OP   | 12  | 12   | 0c  | 00001100 | L    | S,S  |     |     |     | clear_attr object attribute              |
|       |     | 44   | 2c  | 00101100 | L    | S,V  |     |     |     |                                          |
|       |     | 76   | 4c  | 01001100 | L    | V,S  |     |     |     |                                          |
|       |     | 108  | 6c  | 01101100 | L    | V,V  |     |     |     |                                          |
|       |     | 204  | cc  | 11001100 | V    | ?    |     |     |     |                                          |
| 2OP   | 13  | 13   | 0d  | 00001101 | L    | S,S  |     |     |     | store (variable) value                   |
|       |     | 45   | 2d  | 00101101 | L    | S,V  |     |     |     |                                          |
|       |     | 77   | 4d  | 01001101 | L    | V,S  |     |     |     |                                          |
|       |     | 109  | 6d  | 01101101 | L    | V,V  |     |     |     |                                          |
|       |     | 205  | cd  | 11001101 | V    | ?    |     |     |     |                                          |
| 2OP   | 14  | 14   | 0e  | 00001110 | L    | S,S  |     |     |     | insert_obj object destination            |
|       |     | 46   | 2e  | 00101110 | L    | S,V  |     |     |     |                                          |
|       |     | 78   | 4e  | 01001110 | L    | V,S  |     |     |     |                                          |
|       |     | 110  | 6e  | 01101110 | L    | V,V  |     |     |     |                                          |
|       |     | 206  | ce  | 11001110 | V    | ?    |     |     |     |                                          |
| 2OP   | 15  | 15   | 0f  | 00001111 | L    | S,S  |     | \*  |     | loadw array word-index → (result)        |
|       |     | 47   | 2f  | 00101111 | L    | S,V  |     |     |     |                                          |
|       |     | 79   | 4f  | 01001111 | L    | V,S  |     |     |     |                                          |
|       |     | 111  | 6f  | 01101111 | L    | V,V  |     |     |     |                                          |
|       |     | 207  | cf  | 11001111 | V    | ?    |     |     |     |                                          |
| 2OP   | 16  | 16   | 10  | 00010000 | L    | S,S  |     | \*  |     | loadb array byte-index → (result)        |
|       |     | 48   | 30  | 00110000 | L    | S,V  |     |     |     |                                          |
|       |     | 80   | 50  | 01010000 | L    | V,S  |     |     |     |                                          |
|       |     | 112  | 70  | 01110000 | L    | V,V  |     |     |     |                                          |
|       |     | 208  | d0  | 11010000 | V    | ?    |     |     |     |                                          |
| 2OP   | 17  | 17   | 11  | 00010001 | L    | S,S  |     | \*  |     | get_prop object property → (result)      |
|       |     | 49   | 31  | 00110001 | L    | S,V  |     |     |     |                                          |
|       |     | 81   | 51  | 01010001 | L    | V,S  |     |     |     |                                          |
|       |     | 113  | 71  | 01110001 | L    | V,V  |     |     |     |                                          |
|       |     | 209  | d1  | 11010001 | V    | ?    |     |     |     |                                          |
| 2OP   | 18  | 18   | 12  | 00010010 | L    | S,S  |     | \*  |     | get_prop_addr object property → (result) |
|       |     | 50   | 32  | 00110010 | L    | S,V  |     |     |     |                                          |
|       |     | 82   | 52  | 01010010 | L    | V,S  |     |     |     |                                          |
|       |     | 114  | 72  | 01110010 | L    | V,V  |     |     |     |                                          |
|       |     | 210  | d2  | 11010010 | V    | ?    |     |     |     |                                          |
| 2OP   | 19  | 19   | 13  | 00010011 | L    | S,S  |     | \*  |     | get_next_prop object property → (result) |
|       |     | 51   | 33  | 00110011 | L    | S,V  |     |     |     |                                          |
|       |     | 83   | 53  | 01010011 | L    | V,S  |     |     |     |                                          |
|       |     | 115  | 73  | 01110011 | L    | V,V  |     |     |     |                                          |
|       |     | 211  | d3  | 11010011 | V    | ?    |     |     |     |                                          |
| 2OP   | 20  | 20   | 14  | 00010100 | L    | S,S  |     | \*  |     | add a b → (result)                       |
|       |     | 52   | 34  | 00110100 | L    | S,V  |     |     |     |                                          |
|       |     | 84   | 54  | 01010100 | L    | V,S  |     |     |     |                                          |
|       |     | 116  | 74  | 01110100 | L    | V,V  |     |     |     |                                          |
|       |     | 212  | d4  | 11010100 | V    | ?    |     |     |     |                                          |
| 2OP   | 21  | 21   | 15  | 00010101 | L    | S,S  |     | \*  |     | sub a b → (result)                       |
|       |     | 53   | 35  | 00110101 | L    | S,V  |     |     |     |                                          |
|       |     | 85   | 55  | 01010101 | L    | V,S  |     |     |     |                                          |
|       |     | 117  | 75  | 01110101 | L    | V,V  |     |     |     |                                          |
|       |     | 213  | d5  | 11010101 | V    | ?    |     |     |     |                                          |
| 2OP   | 22  | 22   | 16  | 00010110 | L    | S,S  |     | \*  |     | mul a b → (result)                       |
|       |     | 54   | 36  | 00110110 | L    | S,V  |     |     |     |                                          |
|       |     | 86   | 56  | 01010110 | L    | V,S  |     |     |     |                                          |
|       |     | 118  | 76  | 01110110 | L    | V,V  |     |     |     |                                          |
|       |     | 214  | d6  | 11010110 | V    | ?    |     |     |     |                                          |
| 2OP   | 23  | 23   | 17  | 00010111 | L    | S,S  |     | \*  |     | div a b → (result)                       |
|       |     | 55   | 37  | 00110111 | L    | S,V  |     |     |     |                                          |
|       |     | 87   | 57  | 01010111 | L    | V,S  |     |     |     |                                          |
|       |     | 119  | 77  | 01110111 | L    | V,V  |     |     |     |                                          |
|       |     | 215  | d7  | 11010111 | V    | ?    |     |     |     |                                          |
| 2OP   | 24  | 24   | 18  | 00011000 | L    | S,S  |     | \*  |     | mod a b → (result)                       |
|       |     | 56   | 38  | 00111000 | L    | S,V  |     |     |     |                                          |
|       |     | 88   | 58  | 01011000 | L    | V,S  |     |     |     |                                          |
|       |     | 120  | 78  | 01111000 | L    | V,V  |     |     |     |                                          |
|       |     | 216  | d8  | 11011000 | V    | ?    |     |     |     |                                          |
| 2OP   | 25  | 25   | 19  | 00011001 | L    | S,S  | 4   | \*  |     | call_2s routine arg1 → (result)          |
|       |     | 57   | 39  | 00111001 | L    | S,V  |     |     |     |                                          |
|       |     | 89   | 59  | 01011001 | L    | V,S  |     |     |     |                                          |
|       |     | 121  | 79  | 01111001 | L    | V,V  |     |     |     |                                          |
|       |     | 217  | d9  | 11011001 | V    | ?    |     |     |     |                                          |
| 2OP   | 26  | 26   | 1a  | 00011010 | L    | S,S  | 5   |     |     | call_2n routine arg1                     |
|       |     | 58   | 3a  | 00111010 | L    | S,V  |     |     |     |                                          |
|       |     | 90   | 5a  | 01011010 | L    | V,S  |     |     |     |                                          |
|       |     | 122  | 7a  | 01111010 | L    | V,V  |     |     |     |                                          |
|       |     | 218  | da  | 11011010 | V    | ?    |     |     |     |                                          |
| 2OP   | 27  | 27   | 1b  | 00011011 | L    | S,S  | 5   |     |     | set_colour foreground background         |
|       |     | 59   | 3b  | 00111011 | L    | S,V  | 6   |     |     | set_colour foreground background window  |
|       |     | 91   | 5b  | 01011011 | L    | V,S  |     |     |     |                                          |
|       |     | 123  | 7b  | 01111011 | L    | V,V  |     |     |     |                                          |
|       |     | 219  | db  | 11011011 | V    | ?    |     |     |     |                                          |
| 2OP   | 28  | 28   | 1c  | 00011100 | L    | S,S  | 5/6 |     |     | throw value stack-frame                  |
|       |     | 60   | 3c  | 00111100 | L    | S,V  |     |     |     |                                          |
|       |     | 92   | 5c  | 01011100 | L    | V,S  |     |     |     |                                          |
|       |     | 124  | 7c  | 01111100 | L    | V,V  |     |     |     |                                          |
|       |     | 220  | dc  | 11011100 | V    | ?    |     |     |     |                                          |
| 2OP   | 29  | 29   | 1d  | 00011101 | L    | S,S  | ―   | ―   | ―   | ―                                        |
|       |     | 61   | 3d  | 00111101 | L    | S,V  |     |     |     |                                          |
|       |     | 93   | 5d  | 01011101 | L    | V,S  |     |     |     |                                          |
|       |     | 125  | 7d  | 01111101 | L    | V,V  |     |     |     |                                          |
|       |     | 221  | dd  | 11011101 | V    | ?    |     |     |     |                                          |
| 2OP   | 30  | 30   | 1e  | 00011110 | L    | S,S  | ―   | ―   | ―   | ―                                        |
|       |     | 62   | 3e  | 00111110 | L    | S,V  |     |     |     |                                          |
|       |     | 94   | 5e  | 01011110 | L    | V,S  |     |     |     |                                          |
|       |     | 126  | 7e  | 01111110 | L    | V,V  |     |     |     |                                          |
|       |     | 222  | de  | 11011110 | V    | ?    |     |     |     |                                          |
| 2OP   | 31  | 31   | 1f  | 00011111 | L    | S,S  | ―   | ―   | ―   | ―                                        |
|       |     | 63   | 3f  | 00111111 | L    | S,V  |     |     |     |                                          |
|       |     | 95   | 5f  | 01011111 | L    | V,S  |     |     |     |                                          |
|       |     | 127  | 7f  | 01111111 | L    | V,V  |     |     |     |                                          |
|       |     | 223  | df  | 11011111 | V    | ?    |     |     |     |                                          |

## VAR opcodes

| Count | Num | Dec. | Hex | Binary   | Form | Args | V   | St  | Br  | Instruction and syntax                        |
| ----- | --- | ---- | --- | -------- | ---- | ---- | --- | --- | --- | --------------------------------------------- |
| VAR   | 0   | 224  | e0  | 11100000 | V    | ?    | 1   | \*  |     | call routine …​0 to 3 args…​ → (result)       |
|       |     |      |     |          |      |      | 4   |     |     | call_vs routine …​0 to 3 args…​ → (result)    |
| VAR   | 1   | 225  | e1  | 11100001 | V    | ?    |     |     |     | storew array word-index value                 |
| VAR   | 2   | 226  | e2  | 11100010 | V    | ?    |     |     |     | storeb array byte-index value                 |
| VAR   | 3   | 227  | e3  | 11100011 | V    | ?    |     |     |     | put_prop object property value                |
| VAR   | 4   | 228  | e4  | 11100100 | V    | ?    | 1   |     |     | sread text parse                              |
|       |     |      |     |          |      |      | 4   |     |     | sread text parse time routine                 |
|       |     |      |     |          |      |      | 5   | \*  |     | aread text parse time routine → (result)      |
| VAR   | 5   | 229  | e5  | 11100101 | V    | ?    |     |     |     | print_char output-character-code              |
| VAR   | 6   | 230  | e6  | 11100110 | V    | ?    |     |     |     | print_num value                               |
| VAR   | 7   | 231  | e7  | 11100111 | V    | ?    |     | \*  |     | random range → (result)                       |
| VAR   | 8   | 232  | e8  | 11101000 | V    | ?    |     |     |     | push value                                    |
| VAR   | 9   | 233  | e9  | 11101001 | V    | ?    | 1   |     |     | pull (variable)                               |
|       |     |      |     |          |      |      | 6   | \*  |     | pull stack → (result)                         |
| VAR   | 10  | 234  | ea  | 11101010 | V    | ?    | 3   |     |     | split_window lines                            |
| VAR   | 11  | 235  | eb  | 11101011 | V    | ?    | 3   |     |     | set_window window                             |
| VAR   | 12  | 236  | ec  | 11101100 | V    | ?    | 4   | \*  |     | call_vs2 routine …​0 to 7 args…​ → (result)   |
| VAR   | 13  | 237  | ed  | 11101101 | V    | ?    | 4   |     |     | erase_window window                           |
| VAR   | 14  | 238  | ee  | 11101110 | V    | ?    | 4/- |     |     | erase_line value                              |
|       |     |      |     |          |      |      | 6   |     |     | erase_line pixels                             |
| VAR   | 15  | 239  | ef  | 11101111 | V    | ?    | 4   |     |     | set_cursor line column                        |
|       |     |      |     |          |      |      | 6   |     |     | set_cursor line column window                 |
| VAR   | 16  | 240  | f0  | 11110000 | V    | ?    | 4/6 |     |     | get_cursor array                              |
| VAR   | 17  | 241  | f1  | 11110001 | V    | ?    | 4   |     |     | set_text_style style                          |
| VAR   | 18  | 242  | f2  | 11110010 | V    | ?    | 4   |     |     | buffer_mode flag                              |
| VAR   | 19  | 243  | f3  | 11110011 | V    | ?    | 3   |     |     | output_stream number                          |
|       |     |      |     |          |      |      | 5   |     |     | output_stream number table                    |
|       |     |      |     |          |      |      | 6   |     |     | output_stream number table width              |
| VAR   | 20  | 244  | f4  | 11110100 | V    | ?    | 3   |     |     | input_stream number                           |
| VAR   | 21  | 245  | f5  | 11110101 | V    | ?    | 5/3 |     |     | sound_effect number effect volume routine     |
| VAR   | 22  | 246  | f6  | 11110110 | V    | ?    | 4   | \*  |     | read_char 1 time routine → (result)           |
| VAR   | 23  | 247  | f7  | 11110111 | V    | ?    | 4   | \*  | \*  | scan_table x table len form → (result)        |
| VAR   | 24  | 248  | f8  | 11111000 | V    | ?    | 5/6 | \*  |     | not value → (result)                          |
| VAR   | 25  | 249  | f9  | 11111001 | V    | ?    | 5   |     |     | call_vn routine …​up to 3 args…​              |
| VAR   | 26  | 250  | fa  | 11111010 | V    | ?    | 5   |     |     | call_vn2 routine …​up to 7 args…​             |
| VAR   | 27  | 251  | fb  | 11111011 | V    | ?    | 5   |     |     | tokenise text parse dictionary flag           |
| VAR   | 28  | 252  | fc  | 11111100 | V    | ?    | 5   |     |     | encode_text zscii-text length from coded-text |
| VAR   | 29  | 253  | fd  | 11111101 | V    | ?    | 5   |     |     | copy_table first second size                  |
| VAR   | 30  | 254  | fe  | 11111110 | V    | ?    | 5   |     |     | print_table zscii-text width height skip      |
| VAR   | 31  | 255  | ff  | 11111111 | V    | ?    | 5   |     | \*  | check_arg_count argument-number               |

## EXT opcodes

| Count | Num | Dec. | Hex | Binary   | Form | Args | V    | St  | Br  | Instruction and syntax                          |
| ----- | --- | ---- | --- | -------- | ---- | ---- | ---- | --- | --- | ----------------------------------------------- |
| EXT   | ―   | 190  | be  | 10111110 | E    |      |      |     |     | Extended opcode sentinel value                  |
| EXT   | 0   | 0    | 00  | 00000000 | E    | ?    | 5    | \*  |     | save table bytes name prompt → (result)         |
| EXT   | 1   | 1    | 01  | 00000001 | E    | ?    | 5    | \*  |     | restore table bytes name prompt → (result)      |
| EXT   | 2   | 2    | 02  | 00000010 | E    | ?    | 5    | \*  |     | log_shift number places → (result)              |
| EXT   | 3   | 3    | 03  | 00000011 | E    | ?    | 5/-  | \*  |     | art_shift number places → (result)              |
| EXT   | 4   | 4    | 04  | 00000100 | E    | ?    | 5    | \*  |     | set_font font → (result)                        |
|       |     |      |     |          |      |      | 6/-  | \*  |     | set_font font window → (result)                 |
| EXT   | 5   | 5    | 05  | 00000101 | E    | ?    | 6    |     |     | draw_picture picture-number y x                 |
| EXT   | 6   | 6    | 06  | 00000110 | E    | ?    | 6    |     | \*  | picture_data picture-number array ?(label)      |
| EXT   | 7   | 7    | 07  | 00000111 | E    | ?    | 6    |     |     | erase_picture picture-number y x                |
| EXT   | 8   | 8    | 08  | 00001000 | E    | ?    | 6    |     |     | set_margins left right window                   |
| EXT   | 9   | 9    | 09  | 00001001 | E    | ?    | 5    | \*  |     | save_undo → (result)                            |
| EXT   | 10  | 10   | 0a  | 00001010 | E    | ?    | 5    | \*  |     | restore_undo → (result)                         |
| EXT   | 11  | 11   | 0b  | 00001011 | E    | ?    | 5/\* |     |     | print_unicode char-number                       |
| EXT   | 12  | 12   | 0c  | 00001100 | E    | ?    | 5/\* |     |     | check_unicode char-number → (result)            |
| EXT   | 13  | 13   | 0d  | 00001101 | E    | ?    | 5/\* |     |     | set_true_colour foreground background           |
|       |     |      |     |          |      |      | 6/\* |     |     | set_true_colour foreground background window    |
| EXT   | ―   | 14   | 0e  | 00001110 | E    | ?    |      |     |     | ―                                               |
| EXT   | ―   | 15   | 0f  | 00001111 | E    | ?    |      |     |     | ―                                               |
| EXT   | 16  | 16   | 10  | 00010000 | E    | ?    | 6    |     |     | move_window window y x                          |
| EXT   | 17  | 17   | 11  | 00010001 | E    | ?    | 6    |     |     | window_size window y x                          |
| EXT   | 18  | 18   | 12  | 00010010 | E    | ?    | 6    |     |     | window_style window flags operation             |
| EXT   | 19  | 19   | 13  | 00010011 | E    | ?    | 6    | \*  |     | get_wind_prop window property-number → (result) |
| EXT   | 20  | 20   | 14  | 00010100 | E    | ?    | 6    |     |     | scroll_window window pixels                     |
| EXT   | 21  | 21   | 15  | 00010101 | E    | ?    | 6    |     |     | pop_stack items stack                           |
| EXT   | 22  | 22   | 16  | 00010110 | E    | ?    | 6    |     |     | read_mouse array                                |
| EXT   | 23  | 23   | 17  | 00010111 | E    | ?    | 6    |     |     | mouse_window window                             |
| EXT   | 24  | 24   | 18  | 00011000 | E    | ?    | 6    |     | \*  | push_stack value stack ?(label)                 |
| EXT   | 25  | 25   | 19  | 00011001 | E    | ?    | 6    |     |     | put_wind_prop window property-number value      |
| EXT   | 26  | 26   | 1a  | 00011010 | E    | ?    | 6    |     |     | print_form formatted-table                      |
| EXT   | 27  | 27   | 1b  | 00011011 | E    | ?    | 6    |     | \*  | make_menu number table ?(label)                 |
| EXT   | 28  | 28   | 1c  | 00011100 | E    | ?    | 6    |     |     | picture_table table                             |
| EXT   | 29  | 29   | 1d  | 00011101 | E    | ?    | 6/\* | \*  |     | buffer_screen mode → (result)                   |

---

**Source:** [The Z-Machine Standards Document - Opcodes revisited](https://zspec.jaredreisinger.com/zz03-opcodes)
