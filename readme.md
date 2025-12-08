# GOZM – Go Z-Machine Runtime & Engine

GOZM is a Z-Machine interpreter & engine written in Go with a focus on version 3 story files. It executes real games end-to-end (tested with the bundled Infocom classics and custom Inform 6 stories)
while keeping the codebase approachable for anyone curious about the architecture of the original Infocom virtual machine.

## Current Status

- Runs Z3 story files via a terminal frontend (`impl/terminal`) with blocking input, text output, and core engine loop.
- Opcode coverage includes branching, routine calls, stack and object tree handling, text decoding, abbreviations, and dictionary lookup logic as required by V3.
- Command-line debug levels (`-debug 0|1|2`) expose instruction tracing and state dumps to aid reverse engineering and spec validation.
- Comprehensive guide material documents the journey from bytecode loading to full execution.

![web screenshot](./docs/screens/web.png)

## Try the Live Web Version

A live deployment of the web version is available at **https://gozm.benc.dev/** where you can play Z-machine games directly in your browser without any installation. Simply select a story from the menu and start playing classic text adventures with a retro terminal interface.

## Features

- Z-Machine v3 focused interpreter core with structured call stack and object model support.
- Story loader that validates headers, decodes packed addresses, and hydrates initial game memory from Z3 files.
- Text decoding pipeline (ZSCII, abbreviations, dictionary lookup) shared by interpreter and tests.
- Terminal UI providing synchronous input and display, suitable for playing stories directly in the shell.
- **Save/Load support** – Persistent game state with Quetzal-compatible save files, plus browser localStorage integration for the web version.
- **System commands** – Special `/` prefixed commands for save, load, restart, and quit operations.
- Optional debug tooling: instruction tracing, breakpoint hooks, opcode filtering, and pre-generated memory dumps under `misc/` and `test/`.

## Project Layout

- `internal/zmachine/` – machine runtime: instruction dispatch, call stack, object tree, I/O hooks.
- `internal/decode/` – helpers for unpacking V3 headers, operands, and text (abbreviations, ZSCII tables).
- `impl/terminal/` – CLI runner that wires stdin/stdout to the interpreter and exposes debug flags.
- `impl/web/` – WASM entry point for running Z-machine games in the browser with a retro terminal UI.
- `web/` – HTML, CSS, and JavaScript frontend for the WASM build, including story file selection menu.
- [`docs/spec/`](docs/spec/README.md) & [`docs/guide/`](docs/guide/part-01.md) – annotated spec excerpts and a tutorial series documenting implementation details.

## Getting Started

### Prerequisites

- Go 1.25 or newer.
- (Optional) Inform 6 compiler and UnZ tooling if you want to build the sample `.inf` sources yourself.

### Build the CLI Runner

```bash
go build -o bin/gozm ./impl/terminal
```

### Run a Story

```bash
./bin/gozm -file stories/minizork.z3
```

Add `-debug 1` for single-step logging or `-debug 2` for instruction traces. The repository ships with several Infocom-compatible story files under `stories/` and compiler fixtures under `test/` for quick smoke testing.

You can also execute directly with `go run ./impl/terminal -file test/core.z3` during development.

#### System Commands

While playing, you can use system commands prefixed with `/` to control the interpreter:

- `/save` – Save the current game state to a file
- `/load` – Load a previously saved game state
- `/restart` – Restart the current story from the beginning
- `/quit` – Exit the interpreter

Note: Game save files are stored in the current working directory by default.

### Build and Run the Web Version

The web version compiles the Go interpreter to WebAssembly and runs Z-machine games directly in your browser with a retro terminal interface.

```bash
make web
```

This creates `web/main.wasm` which is loaded by `web/index.html`. To serve the web app locally:

```bash
make serve
```

This will launch a development server (using Vite via `npx`) and open the app in your browser. The interface includes:

- A menu bar for selecting from bundled story files (Mini Zork, Moonglow, Catseye, Adventure, etc.).
- A terminal-style display with monospace font rendering game output.
- An input field for entering commands, with input echoed back to the display.
- Story files are fetched on demand from the `web/stories/` directory.
- System commands (`/save`, `/load`, `/restart`, `/quit`) work in the browser just like the terminal version.

The WASM build uses Go's `syscall/js` package to bridge between the Z-machine core and browser JavaScript:

- `impl/web/main.go` – entry point that loads story files via HTTP and initializes the machine.
- `impl/web/webext.go` – implements the `External` interface, routing text output and input through JavaScript callbacks.
- `web/gozm.js` – JavaScript glue code that handles DOM manipulation, input events, and WASM initialization.

### Make Targets

`make build` compiles the terminal binary, `make run STORY=minizork` launches a bundled game, `make story STORY=scratch` rebuilds the matching Inform 6 source and dumps the generated bytecode, `make web` builds the WASM binary, and `make serve` starts a local web server for testing the browser frontend.

## Design Notes

- Instruction Dispatch: `internal/zmachine/step.go` implements the fetch-decode-execute loop keyed by opcode tables defined for V3. Branching helpers live alongside operand decoding for clarity.
- Memory Model: `internal/zmachine/machine.go` loads the packed story into RAM, maps dynamic/static ranges, and surfaces helper methods for address translation.
- Text Handling: `internal/decode/decode.go` maps ZSCII to UTF-8 and surfaces abbreviation expansion used both by the interpreter and tooling.
- IO Abstraction: `internal/zmachine/external.go` defines interfaces so alternative frontends (WASM or scripting) can supply custom input/output streams without altering the core.

For a narrative walkthrough of these pieces, start with the tutorial series:

- [Part 1 – The core](docs/guide/part-01.md)
- [Part 2 – Basic program execution](docs/guide/part-02.md)

More chapters are planned as work continues on streams, save/restore, and alternate frontends.

## Roadmap

- [x] Branching, arithmetic, stack, call/return, and object tree opcodes.
- [x] Abbreviation tables and ZSCII decoding.
- [x] Terminal input pipeline.
- [x] WASM frontend with browser IO and story file selection.
- [x] Sound playback stubs.
- [ ] Output streams (screen/logging) per spec.
- [x] SAVE/RESTORE/RESTART handling.
- [x] Persistent state in browser (localStorage/IndexedDB).
- [ ] Regression test suite driven by official specification transcripts.
- [x] Browser prefs and theme support.
- [ ] Acknowledgements & about and credits screen.

## Tools & References

- UnZ – https://github.com/heasm66/UnZ
- Inform 6 compiler – `sudo apt install inform6-compiler`
- Z-Machine specs – https://zspec.jaredreisinger.com/ & https://inform-fiction.org/zmachine/standards/z1point1/
- Background reading – https://intfiction.org/t/process-of-writing-a-z-machine-interpreter/53231/5 and notes under `misc/`
- ZILF toolkit – https://zilf.io/
