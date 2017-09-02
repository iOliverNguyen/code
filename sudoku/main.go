package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ng-vu/go-stuff/sudoku/ss"
)

var usage = `
  Usage:
    sudoku <input-file>

  Example:
    sudoku inputs/01.txt
`

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println(usage)
		os.Exit(255)
	}

	inputFile := flag.Arg(0)
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Unable to read file:\n  %v\n", err)
		os.Exit(1)
	}

	board, err := ss.Parse(string(data))
	if err != nil {
		fmt.Printf("Input file is not in valid format:\n  %v", err)
		os.Exit(1)
	}

	fmt.Printf("Input Sudoku:\n%v\n", board)

	if err := ss.IsValid(board); err != nil {
		fmt.Printf("The input sudoku is not valid:\n  %v\n", err)
		os.Exit(1)
	}

	result, ok := ss.FindOne(board)
	if !ok {
		fmt.Printf("No result found.\n")
		os.Exit(0)
	}

	fmt.Printf("Found result:\n%v\n", result)
}
