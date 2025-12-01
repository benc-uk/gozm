# 15. Dictionary of opcodes

> The highest ideal of a translation… is achieved when the reader flings it impatiently into the fire, and begins patiently to learn the language for himself.
>
> — Philip Vellacott

## 15.1

The dictionary below is alphabetical and includes entries on every opcode listed in the table above, as well as brief notes on a few opcodes once thought to exist but now disproved.

## 15.2

The Z-machine has the same concept of "table" (as an internal data structure) as Inform. Specifically, a table is an array of words (in dynamic or static memory) of which the initial entry is the number of subsequent words in the table. For example, a table with three entries occupies 8 bytes, arranged as the words 3, x, y, z.

## 15.3

In all cases below where one operand is supposed to be in a particular range, behaviour is undefined if it is not. For instance an interpreter complies with the Standard even if it crashes when an illegal object number (including 0) is given for an object operand. However, see [SA](https://zspec.jaredreisinger.com/A-errors) for guidelines on detecting and dealing with errors.

---

## Opcodes

### `add`

**2OP:20 14** `add a b → (result)`

Signed 16-bit addition.

### `and`

**2OP:9 9** `and a b → (result)`

Bitwise AND.

### `aread`

This is the Inform name for the keyboard-reading opcode under Version 5 and later. (Inform calls the same opcode sread under Versions 3 and 4.) See `read` for the specification.

### `art_shift`

**EXT:3 3** `5/- art_shift number places → (result)`

Does an arithmetic shift of number by the given number of places, shifting left (i.e. increasing) if places is positive, right if negative. In a right shift, the sign bit is preserved as well as being shifted on down. (The alternative behaviour is `log_shift`.)

The places operand must be in the range -15 to +15, otherwise behaviour is undefined.

### `buffer_mode`

**VAR:242 12** `4 buffer_mode flag`

If set to 1, text output on the lower window in stream 1 is buffered up so that it can be word-wrapped properly. If set to 0, it isn't.

In Version 6, this opcode is redundant (the "buffering" window attribute can be set instead). It is used twice in each of Infocom's Version 6 story files, in the `$verify` routine. Frotz responds by setting the current window's "buffering" attribute, while Infocom's own interpreters respond by doing nothing. This standard leaves the result of `buffer_mode` undefined in Version 6.

### `buffer_screen`

**EXT:29 1D** `6/* buffer_screen mode → (result)`

Tells the interpreter how to handle display buffering. If mode is 0, updates must be made as soon as possible. If mode is 1, the interpreter may make changes to a backing store, and need not update the screen. The interpreter is free to ignore the advice, but if so must always act as though the mode is 0 (update the screen as soon as possible).

With `buffer_screen` in either state, an update of the visible display can be forced immediately by issuing `buffer_screen -1`, without altering the current buffering state. Note that `buffer_screen -1` does not flush the text buffer.

The return value is the old `buffer_screen` state.

See [S8](https://zspec.jaredreisinger.com/08-screen) for more details.

**[1.1]** This opcode will only be present in interpreters obeying Standard 1.1 or later, so story files should check the standard number of the interpreter before executing this opcode.

### `call`

**VAR:224 0** `1 call routine …up to 3 args… → (result)`

The only call instruction in Version 3, Inform reads this as `call_vs` in higher versions: it calls the routine with 0, 1, 2 or 3 arguments as supplied and stores the resulting return value. (When the address 0 is called as a routine, nothing happens and the return value is false.)

### `call_1n`

**1OP:143 F** `5 call_1n routine`

Executes routine() and throws away result.

### `call_1s`

**1OP:136 8** `4 call_1s routine → (result)`

Stores routine().

### `call_2n`

**2OP:26 1A** `5 call_2n routine arg1`

Executes routine(arg1) and throws away result.

### `call_2s`

**2OP:25 19** `4 call_2s routine arg1 → (result)`

Stores routine(arg1).

### `call_vn`

**VAR:249 19** `5 call_vn routine …up to 3 args…`

Like call, but throws away result.

### `call_vs`

**VAR:224 0** `4 call_vs routine …up to 3 args… → (result)`

See `call`.

### `call_vn2`

**VAR:250 1A** `5 call_vn2 routine …up to 7 args…`

Call with a variable number (from 0 to 7) of arguments, then throw away the result. This (and `call_vs2`) uniquely have an extra byte of opcode types to specify the types of arguments 4 to 7. Note that it is legal to use these opcodes with fewer than 4 arguments (in which case the second byte of type information will just be `$ff`).

### `call_vs2`

**VAR:236 C** `4 call_vs2 routine …up to 7 args… → (result)`

See `call_vn2`.

### `catch`

**0OP:185 9** `5/6 catch → (result)`

Opposite to throw (and occupying the same opcode that pop used in Versions 3 and 4). `catch` returns the current "stack frame".

### `check_arg_count`

**VAR:255 1F** `5 check_arg_count argument-number`

Branches if the given argument-number (counting from 1) has been provided by the routine call to the current routine. (This allows routines in Versions 5 and later to distinguish between the calls routine(1) and routine(1,0), which would otherwise be impossible to tell apart.)

### `check_unicode`

**EXT:12 C** `5/* check_unicode char-number → (result)`

Determines whether or not the interpreter can print, or receive from the keyboard, the given Unicode character. Bit 0 of the result should be set if and only if the interpreter can print the character; bit 1 if and only if the interpreter can receive it from the keyboard. Bits 2 to 15 are undefined.

**[1.0]** This opcode will only be present in interpreters obeying Standard 1.0 or later, so story files should check the standard number of the interpreter before executing this opcode.

### `clear_attr`

**2OP:12 C** `clear_attr object attribute`

Make object not have the attribute numbered attribute.

### `copy_table`

**VAR:253 1D** `5 copy_table first second size`

If second is zero, then size bytes of first are zeroed.

Otherwise first is copied into second, its length in bytes being the absolute value of size (i.e., size if size is positive, -size if size is negative).

The tables are allowed to overlap. If size is positive, the interpreter must copy either forwards or backwards so as to avoid corrupting first in the copying process. If size is negative, the interpreter must copy forwards even if this corrupts first. (Beyond Zork uses this to fill an array with spaces.)

### `dec`

**1OP:134 6** `dec (variable)`

Decrement variable by 1. This is signed, so 0 decrements to -1.

### `dec_chk`

**2OP:4 4** `dec_chk (variable) value ?(label)`

Decrement variable, and branch if it is now less than the given value.

### `div`

**2OP:23 17** `div a b → (result)`

Signed 16-bit division. Division by zero should halt the interpreter with a suitable error message.

### `draw_picture`

**EXT:5 5** `6 draw_picture picture-number y x`

Displays the picture with the given number. (y,x) coordinates (of the top left of the picture) are each optional, in that a value of zero for y or x means the cursor y or x coordinate in the current window. It is illegal to call this with an invalid picture number.

### `encode_text`

**VAR:252 1C** `5 encode_text zscii-text length from coded-text`

Translates a ZSCII word to Z-encoded text format (stored at coded-text), as if it were an entry in the dictionary. The text begins at from in the zscii-text buffer and is length characters long.

### `erase_line`

**VAR:238 E** `4/6 erase_line value`

**Versions 4 and 5:** if the value is 1, erase from the current cursor position to the end of its line in the current window. If the value is anything other than 1, do nothing.

**Version 6:** if the value is 1, erase from the current cursor position to the end of the its line in the current window. If not, erase the given number of pixels minus one across from the cursor (clipped to stay inside the right margin). The cursor does not move.

### `erase_picture`

**EXT:7 7** `6 erase_picture picture-number y x`

Like `draw_picture`, but paints the appropriate region to the background colour for the given window. It is illegal to call this with an invalid picture number.

### `erase_window`

**VAR:237 D** `4 erase_window window`

Erases window with given number (to background colour); or if -1 it unsplits the screen and clears the lot; or if -2 it clears the screen without unsplitting it. In cases -1 and -2, the cursor may move (see [S8](https://zspec.jaredreisinger.com/08-screen) for precise details).

### "extended"

This byte (decimal 190) is not an instruction, but indicates that the opcode is "extended": the next byte contains the number in the extended set.

### `get_child`

**1OP:130 2** `get_child object → (result) ?(label)`

Get first object contained in given object, branching if this exists, i.e. is not nothing (i.e., is not 0).

### `get_cursor`

**VAR:240 10** `4/6 get_cursor array`

Puts the current cursor row into the word 0 of the given array, and the current cursor column into word 1. (The array is not a table and has no size information in its initial entry.)

### `get_next_prop`

**2OP:19 13** `get_next_prop object property → (result)`

Gives the number of the next property provided by the quoted object. This may be zero, indicating the end of the property list; if called with zero, it gives the first property number present. It is illegal to try to find the next property of a property which does not exist, and an interpreter should halt with an error message (if it can efficiently check this condition).

### `get_parent`

**1OP:131 3** `get_parent object → (result)`

Get parent object (note that this has no "branch if exists" clause).

### `get_prop`

**2OP:17 11** `get_prop object property → (result)`

Read property from object (resulting in the default value if it had no such declared property). If the property has length 1, the value is only that byte. If it has length 2, the first two bytes of the property are taken as a word value. It is illegal for the opcode to be used if the property has length greater than 2, and the result is unspecified.

### `get_prop_addr`

**2OP:18 12** `get_prop_addr object property → (result)`

Get the byte address (in dynamic memory) of the property data for the given object's property. This must return 0 if the object hasn't got the property.

Note that the address retrieved should be the address after the size byte(s) in the property entry.

### `get_prop_len`

**1OP:132 4** `get_prop_len property-address → (result)`

Get length of property data (in bytes) for the given object's property. It is illegal to try to find the property length of a property which does not exist for the given object, and an interpreter should halt with an error message (if it can efficiently check this condition).

`@get_prop_len 0` must return 0. This is required by some Infocom games and files generated by old versions of Inform.

### `get_sibling`

**1OP:129 1** `get_sibling object → (result) ?(label)`

Get next object in tree, branching if this exists, i.e. is not 0.

### `get_wind_prop`

**EXT:19 13** `6 get_wind_prop window property-number → (result)`

Reads the given property of the given window (see [S8](https://zspec.jaredreisinger.com/08-screen)).

### `inc`

**1OP:133 5** `inc (variable)`

Increment variable by 1. (This is signed, so -1 increments to 0.)

### `inc_chk`

**2OP:5 5** `inc_chk (variable) value ?(label)`

Increment variable, and branch if now greater than value.

### `input_stream`

**VAR:244 14** `3 input_stream number`

Selects the current input stream.

### `insert_obj`

**2OP:14 E** `insert_obj object destination`

Moves object O to become the first child of the destination object D. (Thus, after the operation the child of D is O, and the sibling of O is whatever was previously the child of D.) All children of O move with it. (Initially O can be at any point in the object tree; it may legally have parent zero.)

### `je`

**2OP:1 1** `je a b c d ?(label)`

Jump if a is equal to any of the subsequent operands. (Thus `@je a` never jumps and `@je a b` jumps if a = b.)

`je` with just 1 operand is not permitted.

### `jg`

**2OP:3 3** `jg a b ?(label)`

Jump if a > b (using a signed 16-bit comparison).

### `jin`

**2OP:6 6** `jin obj1 obj2 ?(label)`

Jump if object a is a direct child of b, i.e., if parent of a is b.

### `jl`

**2OP:2 2** `jl a b ?(label)`

Jump if a < b (using a signed 16-bit comparison).

### `jump`

**1OP:140 C** `jump ?(label)`

Jump (unconditionally) to the given label. (This is not a branch instruction and the operand is a 2-byte signed offset to apply to the program counter.) It is legal for this to jump into a different routine (which should not change the routine call state), although it is considered bad practice to do so.

The destination of the jump opcode is: `Address after instruction + Offset - 2`

### `jz`

**1OP:128 0** `jz a ?(label)`

Jump if a = 0.

### `load`

**1OP:142 E** `load (variable) → (result)`

The value of the variable referred to by the operand is stored in the result.

### `loadb`

**2OP:16 10** `loadb array byte-index → (result)`

Stores `array→byte-index` (i.e., the byte at address array+byte-index, which must lie in static or dynamic memory).

### `loadw`

**2OP:15 F** `loadw array word-index → (result)`

Stores `array→word-index` (i.e., the word at address array+2\*word-index, which must lie in static or dynamic memory).

### `log_shift`

**EXT:2 2** `5 log_shift number places → (result)`

Does a logical shift of number by the given number of places, shifting left (i.e. increasing) if places is positive, right if negative. In a right shift, the sign is zeroed instead of being shifted on. (See also `art_shift`.)

The places operand must be in the range -15 to +15, otherwise behaviour is undefined.

### `make_menu`

**EXT:27 1B** `6 make_menu number table ?(label)`

Controls menus with numbers greater than 2. If the table supplied is 0, the menu is removed. Otherwise it is a table of tables. Each table is a ZSCII string: the first item being a menu name, subsequent ones the entries.

### `mod`

**2OP:24 18** `mod a b → (result)`

Remainder after signed 16-bit division. Division by zero should halt the interpreter with a suitable error message.

### `mouse_window`

**EXT:23 17** `6 mouse_window window`

Constrain the mouse arrow to sit inside the given window. By default it sits in window 1. Setting to -1 takes all restriction away.

### `move_window`

**EXT:16 10** `6 move_window window y x`

Moves the given window to pixels (y,x): (1,1) being the top left.

### `mul`

**2OP:22 16** `mul a b → (result)`

Signed 16-bit multiplication.

### `new_line`

**0OP:187 B** `new_line`

Print carriage return.

### `nop`

**0OP:180 4** `1/- nop`

Probably the official "no operation" instruction, which, appropriately, was never operated (in any of the Infocom datafiles): it may once have been a breakpoint.

### `not`

**1OP:143 F** `1/4 not value → (result)`  
**VAR:248 18** `5/6 not value → (result)`

Bitwise NOT (i.e., all 16 bits reversed). Note that in Versions 3 and 4 this is a 1OP instruction, but in later Versions it was moved into the extended set to make room for `call_1n`.

### `or`

**2OP:8 8** `or a b → (result)`

Bitwise OR.

### `output_stream`

**VAR:243 13** `3 output_stream number`  
**VAR:243 13** `5 output_stream number table`  
**VAR:243 13** `6 output_stream number table width`

If stream is 0, nothing happens. If it is positive, then that stream is selected; if negative, deselected.

When stream 3 is selected, a table must be given into which text can be printed. The first word always holds the number of characters printed, the actual text being stored at bytes table+2 onward.

In Version 6, a width field may optionally be given: text will then be justified as if it were in the window with that number (if width is zero or positive) or a box -width pixels wide (if negative).

### `picture_data`

**EXT:6 6** `6 picture_data picture-number array ?(label)`

Asks the interpreter for data on the picture with the given number. If the picture number is valid, a branch occurs and information is written to the array: the height in word 0, the width in word 1, in pixels.

Otherwise, if the picture number is zero, the interpreter writes the number of available pictures into word 0 of the array and the release number of the picture file into word 1, and branches if any pictures are available.

### `picture_table`

**EXT:28 1C** `6 picture_table table`

Given a table of picture numbers, the interpreter may if it wishes load or unpack these pictures from disc into a cache for convenient rapid plotting later.

### `piracy`

**0OP:191 F** `5/- piracy ?(label)`

Branches if the game disc is believed to be genuine by the interpreter. Interpreters are asked to be gullible and to unconditionally branch.

### `pop`

**0OP:185 9** `1 pop`

Throws away the top item on the stack.

### `pop_stack`

**EXT:21 15** `6 pop_stack items stack`

The given number of items are thrown away from the top of a stack: by default the system stack, otherwise the one given as a second operand.

### `print`

**0OP:178 2** `print <literal-string>`

Print the quoted (literal) Z-encoded string.

### `print_addr`

**1OP:135 7** `print_addr byte-address-of-string`

Print (Z-encoded) string at given byte address, in dynamic or static memory.

### `print_char`

**VAR:229 5** `print_char output-character-code`

Print a ZSCII character. The operand must be a character code defined in ZSCII for output (see [S3](https://zspec.jaredreisinger.com/03-text)).

### `print_form`

**EXT:26 1A** `6 print_form formatted-table`

Prints a formatted table of the kind written to output stream 3 when formatting is on.

### `print_num`

**VAR:230 6** `print_num value`

Print (signed) number in decimal.

### `print_obj`

**1OP:138 A** `print_obj object`

Print short name of object (the Z-encoded string in the object header, not a property).

### `print_paddr`

**1OP:141 D** `print_paddr packed-address-of-string`

Print the (Z-encoded) string at the given packed address in high memory.

### `print_ret`

**0OP:179 3** `print_ret <literal-string>`

Print the quoted (literal) Z-encoded string, then print a new-line and then return true (i.e., 1).

### `print_table`

**VAR:254 1E** `5 print_table zscii-text width height skip`

Print a rectangle of text on screen spreading right and down from the current cursor position, of given width and height, from the table of ZSCII text given.

### `print_unicode`

**EXT:11 B** `5/* print_unicode char-number`

Print a Unicode character. The given character code must be defined in Unicode.

**[1.0]** This opcode will only be present in interpreters obeying Standard 1.0 or later, so story files should check the standard number of the interpreter before executing this opcode.

### `pull`

**VAR:233 9** `1 pull (variable)`  
**VAR:233 9** `6 pull stack → (result)`

Pulls value off a stack. In Version 6, the stack in question may be specified as a user one; otherwise it is the game stack.

### `push`

**VAR:232 8** `push value`

Pushes value onto the game stack.

### `push_stack`

**EXT:24 18** `6 push_stack value stack ?(label)`

Pushes the value onto the specified user stack, and branching if this was successful. If the stack overflows, nothing happens.

### `put_prop`

**VAR:227 3** `put_prop object property value`

Writes the given value to the given property of the given object. If the property does not exist for that object, the interpreter should halt with a suitable error message. If the property length is 1, then the interpreter should store only the least significant byte of the value.

### `put_wind_prop`

**EXT:25 19** `6 put_wind_prop window property-number value`

Writes a window property (see `get_wind_prop`).

### `quit`

**0OP:186 A** `quit`

Exit the game immediately. (Any "Are you sure?" question must be asked by the game, not the interpreter.)

### `random`

**VAR:231 7** `random range → (result)`

If range is positive, returns a uniformly random number between 1 and range. If range is negative, the random number generator is seeded to that value and the return value is 0. Most interpreters consider giving 0 as range illegal, but correct behaviour is to reseed the generator in as random a way as the interpreter can.

### `read`

**VAR:228 4** `1 sread text parse`  
**VAR:228 4** `4 sread text parse time routine`  
**VAR:228 4** `5 aread text parse time routine → (result)`

This opcode reads a whole command from the keyboard (no prompt is automatically displayed). It is legal for this to be called with the cursor at any position on any window.

In Versions 1 to 3, the status line is automatically redisplayed first.

**In Versions 1 to 4:** byte 0 of the text-buffer should initially contain the maximum number of letters which can be typed, minus 1. The text typed is reduced to lower case and stored in bytes 1 onward, with a zero terminator.

**In Versions 5 and later:** byte 0 of the text-buffer should initially contain the maximum number of letters which can be typed. The interpreter stores the number of characters actually typed in byte 1, and the characters themselves (reduced to lower case) in bytes 2 onward.

In Version 4 and later, if the operands time and routine are supplied (and non-zero) then the routine call `routine()` is made every time/10 seconds during the keyboard-reading process.

Next, lexical analysis is performed on the text (except that in Versions 5 and later, if parse-buffer is zero then this is omitted).

In Version 5 and later, this is a store instruction: the return value is the terminating character.

### `read_char`

**VAR:246 16** `4 read_char 1 time routine → (result)`

Reads a single character from input stream 0 (the keyboard). The first operand must be 1. time and routine are optional and dealt with as in `read` above.

### `read_mouse`

**EXT:22 16** `6 read_mouse array`

The four words in the array are written with the mouse y coordinate, x coordinate, button bits, and a menu word.

### `remove_obj`

**1OP:137 9** `remove_obj object`

Detach the object from its parent, so that it no longer has any parent. (Its children remain in its possession.)

### `restart`

**0OP:183 7** `1 restart`

Restart the game. The only pieces of information surviving from the previous state are the "transcribing to printer" bit and the "use fixed pitch font" bit.

### `restore`

**0OP:182 6** `1 restore ?(label)`  
**0OP:182 6** `4 restore → (result)`  
**EXT:1 1** `5 restore table bytes name prompt → (result)`

See `save`. In Version 3, the branch is never actually made, since either the game has successfully picked up again from where it was saved, or it failed to load the save game file.

From Version 5 it can have optional parameters as `save` does, and returns the number of bytes loaded if so.

### `restore_undo`

**EXT:10 A** `5 restore_undo → (result)`

Like restore, but restores the state saved to memory by `save_undo`.

### `ret`

**1OP:139 B** `ret value`

Returns from the current routine with the value given.

### `ret_popped`

**0OP:184 8** `ret_popped`

Pops top of stack and returns that.

### `rfalse`

**0OP:177 1** `rfalse`

Return false (i.e., 0) from the current routine.

### `rtrue`

**0OP:176 0** `rtrue`

Return true (i.e., 1) from the current routine.

### `save`

**0OP:181 5** `1 save ?(label)`  
**0OP:181 5** `4 save → (result)`  
**EXT:0 0** `5 save table bytes name prompt → (result)`

On Versions 3 and 4, attempts to save the game and branches if successful. From Version 5 it is a store rather than a branch instruction; the store value is 0 for failure, 1 for "save succeeded" and 2 for "the game is being restored and is resuming execution again from here".

It is illegal to use this opcode within an interrupt routine.

**[1.0]** The extension also has (optional) parameters, which save a region of the save area, whose address and length are in bytes, and provides a suggested filename.

**[1.1]** As of Standard 1.1 an additional optional parameter, prompt, is allowed on Version 5 extended save/restore.

### `save_undo`

**EXT:9 9** `5 save_undo → (result)`

Like save, except that the optional parameters may not be specified: it saves the game into a cache of memory held by the interpreter. If the interpreter is unable to provide this feature, it must return -1: otherwise it returns the save return value.

It is illegal to use this opcode within an interrupt routine.

### `scan_table`

**VAR:247 17** `4 scan_table x table len form → (result)`

Is x one of the words in table, which is len words long? If so, return the address where it first occurs and branch. If not, return 0 and don't.

The form is optional: bit 7 is set for words, clear for bytes: the rest contains the length of each field in the table. Thus `$82` is the default.

### `scroll_window`

**EXT:20 14** `6 scroll_window window pixels`

Scrolls the given window by the given number of pixels (a negative value scrolls backwards, i.e., down) writing in blank (background colour) pixels in the new lines.

### `set_attr`

**2OP:11 B** `set_attr object attribute`

Make object have the attribute numbered attribute.

### `set_colour`

**2OP:27 1B** `5 set_colour foreground background`  
**2OP:27 1B** `6 set_colour foreground background window`

If coloured text is available, set text to be foreground-against-background. In version 6, the window argument is optional and is by default the current window.

### `set_cursor`

**VAR:239 F** `4 set_cursor line column`  
**VAR:239 F** `6 set_cursor line column window`

Move cursor in the current window to the position (x,y) (in units) relative to (1,1) in the top left. In Version 6, `set_cursor -1` turns the cursor off, and either `set_cursor -2` or `set_cursor -2 0` turn it back on.

### `set_font`

**EXT:4 4** `5 set_font font → (result)`  
**EXT:4 4** `6 set_font font window → (result)`

If the requested font is available, then it is chosen for the current window, and the store value is the font ID of the previous font (which is always positive). If the font is unavailable, nothing will happen and the store value is 0.

If the font ID requested is 0, the font is not changed, and the ID of the current font is returned.

**[1.1]** In Version 6, `set_font` has an optional window parameter, as for `set_colour`.

### `set_margins`

**EXT:8 8** `6 set_margins left right window`

Sets the margin widths (in pixels) on the left and right for the given window (which are by default 0).

### `set_text_style`

**VAR:241 11** `4 set_text_style style`

Sets the text style to: Roman (if 0), Reverse Video (if 1), Bold (if 2), Italic (4), Fixed Pitch (8). In some interpreters a combination of styles is possible.

**[1.1]** As of Standard 1.1, it is legal to request style combinations in a single `set_text_style` opcode by adding the values (which are powers of two) together.

### `set_true_colour`

**EXT:13 D** `5/* set_true_colour foreground background`  
**EXT:13 D** `6/* set_true_colour foreground background window`

The foreground and background are 15-bit colour values with bits 14-10 for blue, 9-5 for green, and 4-0 for red.

**[1.1]** This opcode will only be present in interpreters obeying Standard 1.1 or later, so story files should check the standard number of the interpreter before executing this opcode.

### `set_window`

**VAR:235 B** `3 set_window window`

Selects the given window for text output.

### `show_status`

**0OP:188 C** `3 show_status`

(In Version 3 only.) Display and update the status line now.

### `sound_effect`

**VAR:245 15** `5/3 sound_effect number effect volume routine`

The given effect happens to the given sound number. The low byte of volume holds the volume level, the high byte the number of repeats.

Note that sound effect numbers 1 and 2 are bleeps (see [S9](https://zspec.jaredreisinger.com/09-sound)).

The effect can be: 1 (prepare), 2 (start), 3 (stop), 4 (finish with).

In Versions 5 and later, the routine is called (with no parameters) after the sound has been finished.

### `split_window`

**VAR:234 A** `3 split_window lines`

Splits the screen so that the upper window has the given number of lines: or, if this is zero, unsplits the screen again. In Version 3 (only) the upper window should be cleared after the split.

### `sread`

This is the Inform name for the keyboard-reading opcode under Versions 3 and 4. (Inform calls the same opcode aread in later Versions.) See `read` for the specification.

### `store`

**2OP:13 D** `store (variable) value`

Set the VARiable referenced by the operand to value.

### `storeb`

**VAR:226 2** `storeb array byte-index value`

`array→byte-index = value`, i.e. stores the given value in the byte at address array+byte-index (which must lie in dynamic memory).

### `storew`

**VAR:225 1** `storew array word-index value`

`array→word-index = value`, i.e. stores the given value in the word at address array+2\*word-index (which must lie in dynamic memory).

### `sub`

**2OP:21 15** `sub a b → (result)`

Signed 16-bit subtraction.

### `test`

**2OP:7 7** `test bitmap flags ?(label)`

Jump if all of the flags in bitmap are set (i.e. if bitmap & flags == flags).

### `test_attr`

**2OP:10 A** `test_attr object attribute ?(label)`

Jump if object has attribute.

### `throw`

**2OP:28 1C** `5/6 throw value stack-frame`

Opposite of catch: resets the routine call state to the state it had when the given stack frame value was 'caught', and then returns.

### `tokenise`

**VAR:251 1B** `5 tokenise text parse dictionary flag`

This performs lexical analysis (see `read` above).

The dictionary and flag operands are optional.

### `verify`

**0OP:189 D** `3 verify ?(label)`

Verification counts a (two byte, unsigned) checksum of the file from `$0040` onwards and compares this against the value in the game header, branching if the two values agree.

### `window_size`

**EXT:17 11** `6 window_size window y x`

Change size of window in pixels. (Does not change the current display.)

### `window_style`

**EXT:18 12** `6 window_style window flags operation`

Changes attributes for a given window. A bitmap of attributes is given, in which the bits are: 0—keep text within margins, 1—scroll when at bottom, 2—copy text to output stream 2 (the printer), 3—buffer text to word-wrap it between the margins of the window.

The operation, by default, is 0, meaning "set to these settings". 1 means "set the bits supplied". 2 means "clear the ones supplied", and 3 means "reverse the bits supplied".

---

_Source: [Z-Machine Standards Document - Section 15: Dictionary of opcodes](https://zspec.jaredreisinger.com/15-opcodes)_
