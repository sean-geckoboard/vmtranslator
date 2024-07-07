package src

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CommandType int

const (
	CArithmetic CommandType = 0
	CPush       CommandType = 1
	CPop        CommandType = 2
	CLabel      CommandType = 3
	CGoto       CommandType = 4
	CIf         CommandType = 5
	CFunction   CommandType = 6
	CReturn     CommandType = 7
	CCall       CommandType = 8
)

type command struct {
	commandType CommandType
	arg1        string
	arg2        int
}

type Parser struct {
	file         *os.File
	scanner      *bufio.Scanner
	hasMorelines bool
	command      *command
}

func NewParser(filename string) (*Parser, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	return &Parser{
		file:         file,
		scanner:      scanner,
		hasMorelines: true,
		command:      nil,
	}, nil
}

func (p *Parser) Close() error {
	return p.file.Close()
}

func (p *Parser) HasMoreLines() bool {
	return p.hasMorelines
}

func (p *Parser) Advance() {
	hasMorelines := p.hasMorelines
	hasValidLine := false
	var nextLine string
	for !hasValidLine && hasMorelines {
		hasMorelines = p.scanner.Scan()
		if !hasMorelines {
			err := p.scanner.Err()
			if err != nil {
				fmt.Print(err)
				panic(err)
			}
		}
		nextLine = p.scanner.Text()
		nextLine = strings.TrimSpace(nextLine)
		if len(nextLine) < 1 {
			continue
		}
		if nextLine[0] == '/' {
			continue
		}
		// valid at this point
		hasValidLine = true
	}
	p.hasMorelines = hasMorelines
	if hasMorelines && hasValidLine {
		c := parse(nextLine)
		p.command = &c
	} else {
		p.command = nil
	}
}

func parse(line string) command {
	tokens := strings.Split(line, " ")
	var c command
	switch tokens[0] {
	case "add", "sub", "eq", "gt", "lt", "and", "or", "neg", "not":
		c.commandType = CArithmetic
		c.arg1 = tokens[0]
	case "push":
		c.commandType = CPush
		c.arg1 = tokens[1]
		intVal, err := strconv.Atoi(tokens[2])
		if err != nil {
			panic(err)
		}
		c.arg2 = intVal
	case "pop":
		c.commandType = CPop
		c.arg1 = tokens[1]
		intVal, err := strconv.Atoi(tokens[2])
		if err != nil {
			panic(err)
		}
		c.arg2 = intVal
	}
	return c
}

func (p *Parser) Command() *command {
	return p.command
}
