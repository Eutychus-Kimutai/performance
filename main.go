package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

var regW1 = []string{"AX", "CX", "DX", "BX", "SP", "BP", "SI", "DI"}

var regW0 = []string{"AL", "CL", "DL", "BL", "AH", "CH", "DH", "BH"}
var regEA = []string{"BX+SI", "BX+DI", "BP+SI", "BP+DI", "SI", "DI", "BP", "BX"}
var regEAMod01 = []string{
	"BX+SI",
	"BX+DI",
	"BP+SI",
	"BP+DI",
	"SI",
	"DI",
	"BP",
	"BX",
}
var regEAMod10 = []string{
	"BX+SI",
	"BX+DI",
	"BP+SI",
	"BP+DI",
	"SI",
	"DI",
	"BP",
	"BX",
}

type regObj struct {
	Reg string
	Val string
}

var regArr = []regObj{}
var registerState = ""

func writeToFile(file *os.File, content string) {
	_, err := file.WriteString(content)
	if err != nil {
		log.Fatalf("Failed to write to output file: %v", err)

	}
}

func modrmsize(mod, rm byte) int {
	switch mod {
	case 0x00:
		if rm == 0x06 {
			return 3 // 16-bit displacement
		}
		return 1 // No displacement
	case 0x01:
		return 2 // 8-bit displacement
	case 0x02:
		return 3 // 16-bit displacement
	default:
		return 1 // Register addressing mode, no displacement

	}
}

func immediateSize(opcode byte) int {
	if (opcode >= 0xb0 && opcode <= 0xb7) || opcode == 0xC6 || opcode == 0x80 || opcode == 0x83 || opcode == 0xc6 {
		return 1 // 8-bit immediate
	} else if opcode == 0xC7 || opcode == 0x81 || (opcode >= 0xb8 && opcode <= 0xbf) {
		return 2 // 16-bit immediate
	} else if opcode == 0x04 || opcode == 0x2c || opcode == 0x3c {
		return 1 // 8-bit immediate for ADD AL, imm8
	} else if opcode == 0x05 || opcode == 0x2d || opcode == 0x3d {
		return 2 // 16-bit immediate for ADD AX, imm16
	}
	return 0 // No immediate
}

func main() {
	source := ""
	if len(os.Args) > 1 {

		source = os.Args[1]
	} else {
		scanner := bufio.NewScanner(os.Stdin)

		if scanner.Scan() {
			source = scanner.Text()
		}
	}

	path := fmt.Sprintf("./%s", source)
	str, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// create a new file to write the output for the disassembled code
	outputFile, err := os.Create("output.asm")
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outputFile.Close()

	for i := 0; i < len(str); {
		opcode := str[i]
		d := (opcode & 0x02) >> 1
		w := opcode & 0x01
		var regmap []string
		if w == 0 {
			regmap = regW0
		} else {
			regmap = regW1
		}
		if i+1 >= len(str) {
			fmt.Printf("Reached end of bytecode, no more bytes to read for ModR/M byte")
			break
		}

		modrm := str[i+1]
		mod := (modrm & 0xC0) >> 6
		reg := (modrm & 0x38) >> 3
		rm := modrm & 0x07

		switch {

		// Handle MOV instructions
		case opcode >= 0x88 && opcode <= 0x8B:
			if !endByteGuard(i, str) {
				i++
				break
			}
			i += regMemInstruction(instructionSpec{name: "MOV", dirAccess: true}, opcode, str, i, outputFile, reg, rm, d, mod)
		case opcode == 0xC6 || opcode == 0xC7:
			if !endByteGuard(i, str) {
				i++
				break
			}
			i += immRegMemInstruction(instructionSpec{name: "MOV", dirAccess: true}, opcode, str, i, outputFile, reg, rm, d, mod, regmap)
			// handle add instructions
		case opcode >= 0xB0 && opcode <= 0xBF:
			i += movImmediateReg(opcode, str, i, outputFile, mod)

		case opcode <= 0x03:
			if !endByteGuard(i, str) {
				i++
				break
			}
			i += regMemInstruction(instructionSpec{name: "ADD", dirAccess: true}, opcode, str, i, outputFile, reg, rm, d, mod)
		case opcode == 0x83 || opcode == 0x81 || opcode == 0x80:
			if !endByteGuard(i, str) {
				i++
				break
			}
			// determine if it's an add or sub instruction based on the reg field in the ModR/M byte
			spec := instructionSpec{}
			switch reg {
			case 0:
				spec.name = "ADD"
			case 5:
				spec.name = "SUB"
			case 7:
				spec.name = "CMP"

			}
			i += immRegMemInstruction(spec, opcode, str, i, outputFile, reg, rm, d, mod, regmap)
		case opcode >= 0x28 && opcode <= 0x2B:
			if !endByteGuard(i, str) {
				i++
				break
			}
			i += regMemInstruction(instructionSpec{name: "SUB", dirAccess: true}, opcode, str, i, outputFile, reg, rm, d, mod)
		case opcode >= 0x2c && opcode <= 0x2d:
			immediateSize := immediateSize(opcode)
			if i+immediateSize >= len(str) {
				fmt.Printf("Reached end of bytecode, no more bytes to read for ModR/M byte")
				i++
				break
			}
			i += immRegInstruction(instructionSpec{name: "SUB", dirAccess: true}, opcode, str, i, outputFile, reg, rm, d, mod, regmap)
		case opcode == 0x04 || opcode == 0x05:
			immediateSize := immediateSize(opcode)
			if i+immediateSize >= len(str) {
				fmt.Printf("Reached end of bytecode, no more bytes to read for ModR/M byte")
				i++
				break
			}

			i += immRegInstruction(instructionSpec{name: "ADD", dirAccess: true}, opcode, str, i, outputFile, reg, rm, d, mod, regmap)
		// Handle CMP instructions
		case opcode >= 0x38 && opcode <= 0x3B:
			if !endByteGuard(i, str) {
				i++
				break
			}
			i += regMemInstruction(instructionSpec{name: "CMP", dirAccess: true}, opcode, str, i, outputFile, reg, rm, d, mod)
		case opcode == 0x3C || opcode == 0x3D:
			immediateSize := immediateSize(opcode)
			if i+immediateSize >= len(str) {
				fmt.Printf("Reached end of bytecode, no more bytes to read for ModR/M byte")
				i++
				break
			}
			i += immRegInstruction(instructionSpec{name: "CMP", dirAccess: true}, opcode, str, i, outputFile, reg, rm, d, mod, regmap)
		case (opcode >= 0x70 && opcode <= 0x7F) || (opcode >= 0xE0 && opcode <= 0xE3):
			if i+1 >= len(str) {
				fmt.Printf("Reached end of bytecode, no more bytes to read for jump displacement")
				i++
				break
			}
			i += jump(opcode, str, i, outputFile)

		default:
			i += 1

		}
	}
	for _, reg := range regArr {
		val, _ := strconv.Atoi(reg.Val)

		registerState += fmt.Sprintf("%s: 0x%04x  (%s)\n", reg.Reg, uint16(val), reg.Val)
	}
	log.Print(regArr)
	log.Printf("Register state:\n%s", registerState)
}
