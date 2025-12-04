# Z-Machine Standards Document

This directory contains the complete Z-Machine Standards Document, converted to markdown format for easy reference. The Z-machine is a virtual machine designed by Infocom for text adventure games, most famously used for the Zork series.

I claim no copyright over the original content, which was created by Graham Nelson. This markdown version is scraped from the modernized presentation by Jared Reisinger at https://zspec.jaredreisinger.com/.
It was mirrored here for easier access and so AI tools can reference it.
If you would like to have this version removed for any reason, please open an issue or contact the maintainers.

## Introduction

- [Preface](001-preface.md) - Introduction to the Z-machine and standardization
- [Overview of Z-machine architecture](002-overview.md) - High-level architectural overview

## Core Specification

### 1. Memory and Data Structures

- [Section 1: Memory Map](01-memory-map.md) - Organization of Z-machine memory
- [Section 2: Numbers and Arithmetic](02-numbers.md) - Number representation and operations
- [Section 3: Text and Characters](03-text.md) - Text encoding and string compression

### 2. Execution Model

- [Section 4: Instructions](04-instructions.md) - Instruction format and encoding
- [Section 5: Routines and the Call Stack](05-routines.md) - Subroutine calls and returns
- [Section 6: Game State](06-game-state.md) - Save, restore, and undo operations

### 3. Input/Output

- [Section 7: Output Streams and File Handling](07-output.md) - Text output and redirection
- [Section 9: Sound Effects](09-sound.md) - Sound effect system
- [Section 10: Input Streams and Devices](10-input.md) - Keyboard, file, and mouse input

### 4. Data Tables

- [Section 11: Header](11-header.md) - Story file header format
- [Section 12: Objects](12-objects.md) - Object table and properties
- [Section 13: Dictionary and Lexical Analysis](13-dictionary.md) - Word dictionary and parsing

### 5. Instruction Reference

- [Section 14: Opcode Table](14-opcode-table.md) - Complete opcode tables by type
- [Section 15: Dictionary of Opcodes](15-opcodes.md) - Detailed opcode specifications
- [Section 16: Font 3 and Character Graphics](16-font3.md) - Character graphics system

## Appendices

- [Appendix A: Error Messages](A-errors.md) - Error handling guidelines
- [Appendix B: Conventional Contents of the Header](B-conventional-header.md) - Standard header usage
- [Appendix C: Resources Available](C-resources.md) - Interpreters, tools, and compilers
- [Appendix D: A Short History of the Z-machine](D-history.md) - Evolution from Version 1 to 8
- [Appendix E: Statistics](E-statistics.md) - Size and complexity statistics
- [Appendix F: Canonical Story Files](F-canonical-story-files.md) - Known Infocom releases

## Editor's Notes

- [Editor's Note](ZZ01-editors-note.md) - About this modernized version
- [Conventions Used in This Document](ZZ02-conventions.md) - Formatting and notation
- [Opcodes Revisited](ZZ03-opcodes.md) - Alternative opcode organization

## About

This specification describes the Z-machine virtual machine architecture across all versions (1-8), with particular focus on the most common versions (3, 5, and 8). The original specification was created by Graham Nelson through reverse-engineering of Infocom's story files. This markdown conversion is based on Jared Reisinger's modernized presentation.

### Key Version Differences

- **Version 1-2**: Early Apple II and TRS-80 games (limited features)
- **Version 3**: "Standard" series games (most common format, 128K limit)
- **Version 4**: "Plus" series games (larger memory, extended features)
- **Version 5**: "Advanced" series games (256K limit, enhanced I/O)
- **Version 6**: Graphics and mouse support
- **Version 7-8**: Modern Inform-compiled large games

## See Also

- Original canonical version: http://inform-fiction.org/zmachine/standards/z1point1
- Modernized web version: https://zspec.jaredreisinger.com/
- Inform compiler: http://inform-fiction.org/
- IF Archive: https://www.ifarchive.org/

---

_These documents are derived from The Z-Machine Standards Document by Graham Nelson, with modernization by Jared Reisinger._
