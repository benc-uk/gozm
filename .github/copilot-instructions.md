You are implementing a Go package for a Z-machine which is a virtual machine for text adventure games.
The package should include functionality for decoding Z-machine strings and handling Z-machine data representations.

The implementation should focus on Z-machine version 3 only

The code is broken as follows:

- internal/zmachine/\*.go - for the Z-machine implementation
- internal/decode/\*.go - for decoding Z-machine data formats and strings
- impl/terminal\*.go - for the terminal/CLI input/output implementation and executable app
- impl/web/\*.go - for the web/WASM input/output implementation, with a WASM entry point
  and two implementations one for CLI/terminal and another for web/WASM in the

# Guidance

- When updating or adding Go code in the internal/zmachine or internal/decode packages, read the specification in docs/specs/\*.md to understand the Z-machine architecture and data formats.
- When working on Story files these have a .inf extension and are written in Inform 6 language. The language manual is externally hosted at https://www.inform-fiction.org/manual/html/contents.html
- If you need to execute or test the Z-machine implementation, use the CLI terminal implementation in impl/terminal which provides a command line interface to load and run Z-machine story files
  - It is an interactive terminal application so if you run it, it needs to be wrapped in a timeout
  - The story files is accepted as a command line argument -file <storyfile.inf>
  - There are three debug levels: 0 (none), 1 (some), 2 (trace). Set with -debug <0|1|2>
