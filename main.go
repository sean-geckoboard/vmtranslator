package main

import (
	"fmt"
	"os"
	"strings"

	"cloncode.com/vmtranslator/src"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("missing arguments, expected: ./vmtranslator FileIn.vm")
		return
	}
	inFileName := os.Args[1]
	outFileName := strings.Split(inFileName, ".vm")[0] + ".asm"
	if outFileName == "" {
		fmt.Println("failed to parse input file, expected: FileIn.vm ")
		return
	}

	err := src.Translate(inFileName, outFileName)
	if err != nil {
		fmt.Printf("err: %s", err)
	}

}
