# Z Machine Implementation in Go

This repository contains an implementation of the Z Machine, a virtual machine designed for running text adventure games, in the Go programming language. The Z Machine was originally developed by Infocom for their interactive fiction games.

This is EXTREMELY early work-in-progress. The goal is to create a fully functional Z Machine interpreter that can run classic Infocom games as well as other Z Machine-compatible titles.

The target is to support Z Machine version 3 only at this time.

# Step by Step Guide - Part 1

I found getting started with a minimal Z Machine implementation quite challenging due to the complexity and a specification that assumes a lot of prior knowledge. Getting a simple "Hello World" or basic arithmetic example running is like facing a wall of complete unknowns.

I found so little docs or examples that I thought it would be helpful to document my journey step-by-step for others who might want to embark on a similar path. Step 1 is to get a minimal Z Machine story file running that does some basic arithmetic operations. The absolute minimum to get started.

To see the code that implements this step, check the branch `bare-bones`.

```php
Global score = 30;

[Main;
  score = score + 7;
];
```

# Appendix

## Tools

- UnZ https://github.com/heasm66/UnZ
- Inform6 Compiler `sudo apt install inform6-compiler` (Debian/Ubuntu)

## Sources & References

- Spec (Nice) https://zspec.jaredreisinger.com/
- Spec (official) https://inform-fiction.org/zmachine/standards/z1point1/
- Forum post https://intfiction.org/t/process-of-writing-a-z-machine-interpreter/53231/5

Tools

- ZILF https://zilf.io/
