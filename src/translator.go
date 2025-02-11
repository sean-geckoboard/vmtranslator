package src

import (
	"fmt"
)

func Translate(inFileName string, cw *CodeWriter) error {
	p, err := NewParser(inFileName)
	if err != nil {
		return fmt.Errorf("new parser: %w", err)
	}
	defer p.Close()

	cw.setFileName(inFileName)

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
			case CLabel:
				cw.WriteLabel(c.arg1)
			case CGoto:
				cw.WriteGoto(c.arg1)
			case CIf:
				cw.WriteIf(c.arg1)
			case CFunction:
				cw.WriteFunction(c.arg1, c.arg2)
			case CCall:
				cw.WriteCall(c.arg1, c.arg2)
			case CReturn:
				cw.WriteReturn()
			default:
				fmt.Printf("cannot handle command of type: %d\n", c.commandType)
			}
		} else {
			panic("error: nil command ??")
		}
	}

	return nil
}
