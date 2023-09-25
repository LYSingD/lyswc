package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {

	// flag.*() returns a Pointer
	mainFs := flag.NewFlagSet("mainFlagSet", flag.ContinueOnError)
	mainFs.SetOutput(ioutil.Discard)

	bytesCounterPtr := mainFs.Bool("c", false, "The number of bytes in each input file is written to the standard output.")
	linesCounterPtr := mainFs.Bool("l", false, "The number of lines in each input file is written to the standard output.")
	wordsCounterPtr := mainFs.Bool("w", false, "The number of words in each input file is written to the standard output.")
	// Parse the command-line arguments with the custom FlagSet
	err := mainFs.Parse(os.Args[1:])
	// fmt.Println("tail:", mainFs.Args())

	if err != nil {
		errString := err.Error()
		err_splitter := strings.Split(errString, " ")
		invalid_option_index := len(err_splitter) - 1
		invalid_option := err_splitter[invalid_option_index][1:]
		fmt.Printf("lyswc: illegal option -- %s", invalid_option)
		return
	}

	args := mainFs.Args()

	if len(args) < 1 {
		fmt.Println("Usage: lyswc <filepath>")
		return
	}

	filePath := args[0]
	result := ""

	if *linesCounterPtr || *wordsCounterPtr {
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("lyswc: %s", err.Error())
			return
		}
		defer file.Close()

		// Create a new Scanner
		scanner := bufio.NewScanner(file)
		lineCounter := 0
		wordCounter := 0
		for scanner.Scan() {
			line := scanner.Text()

			lineCounter++

			words := strings.Fields(line)
			wordCounter += len(words)
		}

		if *linesCounterPtr {
			result += fmt.Sprintf("%8d ", lineCounter)
		}

		if *wordsCounterPtr {
			result += fmt.Sprintf("%8d ", wordCounter)
		}

	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("lyswc: %s", err.Error())
		return
	}

	fileName := fileInfo.Name()
	if *bytesCounterPtr {
		fileSize := fileInfo.Size()
		// printed with a width of 8 characters. If fileSize has fewer than 8 characters, it will be right-aligned and padded with spaces on the left.

		result += fmt.Sprintf("%8d ", fileSize)
	}

	result += fmt.Sprintf("%s \n", fileName)
	fmt.Println(result)
}
