package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type ByteCounter struct {
	count int64
}

func (bc *ByteCounter) Write(p []byte) (int, error) {
	bc.count += int64(len(p))
	return len(p), nil
}

func main() {

	// flag.*() returns a Pointer
	mainFs := flag.NewFlagSet("mainFlagSet", flag.ContinueOnError)
	mainFs.SetOutput(ioutil.Discard)

	bytesCounterPtr := mainFs.Bool("c", false, "The number of bytes in each input file is written to the standard output.")
	linesCounterPtr := mainFs.Bool("l", false, "The number of lines in each input file is written to the standard output.")
	wordsCounterPtr := mainFs.Bool("w", false, "The number of words in each input file is written to the standard output.")
	charactersCounterPtr := mainFs.Bool(
		"m",
		false,
		"The number of characters in each input file is written to the standard output. If the current locale does not support multibyte characters, this is equivalent to the -c option. This will cancel out any prior usage of the -c option.",
	)
	// Parse the command-line arguments with the custom FlagSet
	err := mainFs.Parse(os.Args[1:])

	if err != nil {
		errString := err.Error()
		err_splitter := strings.Split(errString, " ")
		invalid_option_index := len(err_splitter) - 1
		invalid_option := err_splitter[invalid_option_index][1:]
		fmt.Printf("lyswc: illegal option -- %s", invalid_option)
		return
	}

	hasFlag := mainFs.NFlag() > 0

	if !hasFlag {
		*bytesCounterPtr = true
		*linesCounterPtr = true
		*wordsCounterPtr = true
	}

	var reader io.Reader

	var filePath string

	var buffer []byte

	args := mainFs.Args()

	inputInfo, _ := os.Stdin.Stat()
	isInputFromStdin := inputInfo.Mode()&os.ModeCharDevice == 0
	if isInputFromStdin {
		buffer, _ = ioutil.ReadAll(os.Stdin)
		reader = strings.NewReader(string(buffer))
	} else {
		if len(args) < 1 {
			fmt.Println("Usage: lyswc <filepath>")
			return
		}
		filePath = args[0]
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("lyswc: %s", err.Error())
			return
		}
		defer file.Close()
		reader = file
	}

	result := ""
	// Create a new Scanner
	scanner := bufio.NewScanner(reader)
	if *linesCounterPtr || *wordsCounterPtr {
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

	// Reset Reader and read from the beginning
	// reader.(*os.File) is actually assert that "reader" is type of "os.File"
	if isInputFromStdin {
		reader = strings.NewReader(string(buffer))
	} else {
		file, _ := reader.(*os.File)
		file.Seek(0, io.SeekStart)
	}

	if *charactersCounterPtr {
		// Mimicking wc -m
		anotherScanner := bufio.NewScanner(reader)
		anotherScanner.Split(bufio.ScanRunes)
		runesCounter := 0
		for anotherScanner.Scan() {
			runesCounter++
		}
		result += fmt.Sprintf("%8d ", runesCounter)
	} else if *bytesCounterPtr {
		// Mimicking wc -c
		var bc ByteCounter
		byteCounter := &bc
		_, err := io.Copy(byteCounter, reader)
		if err != nil {
			fmt.Printf("lyswc: byteCounter error: %s", err.Error())
		}
		result += fmt.Sprintf("%8d ", byteCounter.count)
	}

	if isInputFromStdin {
		fmt.Println(result)
	} else {
		fmt.Println(result, filePath)
	}
}
