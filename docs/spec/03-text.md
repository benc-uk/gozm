# 3. How text and characters are encoded

![icon03](https://zspec.jaredreisinger.com/_images/icon03.gif)

> This technique is similar to the five-bit Baudot code, which was used by early Teletypes before ASCII was invented.
>
> — Marc S. Blank and S. W. Galley, How to Fit a Large Program Into a Small Machine

## Contents

- [3.1 Text](#31-text)
- [3.2 Alphabets](#32-alphabets)
  - [3.2.1 Current alphabet](#321-current-alphabet)
  - [3.2.2 Versions 1 and 2 shift characters](#322-versions-1-and-2-shift-characters)
  - [3.2.3 Versions 3+ shift characters](#323-versions-3-shift-characters)
  - [3.2.4 Shift sequences](#324-shift-sequences)
- [3.3 Abbreviations](#33-abbreviations)
  - [3.3.1 Abbreviation rules](#331-abbreviation-rules)
- [3.4 ZSCII escape](#34-zscii-escape)
- [3.5 Alphabet table](#35-alphabet-table)
  - [3.5.1 Space character](#351-space-character)
  - [3.5.2 Version 1 newline](#352-version-1-newline)
  - [3.5.3 Versions 2-4 alphabet table](#353-versions-2-4-alphabet-table)
  - [3.5.4 Version 1 alphabet table](#354-version-1-alphabet-table)
  - [3.5.5 Versions 5+ custom alphabet tables](#355-versions-5-custom-alphabet-tables)
- [3.6 Padding and incompleteness](#36-padding-and-incompleteness)
  - [3.6.1 Incomplete constructions](#361-incomplete-constructions)
- [3.7 Dictionary truncation](#37-dictionary-truncation)
  - [3.7.1 Versions 1-2 shift locks](#371-versions-1-2-shift-locks)
- [3.8 Definition of ZSCII and Unicode](#38-definition-of-zscii-and-unicode)
  - [3.8.1 ZSCII code range](#381-zscii-code-range)
  - [3.8.2 Control codes 0-31](#382-control-codes-0-31)
  - [3.8.3 ASCII codes 32-126](#383-ascii-codes-32-126)
  - [3.8.4 Input codes 129-154](#384-input-codes-129-154)
  - [3.8.5 Extra characters 155-251](#385-extra-characters-155-251)
  - [3.8.6 Mouse and menu codes 252-254](#386-mouse-and-menu-codes-252-254)
  - [3.8.7 Code 255 undefined](#387-code-255-undefined)
- [Remarks](#remarks)

## 3.1 Text

Z-machine text is a sequence of ZSCII character codes (ZSCII is a system similar to ASCII: see §3.8 below). These ZSCII values are encoded into memory using a string of Z-characters. The process of converting between Z-characters and ZSCII values is given in §§3.2 to 3.7 below.

## 3.2 Alphabets

Text in memory consists of a sequence of 2-byte words. Each word is divided into three 5-bit 'Z-characters', plus 1 bit left over, arranged as:

```
  --first byte--   --second byte--
  7 6 5 4 3 2 1 0  7 6 5 4 3 2 1 0
 [bit] [1st] [2nd]     [3rd]
```

The bit is set only on the last 2-byte word of the text, and so marks the end.

### 3.2.1 Current alphabet

There are three 'alphabets', A0 (lower case), A1 (upper case) and A2 (punctuation) and during printing one of these is current at any given time. Initially A0 is current. The meaning of a Z-character may depend on which alphabet is current.

### 3.2.2 Versions 1 and 2 shift characters

In Versions 1 and 2, the current alphabet can be any of the three. The Z-characters 2 and 3 are called 'shift' characters and change the alphabet for the next character only. The new alphabet depends on what the current one is:

| Z-char | From A0 → | From A1 → | From A2 → |
| ------ | --------- | --------- | --------- |
| 2      | A1        | A2        | A0        |
| 3      | A2        | A0        | A1        |

Z-characters 4 and 5 permanently change alphabet, according to the same table, and are called 'shift lock' characters.

### 3.2.3 Versions 3+ shift characters

In Versions 3 and later, the current alphabet is always A0 unless changed for 1 character only: Z-characters 4 and 5 are shift characters. Thus 4 means "the next character is in A1" and 5 means "the next is in A2". There are no shift lock characters.

### 3.2.4 Shift sequences

An indefinite sequence of shift or shift lock characters is legal (but prints nothing).

## 3.3 Abbreviations

In Versions 3 and later, Z-characters 1, 2 and 3 represent abbreviations, sometimes also called 'synonyms' (for traditional reasons): the next Z-character indicates which abbreviation string to print. If z is the first Z-character (1, 2 or 3) and x the subsequent one, then the interpreter must look up entry 32(z-1)+x in the abbreviations table and print the string at that word address. In Version 2, Z-character 1 has this effect (but 2 and 3 do not, so there are only 32 abbreviations).

### 3.3.1 Abbreviation rules

Abbreviation string-printing follows all the rules of this section except that an abbreviation string must not itself use abbreviations and must not end with an incomplete multi-Z-character construction (see §3.6.1 below).

## 3.4 ZSCII escape

Z-character 6 from A2 means that the two subsequent Z-characters specify a ten-bit ZSCII character code: the next Z-character gives the top 5 bits and the one after the bottom 5.

## 3.5 Alphabet table

The remaining Z-characters are translated into ZSCII character codes using the "alphabet table".

### 3.5.1 Space character

The Z-character 0 is printed as a space (ZSCII 32).

### 3.5.2 Version 1 newline

In Version 1, Z-character 1 is printed as a new-line (ZSCII 13).

### 3.5.3 Versions 2-4 alphabet table

In Versions 2 to 4, the alphabet table for converting Z-characters into ZSCII character codes is as follows:

```
       06 07 08 09 0a 0b 0c 0d 0e 0f 10 11 12 13 14 15 16 17 18 19 1a 1b 1c 1d 1e 1f
  A0   a  b  c  d  e  f  g  h  i  j  k  l  m  n  o  p  q  r  s  t  u  v  w  x  y  z
  A1   A  B  C  D  E  F  G  H  I  J  K  L  M  N  O  P  Q  R  S  T  U  V  W  X  Y  Z
  A2   ^  0  1  2  3  4  5  6  7  8  9  .  ,  !  ?  _  #  '  "  /  \  -  :  (  )
```

(Character 6 in A2 is printed as a space here, but is not translated using the alphabet table: see §3.4 above. Character 7 in A2, written here as a circumflex `^`, is a new-line.) For example, in alphabet A1 the Z-character 12 (`$0c`) is translated as a capital `G` (ZSCII character code 71).

### 3.5.4 Version 1 alphabet table

Version 1 has a slightly different A2 row in its alphabet table (new-line is not needed, making room for the < character):

```
       06 07 08 09 0a 0b 0c 0d 0e 0f 10 11 12 13 14 15 16 17 18 19 1a 1b 1c 1d 1e 1f
  A2   0  1  2  3  4  5  6  7  8  9  .  ,  !  ?  _  #  '  "  /  \  <  -  :  (  )
```

### 3.5.5 Versions 5+ custom alphabet tables

In Versions 5 and later, the interpreter should look at the word at $34 in the header. If this is zero, then the alphabet table drawn out in §3.5.3 continues in use. Otherwise it is interpreted as the byte address of an alphabet table specific to this story file.

#### 3.5.5.1

Such an alphabet table consists of 78 bytes arranged as 3 blocks of 26 ZSCII values, translating Z-characters 6 to 31 for alphabets A0, A1 and A2. Z-characters 6 and 7 of A2, however, are still translated as escape and newline codes (as above).

## 3.6 Padding and incompleteness

Since the end-bit only comes up once every three Z-characters, a string may have to be 'padded out' with null values. This is conventionally achieved with a sequence of 5's, though a sequence of (for example) 4's would work equally well.

### 3.6.1 Incomplete constructions

It is legal for the string to end while a multi-Z-character construction is incomplete: for instance, after only the top half of an ASCII value has been given. The partial construction is simply ignored. (This can happen in printing dictionary words which have been guillotined to the dictionary resolution of 6 or 9 Z-characters.)

## 3.7 Dictionary truncation

When an interpreter is encrypting typed-in text to match against dictionary words, the following restrictions apply. Text should be converted to lower case (as a result A1 will not be needed unless the game provides its own alphabet table). Abbreviations may not be used. The pad character, if needed, must be 5. The total string length must be 6 Z-characters (in Versions 1 to 3) or 9 (Versions 4 and later): any multi-Z-character constructions should be left incomplete (rather than omitted) if there's no room to finish them. For example, "i" is encrypted as:

```
14, 5, 5, 5, 5, 5, 5, 5, 5
$48a5 $14a5 $94a5
```

### 3.7.1 Versions 1-2 shift locks

In Versions 1 and 2 only, when encoding text for dictionary words, shift-lock Z-characters 4 and 5 are used instead of the single-shift Z-characters 2 and 3 when the next two characters come from the same alphabet.

## 3.8 Definition of ZSCII and Unicode

The character set of the Z-machine is called ZSCII (Zork Standard Code for Information Interchange; pronounced to rhyme with "xyzzy"). ZSCII codes are 10-bit unsigned values between 0 and 1023. Story files may only legally use the values which are defined below. Note that some values are defined only for input and some only for output.

| Code     | Character               | Input/Output |
| -------- | ----------------------- | ------------ |
| 0        | null                    | Output       |
| 1–7      | ―                       |              |
| 8        | delete                  | Input        |
| 9        | tab (V6)                | Output       |
| 10       | ―                       |              |
| 11       | sentence space (V6)     | Output       |
| 12       | ―                       |              |
| 13       | newline                 | Input/Output |
| 14–26    | ―                       |              |
| 27       | escape                  | Input        |
| 28–31    | ―                       |              |
| 32–126   | standard ASCII          | Input/Output |
| 127–128  | ―                       |              |
| 129–132  | cursor u/d/l/r          | Input        |
| 133–144  | function keys f1 to f12 | Input        |
| 145–154  | keypad 0 to 9           | Input        |
| 155–251  | extra characters        | Input/Output |
| 252      | menu click (V6)         | Input        |
| 253      | double-click (V6)       | Input        |
| 254      | single-click            | Input        |
| 255–1023 | ―                       |              |

### 3.8.1 ZSCII code range

The codes 256 to 1023 are undefined, so that for all practical purposes ZSCII is an 8-bit unsigned code.

### 3.8.2 Control codes 0-31

The codes 0 to 31 are undefined except as follows:

#### 3.8.2.1

ZSCII code 0 ("null") is defined for output but has no effect in any output stream. (It is also used as a value meaning "no character" when reporting terminating character codes, but is not formally defined for input.)

#### 3.8.2.2

ZSCII code 8 ("delete") is defined for input only.

#### 3.8.2.3

ZSCII code 9 ("tab") is defined for output in Version 6 only. At the start of a screen line this should print a paragraph indentation suitable for the font being used: if it is printed in the middle of a screen line, it should be converted to a space (Infocom's own interpreters do not do this, however).

#### 3.8.2.4

ZSCII code 11 ("sentence space") is defined for output in Version 6 only. This should be printed as a suitable gap between two sentences (in the same way that typographers normally place larger spaces after the full stops ending sentences than after words or commas).

#### 3.8.2.5

ZSCII code 13 ("carriage return") is defined for input and output.

#### 3.8.2.6

ZSCII code 27 ("escape" or "break") is defined for input only.

### 3.8.3 ASCII codes 32-126

ZSCII codes between 32 ("space") and 126 ("tilde") are defined for input and output, and agree with standard ASCII (as well as all of the ISO 8859 character sets and Unicode). Specifically:

```
    00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f
20     !  "  #  $  %  &  '  (  )  *  +  ,  -  .  /
30  0  1  2  3  4  5  6  7  8  9  :  ;  <  =  >  ?
40  @  A  B  C  D  E  F  G  H  I  J  K  L  M  N  O
50  P  Q  R  S  T  U  V  W  X  Y  Z  [  \  ]  ^  _
60  `  a  b  c  d  e  f  g  h  i  j  k  l  m  n  o
70  p  q  r  s  t  u  v  w  x  y  z  {  |  }  ~
```

Note that code `$23` (35 decimal) is a hash mark, not a pound sign. (Code `$7c` (124 decimal) is a vertical stroke.)

#### 3.8.3.1

ZSCII codes 127 ("delete" in some forms of ASCII) and 128 are undefined.

### 3.8.4 Input codes 129-154

ZSCII codes 129 to 154 are defined for input only:

| Code | Key          |
| ---- | ------------ |
| 129  | cursor up    |
| 130  | cursor down  |
| 131  | cursor left  |
| 132  | cursor right |
| 133  | f1           |
| 134  | f2           |
| …​   |              |
| 144  | f12          |
| 145  | keypad 0     |
| 146  | keypad 1     |
| …​   |              |
| 154  | keypad 9     |

### 3.8.5 Extra characters 155-251

The block of codes between 155 and 251 are the "extra characters" and are used differently by different story files. Some will need accented Latin characters (such as French E-acute), others unusual punctuation (Spanish question mark), others new alphabets (Cyrillic or Hebrew); still others may want dingbat characters, mathematical or musical symbols, and so on.

#### 3.8.5.1

**[1.0]** To define which characters are required, the Unicode (or ISO 10646–1) Basic Multilingual Plane character set is used: characters are specified by unsigned 16-bit codes. These values agree with ISO 8859 Latin-1 in the range 0 to 255, and with ASCII and ZSCII in the range 32 to 126. The Unicode standard leaves a range of values, the Private Use Area, free: however, an Internet group called the ConScript Unicode Registry is organising a standard mapping of invented scripts (such as Klingon, or Tolkien's Elvish) into the Private Use Area, and this should be considered part of the Unicode standard for Z-machine purposes.

The Z-machine does not provide access to non-BMP characters (ie characters outside the range U+0000 to U+FFFF).

#### 3.8.5.2

**[1.0]** The story file chooses its stock of extra characters with a "Unicode translation table" as follows. Under Versions 1 to 4, the "default table" is always used (see below). In Version 5 or later, if Word 3 of the header extension table is present and non-zero then it is interpreted as the byte address of the Unicode translation table. If Word 3 is absent or zero, the default table is used.

##### 3.8.5.2.1

The table consists of one byte giving a number N, followed by N two-byte words.

##### 3.8.5.2.2

This indicates that ZSCII characters 155 to 155+N-1 are defined for both input and output. (It's possible for N to be zero, leaving the whole range 155 to 251 undefined.)

##### 3.8.5.2.3

The words in the table give Unicode character codes for each of the ZSCII characters 155 to 155+N-1 in turn.

#### 3.8.5.3

The default table is as shown in Table 1 (see the default character mapping in the remarks section).

#### 3.8.5.4

The defined extra characters are entirely normal ZSCII characters. They can appear in a story file's alphabet table, in an array created by print stream 3 and so on.

##### 3.8.5.4.1

**[1.0]** The interpreter is required to be able to print representations of every defined Unicode character under `$0100` (i.e. of every defined ISO 8859–1 Latin1 character). If no suitable letter forms are available, textual equivalents may be used (such as "ss" in place of German sharp "s", `ß`).

##### 3.8.5.4.2

Normally, and where sensibly possible, all punctuation and letter characters in ISO 8859–1 Latin1 should be readable from the interpreter's keyboard. (However, some interpreters may want to provide alternative keyboard mappings, or to run in a different ISO 8859 set: Cyrillic, for example.)

##### 3.8.5.4.3

**[1.0]** An interpreter is not required to have suitable letter-forms for printing Unicode characters `$0100` to `$FFFF`. (It may, if it chooses, allow the user to configure certain fonts for certain Unicode ranges; but this is not required.) If a Unicode character must be printed which an interpreter has no letter-form for, a question mark should be printed instead.

##### 3.8.5.4.4

The Z-machine is not required to handle complex Unicode formatting like combining characters, bidirectional formatting and unusual line-wrapping rules.

In Versions other than 6, interpreters may either handle these features, or not, in window 0. In window 1, and all version 6 windows, they should be ignored.

##### 3.8.5.4.5

Unicode characters U+0000 to U+001F and U+007F to U+009F are control codes, and must not be used.

### 3.8.6 Mouse and menu codes 252-254

ZSCII codes 252 to 254 are defined for input only:

| Code | Event              |
| ---- | ------------------ |
| 252  | menu click         |
| 253  | mouse double-click |
| 254  | mouse single-click |

Menu clicks are available only in Version 6. A single click, or the first click of a double-click, is passed in as 254. The second click of a double-click is passed in as 253. In Versions 5 and later it is recommended that an interpreter should only send code 254, whether the mouse is clicked once or twice.

### 3.8.7 Code 255 undefined

ZSCII code 255 is undefined. (This value is needed in the "terminating characters table" as a wildcard, indicating "any Input-only character with code 128 or above". However, it cannot itself be printed or read from the keyboard.)

## Remarks

In practice the text compression factor is not really very good: for instance, 155000 characters of text squashes into 99000 bytes. (Text usually accounts for about 75% of a story file.) Encoding does at least encrypt the text so that casual browsers can't read it. Well-chosen abbreviations will reduce total story file size by 10% or so.

The German translation of Zork I uses an alphabet table to make accented letters (from the standard extra characters set) efficient in dictionary words. In Version 6, Shogun also uses an alphabet table.

Unicode translation tables are new in Standard 1.0: in Standard 0.2, the extra characters were always mapped using the default Unicode translation table.

Note that if a random stretch of memory is accidentally printed as a string (due to an error in the story file), illegal ZSCII codes may well be printed using the 4-Z-character escape sequence. It's helpful for interpreters to filter out any such illegal codes so that the resulting on-screen mess will not cause trouble for the terminal (e.g. by causing the interpreter to print ASCII 12, clear screen, or 7, bell sound).

The continental European quotation marks << and >> should have spacing which looks sensible either in French style <<Merci!>> or in German style >>Danke!<<.

Ideally, an interpreter should be able to read time delays (for timed input) from stream 1 (i.e., from a script file). See the remarks in S7.

The Beyond Zork story file is capable of receiving both mouse-click codes (253 and 254), listing both in its terminating characters table and treating them equally.

The extant Infocom games in Versions 4 and 5 use the control characters 1 to 31 only as follows: they all accept 10 or 13 as equivalent, except that Bureaucracy will only accept 13. Bureaucracy needs either 127 or 8 to be a delete code. No other codes are used.

In the past, not just in the Z-machine world, there has been general confusion over the rendering of ASCII/ZSCII/Latin-1/Unicode characters `$27` and `$60`. For the Z-machine, the traditional interpretations of right-single-quote/apostrophe and left-single-quote are preferred over the modern neutral-single-quote and grave accent—see Table 2A of the Inform Designer's Manual. `$22` is a neutral double-quote.

An alternative rendering is to interpret both `$27` and `$60` as neutral quotes, but interpreting `$60` as a grave accent is to be avoided.

No doubt aware of this confusion, Infocom never used character `$60`, and used `$27` almost exclusively as an apostrophe—hardly any single quotes appear in Infocom games. Modern authors would do well to follow their lead.

In Version 3 and later, many of Infocom's interpreters (and some subsequent interpreters, such as ITF's) treat two consecutive Z-characters 4 or 5 as shift locks, contrary to the Standard. As a result, story files should not use multiple consecutive 4 or 5 codes except for padding at the end of strings and dictionary words. In any case, these shift locks are not used in dictionary words, or any of Infocom's story files.

To handle languages like Arabic or Hebrew, text would have to be output "visually", with manual line breaks (possibly via an in-game formatting engine).

Far eastern languages are generally straightforward, except they usually use no spaces, and line wraps can occur almost anywhere. The easiest to way to handle this would be for the game to turn off buffering. A more sophisticated game might include its own formatting engine. Also, fixed-space output is liable to be problematical with most Far Eastern fonts, which use a mixture of "full width" and "half width" forms—all half-width characters would have to be forced to full width.

---

Source: [Z-Machine Standards Document - Text Encoding](https://zspec.jaredreisinger.com/03-text)
