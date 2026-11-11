package main

import (
	"fmt"
	"log"
	"strconv"
)

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
	if regDest == "" {
		return
	}
	destVal := regGet(regDest)
	srcVal := value
	if regSrc != "" {
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
	zf := "0"
	if r == 0 {
		zf = "1"
	}
	sf := "0"
	if r>>15 != 0 {
		sf = "1"
	}
	log.Printf("Flags after %s: ZF=%s, SF=%s", spec, zf, sf)
}
