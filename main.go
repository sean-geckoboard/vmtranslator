package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cloncode.com/vmtranslator/src"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("missing arguments, expected: ./vmtranslator FileIn.vm")
		return
	}
	inFileName := os.Args[1]
	if inFileName == "" {
		fmt.Println("failed to parse inFileName ")
		return
	}

	// output file
	outFileName := fmt.Sprintf("%s.asm", getJustFileName(inFileName))
	fmt.Printf("output file: %s\n", outFileName)
	cw := src.NewCodeWriter(outFileName)
	defer cw.Close()

	// input files
	var inFiles []string
	if isDir(inFileName) {
		fmt.Printf("translating directory: %s\n", inFileName)
		inFiles = getFilesInDir(inFileName)
		for _, v := range inFiles {
			fmt.Printf("translating file: %s\n", v)
		}
	} else {
		fmt.Printf("translating file: %s\n", inFileName)
		inFiles = []string{inFileName}
	}

	// translate
	for _, inFile := range inFiles {
		err := src.Translate(inFile, cw)
		if err != nil {
			fmt.Printf("err: %s", err)
			return
		}
	}

	fmt.Println("done, success")
}

func isDir(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("err reading file: %w", err)
	}

	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("err stat file: %w", err)
	}

	return fileInfo.IsDir()
}

func getFilesInDir(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("failed to read dir: %w", err)
		return nil
	}

	fileNames := []string{}
	for _, e := range entries {
		fileNames = append(fileNames, fmt.Sprintf("%s/%s", dir, e.Name()))
	}
	return fileNames
}

func getJustFileName(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
