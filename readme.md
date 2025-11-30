# Z Machine Implementation in Go

This repository contains an implementation of the Z Machine interpreter and runtime, a virtual machine designed for running text adventure games, in the Go programming language. The Z Machine was originally developed by Infocom for their interactive fiction games.

This is EXTREMELY early work-in-progress. The goal is to create a fully functional Z Machine interpreter that can run classic Infocom games as well as other Z Machine-compatible titles.

The target is to support Z Machine version 3 only at this time.

## Tutorial

I found getting started with even a minimal Z Machine implementation quite challenging due to the complexity and a badly written specification that assumes a LOT of prior knowledge. Getting a simple "Hello World" working is like facing a mountain of unknowns.

I found so little information on how to get started I thought it would be helpful to document the start of my development journey step-by-step for others who might want to embark on a similar path. I've split this into a separate parts

- [Part 1 - The core](docs/guide/part-01.md)
- Part 2 - On a different branch

These walk through the process of building a basic Z Machine interpreter in Go, explaining all the gotcha and head-scratchers I encountered.

## Tools

- UnZ https://github.com/heasm66/UnZ
- Inform6 Compiler `sudo apt install inform6-compiler` (Debian/Ubuntu)

## Sources & References

- Spec (Nice) https://zspec.jaredreisinger.com/
- Spec (official) https://inform-fiction.org/zmachine/standards/z1point1/
- Forum post https://intfiction.org/t/process-of-writing-a-z-machine-interpreter/53231/5

Tools

- ZILF https://zilf.io/
