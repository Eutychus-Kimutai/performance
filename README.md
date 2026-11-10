# 8086 Disassembler

A small Go program built while working through Casey Muratori's _Performance-Aware Programming_ course. The goal is to learn how machine code maps back to assembly by writing a disassembler from scratch: give it the raw bytes  and it prints the equivalent assembly instructions.

## Usage

```sh
go run . listing_0044_register_movs   # pass a filename as an argument
go run .                              # or type the filename on stdin
```

The disassembled output is written to `output.asm`. Instructions and the final register values are also logged to the console.

## What it understands so far

- `MOV`, `ADD`, `SUB`, `CMP` — register, memory, and immediate forms
- Short jumps — `Jcc`, `LOOP`, `LOOPZ`, `LOOPNZ`, `JCXZ`
- ModR/M addressing, including 8-bit and 16-bit displacements and direct addressing

## Files

- `main.go` — entry point, opcode dispatch loop, register tracking
- `instr.go` — ModR/M decoding and register/memory/immediate handlers
- `mov.go` — `MOV` immediate-to-register
- `jmp.go` — jump instruction table and handler
- `listing_*` — sample bytecode inputs from the course
- `output.asm` — generated disassembly
