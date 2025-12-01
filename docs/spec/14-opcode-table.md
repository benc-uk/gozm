# 14. Complete table of opcodes

## 14.1 Contents

This table contains all 119 opcodes and, taken with the dictionary in [S15](https://zspec.jaredreisinger.com/15-opcodes), describes exactly what each should do. In addition, it lists which opcodes are actually used in the known Infocom story files, and documents the Inform assembly language syntax.

## Opcode Tables

### 2OP Opcodes (0-31)

| St  | Br  | Opcode | Hex | V   | Instruction                              | Name          |
| --- | --- | ------ | --- | --- | ---------------------------------------- | ------------- |
|     |     | —      | 0   | —   | —                                        |               |
|     | \*  | 2OP:1  | 1   |     | je a b ?(label)                          | je            |
|     | \*  | 2OP:2  | 2   |     | jl a b ?(label)                          | jl            |
|     | \*  | 2OP:3  | 3   |     | jg a b ?(label)                          | jg            |
|     | \*  | 2OP:4  | 4   |     | dec_chk (variable) value ?(label)        | dec_chk       |
|     | \*  | 2OP:5  | 5   |     | inc_chk (variable) value ?(label)        | inc_chk       |
|     | \*  | 2OP:6  | 6   |     | jin obj1 obj2 ?(label)                   | jin           |
|     | \*  | 2OP:7  | 7   |     | test bitmap flags ?(label)               | test          |
| \*  |     | 2OP:8  | 8   |     | or a b → (result)                        | or            |
| \*  |     | 2OP:9  | 9   |     | and a b → (result)                       | and           |
|     | \*  | 2OP:10 | A   |     | test_attr object attribute ?(label)      | test_attr     |
|     |     | 2OP:11 | B   |     | set_attr object attribute                | set_attr      |
|     |     | 2OP:12 | C   |     | clear_attr object attribute              | clear_attr    |
|     |     | 2OP:13 | D   |     | store (variable) value                   | store         |
|     |     | 2OP:14 | E   |     | insert_obj object destination            | insert_obj    |
| \*  |     | 2OP:15 | F   |     | loadw array word-index → (result)        | loadw         |
| \*  |     | 2OP:16 | 10  |     | loadb array byte-index → (result)        | loadb         |
| \*  |     | 2OP:17 | 11  |     | get_prop object property → (result)      | get_prop      |
| \*  |     | 2OP:18 | 12  |     | get_prop_addr object property → (result) | get_prop_addr |
| \*  |     | 2OP:19 | 13  |     | get_next_prop object property → (result) | get_next_prop |
| \*  |     | 2OP:20 | 14  |     | add a b → (result)                       | add           |
| \*  |     | 2OP:21 | 15  |     | sub a b → (result)                       | sub           |
| \*  |     | 2OP:22 | 16  |     | mul a b → (result)                       | mul           |
| \*  |     | 2OP:23 | 17  |     | div a b → (result)                       | div           |
| \*  |     | 2OP:24 | 18  |     | mod a b → (result)                       | mod           |
| \*  |     | 2OP:25 | 19  | 4   | call_2s routine arg1 → (result)          | call_2s       |
|     |     | 2OP:26 | 1A  | 5   | call_2n routine arg1                     | call_2n       |
|     |     | 2OP:27 | 1B  | 5   | set_colour foreground background         | set_colour    |
|     |     |        |     | 6   | set_colour foreground background window  | set_colour    |
|     |     | 2OP:28 | 1C  | 5/6 | throw value stack-frame                  | throw         |
|     |     | —      | 1D  | —   | —                                        |               |
|     |     | —      | 1E  | —   | —                                        |               |
|     |     | —      | 1F  | —   | —                                        |               |

_Opcode numbers 32 to 127: other forms of 2OP with different types._

### 1OP Opcodes (128-143)

| St  | Br  | Opcode  | Hex | V   | Instruction                              | Name         |
| --- | --- | ------- | --- | --- | ---------------------------------------- | ------------ |
|     | \*  | 1OP:128 | 0   |     | jz a ?(label)                            | jz           |
| \*  | \*  | 1OP:129 | 1   |     | get_sibling object → (result) ?(label)   | get_sibling  |
| \*  | \*  | 1OP:130 | 2   |     | get_child object → (result) ?(label)     | get_child    |
| \*  |     | 1OP:131 | 3   |     | get_parent object → (result)             | get_parent   |
| \*  |     | 1OP:132 | 4   |     | get_prop_len property-address → (result) | get_prop_len |
|     |     | 1OP:133 | 5   |     | inc (variable)                           | inc          |
|     |     | 1OP:134 | 6   |     | dec (variable)                           | dec          |
|     |     | 1OP:135 | 7   |     | print_addr byte-address-of-string        | print_addr   |
| \*  |     | 1OP:136 | 8   | 4   | call_1s routine → (result)               | call_1s      |
|     |     | 1OP:137 | 9   |     | remove_obj object                        | remove_obj   |
|     |     | 1OP:138 | A   |     | print_obj object                         | print_obj    |
|     |     | 1OP:139 | B   |     | ret value                                | ret          |
|     |     | 1OP:140 | C   |     | jump ?(label)                            | jump         |
|     |     | 1OP:141 | D   |     | print_paddr packed-address-of-string     | print_paddr  |
| \*  |     | 1OP:142 | E   |     | load (variable) → (result)               | load         |
| \*  |     | 1OP:143 | F   | 1/4 | not value → (result)                     | not          |
|     |     |         |     | 5   | call_1n routine                          | call_1n      |

_Opcode numbers 144 to 175: other forms of 1OP with different types._

### 0OP Opcodes (176-191)

| St  | Br  | Opcode  | Hex | V   | Instruction                     | Name        |
| --- | --- | ------- | --- | --- | ------------------------------- | ----------- |
|     |     | 0OP:176 | 0   |     | rtrue                           | rtrue       |
|     |     | 0OP:177 | 1   |     | rfalse                          | rfalse      |
|     |     | 0OP:178 | 2   |     | print (literal-string)          | print       |
|     |     | 0OP:179 | 3   |     | print_ret (literal-string)      | print_ret   |
|     |     | 0OP:180 | 4   | 1/- | nop                             | nop         |
|     | \*  | 0OP:181 | 5   | 1   | save ?(label)                   | save        |
|     |     |         |     | 4   | save → (result)                 | save        |
|     |     |         |     | 5   | illegal                         |             |
|     | \*  | 0OP:182 | 6   | 1   | restore ?(label)                | restore     |
|     |     |         |     | 4   | restore → (result)              | restore     |
|     |     |         |     | 5   | illegal                         |             |
|     |     | 0OP:183 | 7   |     | restart                         | restart     |
|     |     | 0OP:184 | 8   |     | ret_popped                      | ret_popped  |
|     |     | 0OP:185 | 9   | 1   | pop                             | pop         |
| \*  |     |         |     | 5/6 | catch → (result)                | catch       |
|     |     | 0OP:186 | A   |     | quit                            | quit        |
|     |     | 0OP:187 | B   |     | new_line                        | new_line    |
|     |     | 0OP:188 | C   | 3   | show_status                     | show_status |
|     |     |         |     | 4   | illegal                         |             |
|     | \*  | 0OP:189 | D   | 3   | verify ?(label)                 | verify      |
|     |     | 0OP:190 | E   | 5   | [first byte of extended opcode] | extended    |
|     | \*  | 0OP:191 | F   | 5/- | piracy ?(label)                 | piracy      |

_Opcode numbers 192 to 223: VAR forms of 2OP:0 to 2OP:31._

### VAR Opcodes (224-255)

| St  | Br  | Opcode  | Hex | V   | Instruction                                   | Name            |
| --- | --- | ------- | --- | --- | --------------------------------------------- | --------------- |
| \*  |     | VAR:224 | 0   | 1   | call routine …0 to 3 args… → (result)         | call            |
|     |     |         |     | 4   | call_vs routine …0 to 3 args… → (result)      | call_vs         |
|     |     | VAR:225 | 1   |     | storew array word-index value                 | storew          |
|     |     | VAR:226 | 2   |     | storeb array byte-index value                 | storeb          |
|     |     | VAR:227 | 3   |     | put_prop object property value                | put_prop        |
|     |     | VAR:228 | 4   | 1   | sread text parse                              | sread           |
|     |     |         |     | 4   | sread text parse time routine                 | sread           |
| \*  |     |         |     | 5   | aread text parse time routine → (result)      | aread           |
|     |     | VAR:229 | 5   |     | print_char output-character-code              | print_char      |
|     |     | VAR:230 | 6   |     | print_num value                               | print_num       |
| \*  |     | VAR:231 | 7   |     | random range → (result)                       | random          |
|     |     | VAR:232 | 8   |     | push value                                    | push            |
|     |     | VAR:233 | 9   | 1   | pull (variable)                               | pull            |
| \*  |     |         |     | 6   | pull stack → (result)                         | pull            |
|     |     | VAR:234 | A   | 3   | split_window lines                            | split_window    |
|     |     | VAR:235 | B   | 3   | set_window window                             | set_window      |
| \*  |     | VAR:236 | C   | 4   | call_vs2 routine …0 to 7 args… → (result)     | call_vs2        |
|     |     | VAR:237 | D   | 4   | erase_window window                           | erase_window    |
|     |     | VAR:238 | E   | 4/- | erase_line value                              | erase_line      |
|     |     |         |     | 6   | erase_line pixels                             | erase_line      |
|     |     | VAR:239 | F   | 4   | set_cursor line column                        | set_cursor      |
|     |     |         |     | 6   | set_cursor line column window                 | set_cursor      |
|     |     | VAR:240 | 10  | 4/6 | get_cursor array                              | get_cursor      |
|     |     | VAR:241 | 11  | 4   | set_text_style style                          | set_text_style  |
|     |     | VAR:242 | 12  | 4   | buffer_mode flag                              | buffer_mode     |
|     |     | VAR:243 | 13  | 3   | output_stream number                          | output_stream   |
|     |     |         |     | 5   | output_stream number table                    | output_stream   |
|     |     |         |     | 6   | output_stream number table width              | output_stream   |
|     |     | VAR:244 | 14  | 3   | input_stream number                           | input_stream    |
|     |     | VAR:245 | 15  | 5/3 | sound_effect number effect volume routine     | sound_effect    |
| \*  |     | VAR:246 | 16  | 4   | read_char 1 time routine → (result)           | read_char       |
| \*  | \*  | VAR:247 | 17  | 4   | scan_table x table len form → (result)        | scan_table      |
| \*  |     | VAR:248 | 18  | 5/6 | not value → (result)                          | not             |
|     |     | VAR:249 | 19  | 5   | call_vn routine …up to 3 args…                | call_vn         |
|     |     | VAR:250 | 1A  | 5   | call_vn2 routine …up to 7 args…               | call_vn2        |
|     |     | VAR:251 | 1B  | 5   | tokenise text parse dictionary flag           | tokenise        |
|     |     | VAR:252 | 1C  | 5   | encode_text zscii-text length from coded-text | encode_text     |
|     |     | VAR:253 | 1D  | 5   | copy_table first second size                  | copy_table      |
|     |     | VAR:254 | 1E  | 5   | print_table zscii-text width height skip      | print_table     |
|     | \*  | VAR:255 | 1F  | 5   | check_arg_count argument-number               | check_arg_count |

### EXT (Extended) Opcodes

| St  | Br  | Opcode | Hex | V    | Instruction                                     | Name            |
| --- | --- | ------ | --- | ---- | ----------------------------------------------- | --------------- |
| \*  |     | EXT:0  | 0   | 5    | save table bytes name prompt → (result)         | save            |
| \*  |     | EXT:1  | 1   | 5    | restore table bytes name prompt → (result)      | restore         |
| \*  |     | EXT:2  | 2   | 5    | log_shift number places → (result)              | log_shift       |
| \*  |     | EXT:3  | 3   | 5/-  | art_shift number places → (result)              | art_shift       |
| \*  |     | EXT:4  | 4   | 5    | set_font font → (result)                        | set_font        |
| \*  |     |        |     | 6/-  | set_font font window → (result)                 | set_font        |
|     |     | EXT:5  | 5   | 6    | draw_picture picture-number y x                 | draw_picture    |
|     | \*  | EXT:6  | 6   | 6    | picture_data picture-number array ?(label)      | picture_data    |
|     |     | EXT:7  | 7   | 6    | erase_picture picture-number y x                | erase_picture   |
|     |     | EXT:8  | 8   | 6    | set_margins left right window                   | set_margins     |
| \*  |     | EXT:9  | 9   | 5    | save_undo → (result)                            | save_undo       |
| \*  |     | EXT:10 | A   | 5    | restore_undo → (result)                         | restore_undo    |
|     |     | EXT:11 | B   | 5/\* | print_unicode char-number                       | print_unicode   |
|     |     | EXT:12 | C   | 5/\* | check_unicode char-number → (result)            | check_unicode   |
|     |     | EXT:13 | D   | 5/\* | set_true_colour foreground background           | set_true_colour |
|     |     |        |     | 6/\* | set_true_colour foreground background window    | set_true_colour |
|     |     | —      | E   | —    | —                                               |                 |
|     |     | —      | F   | —    | —                                               |                 |
|     |     | EXT:16 | 10  | 6    | move_window window y x                          | move_window     |
|     |     | EXT:17 | 11  | 6    | window_size window y x                          | window_size     |
|     |     | EXT:18 | 12  | 6    | window_style window flags operation             | window_style    |
| \*  |     | EXT:19 | 13  | 6    | get_wind_prop window property-number → (result) | get_wind_prop   |
|     |     | EXT:20 | 14  | 6    | scroll_window window pixels                     | scroll_window   |
|     |     | EXT:21 | 15  | 6    | pop_stack items stack                           | pop_stack       |
|     |     | EXT:22 | 16  | 6    | read_mouse array                                | read_mouse      |
|     |     | EXT:23 | 17  | 6    | mouse_window window                             | mouse_window    |
|     | \*  | EXT:24 | 18  | 6    | push_stack value stack ?(label)                 | push_stack      |
|     |     | EXT:25 | 19  | 6    | put_wind_prop window property-number value      | put_wind_prop   |
|     |     | EXT:26 | 1A  | 6    | print_form formatted-table                      | print_form      |
|     | \*  | EXT:27 | 1B  | 6    | make_menu number table ?(label)                 | make_menu       |
|     |     | EXT:28 | 1C  | 6    | picture_table table                             | picture_table   |
| \*  |     | EXT:29 | 1D  | 6/\* | buffer_screen mode → (result)                   | buffer_screen   |

## 14.2 Out of range opcodes

Formally, it is illegal for a game to contain an opcode not specified for its version. An interpreter should normally halt with a suitable message.

### 14.2.1

However, extended opcodes in the range EXT:29 to EXT:255 should be simply ignored (perhaps with a warning message somewhere off-screen).

### 14.2.2

**[1.0][1.1]** EXT:11 and EXT:12 were opcodes added in Standard 1.0 and can be generated in code compiled by Inform 6.12 or later. EXT:13 and EXT:29 are new in Standard 1.1. EXT:14 to EXT:15, and EXT:30 to EXT:127, are reserved for future versions of this document to specify.

### 14.2.3

Designers who wish to create their own "new" opcodes, for one specific game only, are asked to use opcode numbers in the range EXT:128 to EXT:255. It is easy to modify Inform to name and assemble such opcodes. (Of course the game will then have to be circulated with a suitably modified interpreter to run it.)

### 14.2.4

Interpreter-writers should ideally make this easy by providing a routine which is called if EXT:128 to EXT:255 are found, so that the minimum possible modification to the interpreter is needed.

## Reading the opcode tables

The two columns St and Br (store and branch) mark whether an instruction stores a result in a variable, and whether it must provide a label to jump to, respectively.

The Opcode is written TYPE:Decimal where the TYPE is the operand count (2OP, 1OP, 0OP or VAR) or else EXT for two-byte opcodes (where the first byte is (decimal) 190). The decimal number is the lowest possible decimal opcode value. The hex number is the opcode number within each TYPE.

The V column gives the Version information. If nothing is specified, the opcode is as stated from Version 1 onwards. Otherwise, it exists only from the version quoted onwards. Before this time, its use is illegal. Some opcodes change their meanings as the Version increases, and these have more than one line of specification. Others become illegal again, and these are marked "illegal".

In a few cases, the Version is given as "3/4″ or some such. The first number is the Version number whose specification the opcode belongs to, and the second is the earliest Version in which the opcode is known actually to be used in an Infocom-produced story file. A dash means that it seems never to have been used (in any of Versions 1 to 6). The notation "5/-" or "6/-" means that the opcode was introduced in this Standards document long after the Infocom era.

The table explicitly marks opcodes which do not exist in any version of the Z-machine as "―": in addition, none of the extended set of codes after EXT:29 have been used.

## Inform assembly language

This section documents Inform 6 assembly language, which is richer than that of Inform 5. The Inform 6 assembler can generate every legal opcode and automatically sets any consequent header bits (for instance, a usage of `set_colour` will set the "colours needed" bit).

One way to get a picture of Inform assembly language is to compile a short program with tracing switched on (using the `-a` or `-t` switches).

1. An Inform statement beginning with an `@` is sent directly to the assembler. In the syntax below, `(variable)` and `(result)` must be variables (or `sp`, a special variable name available only in assembly language, and meaning the stack pointer); `(label)` a label (not a routine name). `(literal-string)` must be literal text in quotation marks `"thus"`. `routine` should be the name of a routine (this assembles to its packed address). Otherwise any Inform constant term (such as `/` or `beetle`) can be given as an operand.

2. It is optional, but sensible, to place a `→` sign before a store-variable. For example, in

   ```
   @mul a 56 -> sp;
   ```

   ("multiply variable `a` by 56, and put the result on the stack") the → can be omitted, but should be included for clarity.

3. A label to branch to should be prefaced with a question mark `?`, as in

   ```
   @je a b ?Equal;      ! Branch to Equal if a == b
   ```

   (If the question mark is omitted, the branch is compiled in the short form, which will only work for very nearby labels and is very seldom useful in code written by hand.) Note that the effect of any branch instruction can be negated using a tilde ~:

   ```
   @je a b ?~Different; ! Branch to Different if a ~= b
   ```

4. Labels are assembled using full stops:

   ```
   .MyLabel;
   ```

   All branches must be to such a label within the same routine. (The Inform assembler imposes the same-routine restriction.)

5. Most operands are assembled in the obvious way: numbers and constant values (like characters) as numbers, variables as variables, `sp` as the value on top of the stack. There are two exceptions. "Call" opcodes expect as first operand the name of a routine to call:

   ```
   @call_1n MyRoutine;
   ```

   but one can also give an indirect address, as a constant or variable, using square brackets:

   ```
   @call_1n [x];        ! Call routine whose address is in x
   ```

   Secondly, seven Z-machine opcodes access variables but by their numbers: thus one should write, say, the constant 0 instead of the variable `sp`. This is inconvenient, so the Inform assembler accepts variable names instead. The operands affected are those marked as `(variable)` in the syntax chart; Inform translates the variable name as a "small constant" operand with that variable's number as value. The affected opcodes are: `inc`, `dec`, `inc_chk`, `dec_chk`, `store`, `pull`, `load`. This is useful, but there is another possibility, of genuinely giving a variable operand. The Inform notation for this involves square brackets again:

   ```
   @inc frog;          ! Increment var "frog"
   @inc [frog];        ! Increment var whose number is in "frog"
   ```

   Infocom story files often use such instructions.

6. The Inform assembler is also written with possible extensions to the Z-machine instruction set in mind. (Of course these can only work if a customised interpreter is used.) Simply give a specification in double-quotes where you would normally give the opcode name. For example,

   ```
   @"1OP:4S" 34 -> i;
   @get_prop_len 34 -> i;
   ```

   are equivalent instructions, since `get_prop_len` is instruction 4 in the 1OP (one-operand) set, and is a Store opcode. The syntax is:

   ```
   "0OP:decimal-number flags"     (range 0 to 15)
   "1OP:decimal-number flags"     (range 0 to 15)
   "2OP:decimal-number flags"     (range 0 to 15)
   "VAR:decimal-number flags"     (range 32 to 63)
   "VAR_LONG:decimal-number flags" (range 32 to 63)
   "EXT:decimal-number flags"     (range 0 to 255)
   "EXT_LONG:decimal-number flags" (range 0 to 255)
   ```

   (EXT_LONG is a logical possibility but has not been used in the Z-machine so far: the assembler provides it in case it might be useful in future.) The possible flags are:

   | Flag | Meaning                                                                                                                     |
   | ---- | --------------------------------------------------------------------------------------------------------------------------- |
   | S    | Store opcode                                                                                                                |
   | B    | Branch opcode                                                                                                               |
   | T    | Text in-line instead of operands (as with print and print_ret)                                                              |
   | I    | "Indirect addressing": first operand is a (variable)                                                                        |
   | Fnn  | Set bit nn in Flags 2 (signalling to the interpreter that an unusual feature has been called for): the number is in decimal |

   For example, `EXT:128BSF14` is an exotic new opcode, number 128 in the extended range, which is both Branch and Store, and the assembly of which causes bit 14 to be set in "Flags 2". See [S14.2](#142-out-of-range-opcodes) for rules on how to number newly created opcodes.

## Remarks

The opcodes EXT:5 to EXT:8 were very likely in Infocom's own Version 5 specification (documentary records of which are lost): they seem to have been partially implemented in existing Infocom interpreters, but do not occur in any existing Version 5 story file. They are here left unspecified.

The notation "5/3" for `sound_effect` is because this plainly Version 5 feature was used also in one solitary Version 3 game, The Lurking Horror (the sound version of which was the last Version 3 release, in September 1987).

The 2OP opcode 0 was possibly intended for setting break-points in debugging (and may be used for this again). It was not `nop`.

`read_mouse` and `make_menu` are believed to have been used only in Journey (based on a check of 11 Version 6 story files).

`picture_table` is used once by Shogun and several times by Zork Zero.

---

_Source: [Z-Machine Standards Document - Section 14: Complete table of opcodes](https://zspec.jaredreisinger.com/14-opcode-table)_
