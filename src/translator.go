package src

import (
	"fmt"
)

func Translate(inFileName string, outFileName string) error {
	p, err := NewParser(inFileName)
	if err != nil {
		return fmt.Errorf("new parser: %w", err)
	}
	defer p.Close()

	cw := NewCodeWriter(outFileName)
	defer cw.Close()

	for p.HasMoreLines() {
		p.Advance()
		if !p.HasMoreLines() {
			return nil
		}

		c := p.Command()
		if c != nil {
			switch c.commandType {
			case CArithmetic:
				cw.WriteArithmetic(c.arg1)
			case CPush:
				cw.WritePush(c.arg1, c.arg2)
			case CPop:
				cw.WritePop(c.arg1, c.arg2)
			default:
				fmt.Printf("cannot handle command of type: %d\n", c.commandType)
			}
		} else {
			fmt.Println("error: nil command ??")
		}
	}

	return nil
}
