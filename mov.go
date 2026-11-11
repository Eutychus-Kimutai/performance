package main

import (
	"fmt"
	"log"
	"os"
)

func movImmediateReg(opcode byte, str []byte, i int, outputFile *os.File, mod byte) int {
	content := ""
	regmap := []string{}
	if opcode >= 0xb8 {
		regmap = regW1
	} else {
		regmap = regW0
	}
	reg := opcode & 0x07 // Extract the register from the opcode
	immLength := immediateSize(opcode)

	var value int16
	if immLength == 2 {
		value = int16(str[i+1]) | int16(str[i+2])<<8
	} else {
		value = int16(str[i+1])
	}
	log.Printf("MOV %s, %d", regmap[reg], value)
	content = fmt.Sprintf("MOV %s, %d\n", regmap[reg], value)
	regEntries(regmap[reg], "", fmt.Sprintf("%d", value), "MOV")
	writeToFile(outputFile, content)
	return 1 + immLength
}
