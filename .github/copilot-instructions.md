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
