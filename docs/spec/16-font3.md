# 16. Font 3 and character graphics

## 16.1

The following table of 8 × 8 bitmaps gives a suitable appearance for font 3. The font must have a fixed pitch and characters must be printed immediately next to each other in all four directions.

The table includes ASCII characters 32-126 with their 8x8 bitmap representations. Each character is defined by an 8-byte pattern where each byte represents one row of pixels (bits 7-0 from left to right).

**Note:** The complete bitmap table includes all printable ASCII characters with their pixel patterns. Due to the length and complexity of reproducing the entire bitmap table here, the specification provides detailed 8x8 pixel patterns for characters 32 (space) through 126 (~).

Key characteristics:

- Fixed pitch font
- 8x8 pixel characters
- Characters print immediately adjacent in all four directions
- Includes standard ASCII printable characters
- Special graphics characters for drawing boxes and lines

## Remarks

Two different versions of font 3 were supplied by Infocom, which we shall call the Amiga and PC forms (the Atari form is the same as for the PC). The arrow shape differed slightly and so did the rune alphabet. Each was an attempt to map the late Anglian ("futhorc") runic alphabet, which has 33 characters, onto our Latin alphabet.

Most of the mappings are straightforward (e.g., Latin A maps to Anglian a), except that:

- Latin C is mapped to Anglian eo
- K to "other k" (previously a z sound)
- Q to Anglian k (the same rune as c)
- V to ea
- X to z
- Z to oe

The PC runes differ as follows:

- G has an ornamental circle making it more look like "other z"
- K maps to Anglian k (or c)
- Q is an Anglian ea (which resembles the late Anglian q)
- V is an oe
- X is an "other k"
- Z is a symbol Infocom seem to have invented themselves

(Though less well drawn the PC runes arguably have a better sound-mapping.)

### Beyond Zork Font Behavior

The font behaviour of Beyond Zork, which does have bit 3 of 'Flags 2′ set, is rather complicated and depends on the interpreter number it finds in the header (see [S11](https://zspec.jaredreisinger.com/11-header)). Specifically:

1. **(Digital terminal)** BZ asks whether the player has a VT220 terminal (a model capable of character graphics) and uses font 3 if and only if the answer is yes. (An in-house convenience: Infocom used a Digital mainframe.)

2. **(Apple IIe)** BZ never uses font 3.

3. **(Macintosh)** BZ always uses font 3.

4. **(Amiga)** BZ always uses font 3.

5. **(Atari ST)** BZ always uses font 3.

6. **(MSDOS)** BZ uses font 3 if it finds bit 3 of 'Flags 2′ set (indicating that a graphical screen mode is in use) and otherwise uses IBM PC graphics codes. These need to be converted back into ASCII. The conversion process used by the Zip interpreter is as follows:

   | Code                               | Converts to                   |
   | ---------------------------------- | ----------------------------- |
   | 179                                | a vertical stroke (ASCII 124) |
   | 186                                | a hash (ASCII 35)             |
   | 196                                | a minus sign (ASCII 45)       |
   | 205                                | an equals sign (ASCII 61)     |
   | all others in the range 179 to 218 | a plus sign (ASCII 43)        |

7. **(Commodore 128)** BZ always uses font 3.

8. **(Commodore 64)** BZ always uses font 3.

9. **(Apple IIc)** BZ uses Apple character graphics (possibly "Mousetext"), but has problems when the units used are not 1x1.

10. **(Apple IIgs)** BZ always uses font 3.

11. **(Tandy)** BZ crashes on the public interpreters.

A similarly tangled process is used in Journey. It is obviously highly unsatisfactory to have to make the decision in the above way, which is why `set_font` is now required to return 0 indicating non-availability of a font.

### Graphics Mode

Stefan Jokisch suggests that Infocom originally intended the graphics bit as a way to develop Version 5 to allow a graphical version in parallel with the normal text one. For instance, when the Infocom MSDOS interpreter starts up, it looks at the graphics flag and:

- if clear, it sets the font width/height to 1/1 (so that screen units are character positions);
- if set, it enters MCGA, a graphical screen mode and sets the font width/height to 8/8 (so that screen units are pixels).

The `COLOR` command in BZ (typed at the keyboard) also behaves differently depending on the interpreter number, which is legal behaviour and has no impact on the specification.

---

_Source: [Z-Machine Standards Document - Section 16: Font 3 and character graphics](https://zspec.jaredreisinger.com/16-font3)_
