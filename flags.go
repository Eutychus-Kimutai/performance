package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Flags struct {
	ZF string // Zero Flag
	SF string // Sign Flag
}

var flags Flags
var memory = make(map[uint16]uint16)

func evalAddr(s string) uint16 {
	// parse [1000] or [BX+SI] or [BX+4]...etc
	if s[0] != '[' || s[len(s)-1] != ']' {
		log.Printf("Invalid memory address: %s", s)
		return 0
	}
	//  resolving register parts via regGet and displacements via strconv.Atoi.
	regPart := s[1 : len(s)-1]
	if strings.Contains(regPart, "+") {
		parts := strings.Split(regPart, "+")
		addr := 0
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if val := regGet(part); val != "" {

				intVal, _ := strconv.Atoi(val)
				addr += intVal
			} else {
				intVal, _ := strconv.Atoi(part)
				addr += intVal
			}
		}
		return uint16(addr)
	} else {
		if val := regGet(regPart); val != "" {
			intVal, _ := strconv.Atoi(val)
			return uint16(intVal)
		} else {
			intVal, _ := strconv.Atoi(regPart)
			return uint16(intVal)
		}
	}

}

func regSet(dest, val string) {
	for i := range regArr {
		if regArr[i].Reg == dest {
			regArr[i].Val = val
			return
		}
	}
	regArr = append(regArr, regObj{Reg: dest, Val: val})
}
func regGet(reg string) string {
	for i := range regArr {
		if regArr[i].Reg == reg {
			return regArr[i].Val
		}
	}
	return ""
}

func regEntries(regDest, regSrc, value, spec string) {
	if !simulate {
		return
	}

	memSrc := strings.HasPrefix(regSrc, "[") && strings.HasSuffix(regSrc, "]")
	if regDest == "" {
		if memSrc {
			v, err := strconv.Atoi(value)
			if err != nil {
				v, _ = strconv.Atoi(regGet(value))
			}
			memory[evalAddr(regSrc)] = uint16(v)
			return
		}

	}
	destVal := regGet(regDest)

	srcVal := value
	if memSrc {
		srcVal = fmt.Sprintf("%d", memory[evalAddr(regSrc)])
	} else if regSrc != "" {
		srcVal = regGet(regSrc)
	}
	switch spec {
	case "MOV":
		regSet(regDest, srcVal)
	case "ADD":
		destInt, _ := strconv.Atoi(destVal)
		srcInt, _ := strconv.Atoi(srcVal)
		result := destInt + srcInt
		regSet(regDest, fmt.Sprintf("%d", result))
		regEntriesWithFlags(result, spec)
	case "SUB":
		destInt, _ := strconv.Atoi(destVal)
		srcInt, _ := strconv.Atoi(srcVal)
		result := destInt - srcInt
		regSet(regDest, fmt.Sprintf("%d", result))
		regEntriesWithFlags(result, spec)
	case "CMP":
		destInt, _ := strconv.Atoi(destVal)
		srcInt, _ := strconv.Atoi(srcVal)
		result := destInt - srcInt
		regEntriesWithFlags(result, spec)
	default:
		log.Printf("Unsupported instruction: %s", spec)
	}

}
func regEntriesWithFlags(result int, spec string) {
	// After each arithmetic instruction, ZF is set to 1 if the result is zero, and 0 if the result is non-zero
	// SF is set to 1 is the high bit of the result is 1, and 0 if it is zero.
	r := uint16(result)
	flags.ZF = "0"
	flags.SF = "0"
	if r == 0 {
		flags.ZF = "1"
		flags.SF = "0"
	}
	if r>>15 != 0 {
		flags.SF = "1"
		flags.ZF = "0"
	}
	log.Printf("zf: %s, sf: %s", flags.ZF, flags.SF)
}
