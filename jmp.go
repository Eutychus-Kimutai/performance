package main

import (
	"fmt"
	"log"
	"os"
)

var jumpNames = map[byte]string{
	0x70: "JO",
	0x71: "JNO",
	0x72: "JB",
	0x73: "JNB",
	0x74: "JE",
	0x75: "JNZ",
	0x76: "JBE",
	0x77: "JA",
	0x78: "JS",
	0x79: "JNS",
	0x7A: "JP",
	0x7B: "JNP",
	0x7C: "JL",
	0x7D: "JNL",
	0x7E: "JLE",
	0x7F: "JG",
	0xE0: "LOOPNZ",
	0xE1: "LOOPZ",
	0xE2: "LOOP",
	0xE3: "JCXZ",
}

func jump(opcode byte, str []byte, i int, outputFile *os.File) int {
	disp := int(int8(str[i+1]))
	target := disp + 2
	name := jumpNames[opcode]

	log.Printf("%s $%+d", name, target)
	content := fmt.Sprintf("%s $%+d\n", name, target)
	writeToFile(outputFile, content)
	return 2
}
