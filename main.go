package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

var regArr = []regObj{}
var registerState = ""
var emit = true
var simulate = false
var ip int

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

	runLoop(str, outputFile, false)
	registerState = ""
	regArr = nil
	flags = Flags{}
	runLoop(str, outputFile, true)
	for _, reg := range regArr {
		val, _ := strconv.Atoi(reg.Val)

		registerState += fmt.Sprintf("%s: 0x%04x  (%s)\n", reg.Reg, uint16(val), reg.Val)
	}
	log.Print(regArr)
	log.Printf("Register state:\n%s", registerState)
	log.Printf("ip: [%d]", ip)
}
