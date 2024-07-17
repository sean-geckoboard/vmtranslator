package src

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type CodeWriter struct {
	outFile         *os.File
	eqCounter       int
	fileName        string
	currentFunction string
}

func NewCodeWriter(outFileName string) *CodeWriter {
	f, err := os.Create(outFileName)
	if err != nil {
		panic(err)
	}

	return &CodeWriter{
		outFile:         f,
		currentFunction: "Global",
	}
}

func (c *CodeWriter) Close() {
	c.outFile.Close()
}

func (c *CodeWriter) setFileName(fileName string) {
	c.fileName = getJustFileName(fileName)
}

func (c *CodeWriter) WriteLabel(label string) {
	lines := []string{
		fmt.Sprintf("(%s.%s$%s)", c.fileName, c.currentFunction, label),
	}
	c.writeLines(lines)
}

func (c *CodeWriter) WriteGoto(label string) {
	lines := []string{
		fmt.Sprintf("@%s.%s$%s", c.fileName, c.currentFunction, label),
		"0;JMP",
	}
	c.writeLines(lines)
}

func (c *CodeWriter) WriteIf(label string) {
	lines := []string{
		"@SP",
		"M=M-1",
		"A=M",
		"D=M",
		fmt.Sprintf("@%s.%s$%s", c.fileName, c.currentFunction, label),
		"D;JNE",
	}
	c.writeLines(lines)
}

func (c *CodeWriter) WriteFunction(functionName string, nVars int) {

}

func (c *CodeWriter) WriteCall(functionName string, nArgs int) {

}

func (c *CodeWriter) WriteReturn() {

}

func (c *CodeWriter) WriteArithmetic(command string) {

	switch command {
	case "add":
		lines := []string{
			"@SP",
			"M=M-1",
			"A=M",
			"D=M",
			"@SP",
			"M=M-1",
			"A=M",
			"M=D+M", // add
			"@SP",
			"M=M+1",
		}
		c.writeLines(lines)
	case "sub":
		lines := []string{
			"@SP",
			"M=M-1",
			"A=M",
			"D=M",
			"@SP",
			"M=M-1",
			"A=M",
			"M=M-D", // sub
			"@SP",
			"M=M+1",
		}
		c.writeLines(lines)
	case "eq":
		lines := []string{
			"@SP",
			"M=M-1",
			"A=M",
			"D=M",
			"@SP",
			"M=M-1",
			"A=M",
			"D=D-M", // 0 if equal
			"M=-1",  // true
			fmt.Sprintf("@eq.end.%d", c.eqCounter),
			"D;JEQ",
			"@SP",
			"A=M",
			"M=0", // false
			fmt.Sprintf("(eq.end.%d)", c.eqCounter),
			"@SP",
			"M=M+1",
		}
		c.writeLines(lines)
		c.eqCounter = c.eqCounter + 1
	case "lt":
		lines := []string{
			"@SP",
			"M=M-1",
			"A=M",
			"D=M",
			"@SP",
			"M=M-1",
			"A=M",
			"D=M-D",
			"M=-1", // true
			fmt.Sprintf("@eq.end.%d", c.eqCounter),
			"D;JLT",
			"@SP",
			"A=M",
			"M=0", // false
			fmt.Sprintf("(eq.end.%d)", c.eqCounter),
			"@SP",
			"M=M+1",
		}
		c.writeLines(lines)
		c.eqCounter = c.eqCounter + 1
	case "gt":
		lines := []string{
			"@SP",
			"M=M-1",
			"A=M",
			"D=M",
			"@SP",
			"M=M-1",
			"A=M",
			"D=M-D",
			"M=-1", // true
			fmt.Sprintf("@eq.end.%d", c.eqCounter),
			"D;JGT",
			"@SP",
			"A=M",
			"M=0", // false
			fmt.Sprintf("(eq.end.%d)", c.eqCounter),
			"@SP",
			"M=M+1",
		}
		c.writeLines(lines)
		c.eqCounter = c.eqCounter + 1
	case "and":
		lines := []string{
			"@SP",
			"M=M-1",
			"A=M",
			"D=M",
			"@SP",
			"M=M-1",
			"A=M",
			"M=D&M", // and
			"@SP",
			"M=M+1",
		}
		c.writeLines(lines)
	case "or":
		lines := []string{
			"@SP",
			"M=M-1",
			"A=M",
			"D=M",
			"@SP",
			"M=M-1",
			"A=M",
			"M=D|M", // and
			"@SP",
			"M=M+1",
		}
		c.writeLines(lines)
	case "neg":
		lines := []string{
			"@SP",
			"M=M-1",
			"A=M",
			"M=-M",
			"@SP",
			"M=M+1",
		}
		c.writeLines(lines)
	case "not":
		lines := []string{
			"@SP",
			"M=M-1",
			"A=M",
			"M=!M",
			"@SP",
			"M=M+1",
		}
		c.writeLines(lines)
	}

}

func (c *CodeWriter) WritePush(segment string, index int) {

	switch segment {
	case "constant":
		lines := []string{
			// index is a value in this case
			fmt.Sprintf("@%d", index),
			"D=A",
			"@SP",
			"A=M",
			"M=D",
			"@SP",
			"M=M+1",
		}
		c.writeLines(lines)
	case "local":
		c.writePushSegment("LCL", index)
	case "argument":
		c.writePushSegment("ARG", index)
	case "this":
		c.writePushSegment("THIS", index)
	case "that":
		c.writePushSegment("THAT", index)
	case "temp":
		c.writePushTemp(index)
	case "pointer":
		c.writePushPointer(index)
	case "static":
		c.writePushStatic(index)
	default:
		fmt.Printf("no matching segment for pop operation: %s", segment)
		return
	}

}

func (c *CodeWriter) writePushSegment(source string, index int) {
	lines := []string{
		fmt.Sprintf("@%d", index),
		"D=A",
		fmt.Sprintf("@%s", source),
		"A=D+M",
		"D=M",
		"@SP",
		"A=M",
		"M=D",
		"@SP",
		"M=M+1",
	}
	c.writeLines(lines)
}

func (c *CodeWriter) writePushTemp(index int) {
	lines := []string{
		fmt.Sprintf("@%d", index),
		"D=A",
		"@R5",
		"A=D+A",
		"D=M",
		"@SP",
		"A=M",
		"M=D",
		"@SP",
		"M=M+1",
	}
	c.writeLines(lines)
}

func (c *CodeWriter) writePushPointer(index int) {
	location := "R3" // this
	if index == 1 {
		location = "R4" // that
	}
	lines := []string{
		fmt.Sprintf("@%s", location),
		"D=M",
		"@SP",
		"A=M",
		"M=D",
		"@SP",
		"M=M+1",
	}
	c.writeLines(lines)
}

func (c *CodeWriter) writePushStatic(index int) {
	lines := []string{
		fmt.Sprintf("@%s.%d", c.fileName, index),
		"D=M",
		"@SP",
		"A=M",
		"M=D",
		"@SP",
		"M=M+1",
	}
	c.writeLines(lines)
}

func (c *CodeWriter) WritePop(segment string, index int) {

	switch segment {
	case "local":
		c.writePopSegment("LCL", index)
	case "argument":
		c.writePopSegment("ARG", index)
	case "this":
		c.writePopSegment("THIS", index)
	case "that":
		c.writePopSegment("THAT", index)
	case "temp":
		c.writePopTemp(index)
	case "pointer":
		c.writePopPointer(index)
	case "static":
		c.writePopStatic(index)
	default:
		fmt.Printf("no matching segment for pop operation: %s", segment)
		return
	}

}

func (c *CodeWriter) writePopSegment(dest string, index int) {
	lines := []string{
		// put index into D
		fmt.Sprintf("@%d", index),
		"D=A",
		// calculate dest address
		fmt.Sprintf("@%s", dest),
		"A=D+M",
		"D=A",
		// store dest address
		"@R13",
		"M=D",
		// pop the value from the stack
		"@SP",
		"M=M-1",
		"A=M",
		"D=M",
		// set memory to the stored address
		"@R13",
		"A=M",
		// set the value at the address
		"M=D",
	}
	c.writeLines(lines)
}

func (c *CodeWriter) writePopTemp(index int) {
	lines := []string{
		// put index into D
		fmt.Sprintf("@%d", index),
		"D=A",
		// calculate dest address
		"@R5",
		"D=D+A",
		// store dest address
		"@R13",
		"M=D",
		// pop the value from the stack
		"@SP",
		"M=M-1",
		"A=M",
		"D=M",
		// set memory to the stored address
		"@R13",
		"A=M",
		// set the value at the address
		"M=D",
	}
	c.writeLines(lines)
}

func (c *CodeWriter) writePopPointer(index int) {
	location := "R3" // this
	if index == 1 {
		location = "R4" // that
	}
	lines := []string{
		// pop the value from the stack
		"@SP",
		"M=M-1",
		"A=M",
		"D=M",
		// set memory to the stored address
		fmt.Sprintf("@%s", location),
		"M=D",
	}
	c.writeLines(lines)
}

func (c *CodeWriter) writePopStatic(index int) {
	lines := []string{
		// pop the value from the stack
		"@SP",
		"M=M-1",
		"A=M",
		"D=M",
		// set the address to the variable
		fmt.Sprintf("@%s.%d", c.fileName, index),
		// set the value at the address
		"M=D",
	}
	c.writeLines(lines)
}

func (c *CodeWriter) writeLines(lines []string) {
	var line string
	for _, l := range lines {
		line = line + l + "\n"
	}
	c.outFile.WriteString(line)
}

func getJustFileName(fileName string) string {
	return strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
}
