package main

import (
	"fmt"
	"log"
	"os"
)

type instructionSpec struct {
	name      string
	dirAccess bool
}

func write(outputFile *os.File, content string, args ...any) {
	line := fmt.Sprintf(content, args...)
	log.Print(line)
	outputFile.WriteString(content + "\n")
}

func regMapFor(w byte) []string {
	if w == 0 {
		return regW0
	}
	return regW1
}

func endByteGuard(i int, str []byte) bool {
	if i+1 > len(str) {
		log.Print("Reached end of bytes, no more bytes to read for the instruction")
		return false
	}
	return true
}

func readImm(str []byte, off, size int) int16 {
	if size == 2 {
		return int16(str[off]) | int16(str[off+1])<<8
	} else {
		return int16(str[off])
	}
}

func regMemInstruction(spec instructionSpec, opcode byte, str []byte, i int, outputFile *os.File, reg, rm, d, mod byte) int {
	n := 1 + modrmsize(mod, rm)
	regmap := regMapFor(opcode & 0x01)

	switch mod {
	case 0x03:
		// Register to Register
		if d == 0 {
			write(outputFile, "%s %s, %s", spec.name, regmap[rm], regmap[reg])
			regEntries(regmap[rm], regmap[reg], "", spec.name)
		} else {
			write(outputFile, "%s %s, %s", spec.name, regmap[reg], regmap[rm])
			regEntries(regmap[reg], regmap[rm], "", spec.name)
		}
	case 0x00:
		if spec.dirAccess && rm == 0x06 {
			// Direct memory access
			addr := int16(str[i+2]) | int16(str[i+3])<<8
			if d == 0 {
				write(outputFile, "%s [%d], %s", spec.name, addr, regmap[reg])
				regEntries("", fmt.Sprintf("[%d]", addr), regmap[reg], spec.name)
			} else {
				write(outputFile, "%s %s, [%d]", spec.name, regmap[reg], addr)
				regEntries(regmap[reg], fmt.Sprintf("[%d]", addr), "", spec.name)
			}
			n += 2 // Additional bytes for direct address
		}
		// Memory access through effective address
		if d == 0 {
			write(outputFile, "%s [%s], %s", spec.name, regEA[rm], regmap[reg])
			regEntries("", regEA[rm], regmap[reg], spec.name)
		} else {
			write(outputFile, "%s %s, [%s]", spec.name, regmap[reg], regEA[rm])
			regEntries(regmap[reg], regEA[rm], "", spec.name)
		}
	case 0x01:
		// Memory access with 8-bit displacement
		disp := int16(int8(str[i+2]))
		var regEAmap = regEAMod01
		if d == 0 {
			write(outputFile, "%s [%s + %d], %s", spec.name, regEAmap[rm], disp, regmap[reg])
			regEntries("", fmt.Sprintf("[%s + %d]", regEAmap[rm], disp), regmap[reg], spec.name)
		} else {
			write(outputFile, "%s %s, [%s + %d]", spec.name, regmap[reg], regEAmap[rm], disp)
			regEntries(regmap[reg], fmt.Sprintf("[%s + %d]", regEAmap[rm], disp), "", spec.name)
		}
	case 0x02:
		// Memory access with 16-bit displacement
		disp := int16(str[i+2]) | int16(str[i+3])<<8
		var regEAmap = regEAMod10
		if d == 0 {
			write(outputFile, "%s [%s + %d], %s", spec.name, regEAmap[rm], disp, regmap[reg])
			regEntries("", fmt.Sprintf("[%s + %d]", regEAmap[rm], disp), regmap[reg], spec.name)
		} else {
			write(outputFile, "%s %s, [%s + %d]", spec.name, regmap[reg], regEAmap[rm], disp)
			regEntries(regmap[reg], fmt.Sprintf("[%s + %d]", regEAmap[rm], disp), "", spec.name)
		}
	default:
		log.Printf("Unsupported mod value: %d", mod)
	}
	return n
}

func immRegMemInstruction(spec instructionSpec, opcode byte, str []byte, i int, outputFile *os.File, reg, rm, d, mod byte, regmap []string) int {
	immSize := immediateSize(opcode)
	immOffset := i + 1 + modrmsize(mod, rm)
	value := readImm(str, immOffset, immSize)

	switch mod {
	case 0x03:
		// Register to Register
		if d == 0 {
			write(outputFile, "%s %s, %d", spec.name, regmap[rm], value)
			regEntries(regmap[rm], "", fmt.Sprintf("%d", value), spec.name)
		} else {
			write(outputFile, "%s %s, %d", spec.name, regmap[reg], value)
			regEntries(regmap[reg], "", fmt.Sprintf("%d", value), spec.name)
		}
	case 0x00:
		if spec.dirAccess && rm == 0x06 {
			// Direct memory access
			addr := int16(str[i+2]) | int16(str[i+3])<<8
			if d == 0 {
				write(outputFile, "%s [%d], %d", spec.name, addr, value)
				regEntries("", fmt.Sprintf("[%d]", addr), fmt.Sprintf("%d", value), spec.name)
			} else {
				write(outputFile, "%s %d, [%d]", spec.name, value, addr)
				regEntries(fmt.Sprintf("%d", value), fmt.Sprintf("[%d]", addr), "", spec.name)
			}
			return 1 + modrmsize(mod, rm) + immSize + 2 // Additional bytes for direct address
		}
		if d == 0 {
			write(outputFile, "%s [%s], %d", spec.name, regEA[rm], value)
			regEntries("", regEA[rm], fmt.Sprintf("%d", value), spec.name)
		} else {
			write(outputFile, "%s %d, [%s]", spec.name, value, regEA[rm])
			regEntries(fmt.Sprintf("%d", value), regEA[rm], "", spec.name)
		}
	case 0x01:
		// Memory access with 8-bit displacement
		disp := int16(int8(str[i+2]))
		var regEAmap = regEAMod01
		if d == 0 {
			write(outputFile, "%s [%s + %d], %d", spec.name, regEAmap[rm], disp, value)
			regEntries("", fmt.Sprintf("[%s + %d]", regEAmap[rm], disp), fmt.Sprintf("%d", value), spec.name)
		} else {
			write(outputFile, "%s %d, [%s + %d]", spec.name, value, regEAmap[rm], disp)
			regEntries(fmt.Sprintf("%d", value), fmt.Sprintf("[%s + %d]", regEAmap[rm], disp), "", spec.name)
		}

	case 0x02:
		// Memory access with 16-bit displacement
		disp := int16(str[i+2]) | int16(str[i+3])<<8
		var regEAmap = regEAMod10
		if d == 0 {
			write(outputFile, "%s [%s + %d], %d", spec.name, regEAmap[rm], disp, value)
			regEntries("", fmt.Sprintf("[%s + %d]", regEAmap[rm], disp), fmt.Sprintf("%d", value), spec.name)
		} else {
			write(outputFile, "%s %d, [%s + %d]", spec.name, value, regEAmap[rm], disp)
			regEntries(fmt.Sprintf("%d", value), fmt.Sprintf("[%s + %d]", regEAmap[rm], disp), "", spec.name)
		}
	default:
		log.Printf("Unsupported mod value: %d", mod)
	}
	return 1 + modrmsize(mod, rm) + immSize
}

func immRegInstruction(spec instructionSpec, opcode byte, str []byte, i int, outputFile *os.File, reg, rm, d, mod byte, regmap []string) int {
	var value int16
	immLength := immediateSize(opcode)

	if immLength == 2 {
		value = int16(str[i+1]) | int16(str[i+2])<<8
	} else {
		value = int16(int8(str[i+1]))
	}
	if d == 0 {
		write(outputFile, "%s %s, %d", spec.name, regmap[rm], value)
		regEntries(regmap[rm], "", fmt.Sprintf("%d", value), spec.name)
	} else {
		write(outputFile, "%s %s, %d", spec.name, regmap[reg], value)
		regEntries(regmap[reg], "", fmt.Sprintf("%d", value), spec.name)
	}
	return 1 + immLength
}
