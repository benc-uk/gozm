# Appendix B. Conventional contents of the header

The header table in [section 11](https://zspec.jaredreisinger.com/11-header) details everything the interpreter needs to know about the header's contents. A few other slots in the header are meaningful but only by convention (since existing Infocom and Inform games write them). These additional slots are described here.

As in [S11](https://zspec.jaredreisinger.com/11-header), Hex means the address, in hexadecimal; V the earliest version in which the feature appeared; Dyn means that the byte or bit may legally be changed by the game during play (no others may legally be changed by the game); Int means that the interpreter may (in some cases must) do so.

## Conventional Header Fields

| Hex | V   | Dyn | Int | Contents                                                                                                                                                                                        |
| --- | --- | --- | --- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | 1   |     |     | **Flags 1:**                                                                                                                                                                                    |
|     | 3   |     | \*  | Bit 3: The legendary "Tandy" bit (see note)                                                                                                                                                     |
| 2   | 1   |     |     | Release number (word)                                                                                                                                                                           |
| 10  | 1   | \*  |     | **Flags 2:**                                                                                                                                                                                    |
|     | 3   |     |     | Bit 4: Set in the Amiga version of The Lurking Horror so presumably to do with sound effects?                                                                                                   |
|     | ?   | ?   | \*  | Bit 10: Possibly set by interpreter to indicate an error with the printer during transcription                                                                                                  |
| 12  | 2   |     |     | Serial code (six characters of ASCII)                                                                                                                                                           |
|     | 3   |     |     | Serial number (ASCII for the compilation date in the form YYMMDD)                                                                                                                               |
| 38  | 6   |     | \*  | Infocom used this as 8 bytes of ASCII: the player's user-name on their own mainframe, useful to identify which person played a saved-game (however, the bytes are all 0 in shipped story files) |
| 3C  |     |     |     | Inform 6 stores 4 bytes of ASCII here, giving the version of Inform used to compile the file: for instance, 6.11.                                                                               |

## Notes

1. In Versions 1 to 3, bits 0 and 7 of 'Flags 1′ are unused. (The meaning of bit 2 has recently been discovered: see section 11.) In later Versions, bits 0, 6 and 7 are unused. In 'Flags 2′, bits 9 and 11–15 are unused. Infocom used up almost the whole header: only the bytes at `$32` and `$33` are unused in any Version, and those are now allocated for standard interpreters to give their Revision numbers.

2. **The Tandy Bit:** Some early Infocom games were sold by the Tandy Corporation, who seem to have been sensitive souls. Zork I pretends not to have sequels if it finds the Tandy bit set. And to quote Paul David Doherty:

   > In The Witness, the Tandy Flag can be set while playing the game, by typing `$DB` and then `$TA`. If it is set, some of the prose will be less offensive. For example, "private dicks" become "private eyes", "bastards" are only "idiots", and all references to "slanteyes" and "necrophilia" are removed.

   We live in an age of censorship.

3. For comment on interpreter numbers, see [S11](https://zspec.jaredreisinger.com/11-header). Infocom's own interpreters were generally rewritten for each of versions 3 to 6. For instance, interpreters known to have been shipped with the Macintosh gave version letters B, C, G, I (Version 3), E, H, I (Version 4), A, B, C (Version 5) and finally 6.1 for Version 6. (Version 6 interpreters seem to have version numbers rather than letters.) See the Infocom fact sheet for fuller details.

4. **Distinguishing Story Files:** Inform 6 story files are easily distinguished from all other story files by their usage of the last four header bytes. Inform 1 to 5 story files are best distinguished from Infocom ones by the serial code date: anything before 930000 is either an Infocom file, or a fake. Clearly there is no point going to any trouble to prevent fakes, but with a little practice it's easy to tell whether Zilch or Inform compiled a file from the style of code generated.

---

_Source: [Z-Machine Standards Document - Appendix B: Conventional contents of the header](https://zspec.jaredreisinger.com/B-conventional-header)_
