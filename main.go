package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type CountResults struct {
	bytesCount      int64
	linesCount      int
	wordsCount      int
	charactersCount int
}

type ByteCounter struct {
	count int64
}

func (bc *ByteCounter) Write(p []byte) (int, error) {
	bc.count += int64(len(p))
	return len(p), nil
}

func main() {
	// flag.*() returns a Pointer
	mainFlagSet := flag.NewFlagSet("mainFlagSet", flag.ContinueOnError)

	// Prevent the default error messages as we're using our own custom error message.
	mainFlagSet.SetOutput(ioutil.Discard)

	bytesCounter, wordsCounter, linesCounter, charactersCounter := parseFlags(mainFlagSet)

	args := mainFlagSet.Args()

	data, filePath, isInputFromStdin := getInputSource(args)

	resultsChan := make(chan CountResults)
	go countAll(data, resultsChan)

	countAllResult := <-resultsChan
	display_output := ""

	if *linesCounter {
		display_output += fmt.Sprintf("%8d", countAllResult.linesCount)
	}

	if *wordsCounter {
		display_output += fmt.Sprintf("%8d", countAllResult.wordsCount)
	}

	if *charactersCounter {
		display_output += fmt.Sprintf("%8d", countAllResult.charactersCount)
	} else if *bytesCounter {
		display_output += fmt.Sprintf("%8d", countAllResult.bytesCount)
	}

	if isInputFromStdin {
		fmt.Println(display_output)
	} else {
		fmt.Println(display_output, filePath)
	}
}

func parseFlags(mainFlagSet *flag.FlagSet) (bytesCounter *bool, linesCounter *bool, wordsCounter *bool, charactersCounter *bool) {

	bytesCounter = mainFlagSet.Bool("c", false, "The number of bytes in each input file is written to the standard output.")
	linesCounter = mainFlagSet.Bool("l", false, "The number of lines in each input file is written to the standard output.")
	wordsCounter = mainFlagSet.Bool("w", false, "The number of words in each input file is written to the standard output.")
	charactersCounter = mainFlagSet.Bool(
		"m",
		false,
		"The number of characters in each input file is written to the standard output. If the current locale does not support multibyte characters, this is equivalent to the -c option. This will cancel out any prior usage of the -c option.",
	)
	// Parse the command-line arguments with the custom FlagSet
	err := mainFlagSet.Parse(os.Args[1:])

	if err != nil {
		errString := err.Error()
		flagName := strings.TrimPrefix(errString, "flag provided but not defined: -")
		fmt.Printf("lyswc: illegal option -- %s", flagName)
		os.Exit(1)
	}

	hasFlag := mainFlagSet.NFlag() > 0

	if !hasFlag {
		*bytesCounter = true
		*linesCounter = true
		*wordsCounter = true
	}

	return bytesCounter, linesCounter, wordsCounter, charactersCounter
}

func getInputSource(args []string) (data []byte, filePath string, isInputFromStdin bool) {
	inputInfo, err := os.Stdin.Stat()
	if err != nil {
		fmt.Printf("lyswc: error: %s", err.Error())
		os.Exit(1)
	}

	isInputFromStdin = inputInfo.Mode()&os.ModeCharDevice == 0
	if isInputFromStdin {
		stdinContents, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Printf("lyswc: error: %s", err.Error())
			os.Exit(1)
		}
		data = stdinContents
	} else {
		if len(args) < 1 {
			fmt.Println("Usage: lyswc <filepath>")
			os.Exit(1)
		}
		filePath = args[0]
		data, err = ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Printf("lyswc: %s", err.Error())
			os.Exit(1)
		}
	}

	return data, filePath, isInputFromStdin
}

func countAll(data []byte, resultsChan chan CountResults) {
	var results CountResults

	// Create channels to collect result from each count functions
	bytesChan := make(chan int64)
	wordsChan := make(chan int)
	linesChan := make(chan int)
	charactersChan := make(chan int)

	// Launch goroutines
	go func() {
		bytesChan <- countBytes(bytes.NewReader(data))
	}()
	go func() {
		wordsChan <- countWords(bytes.NewReader(data))
	}()
	go func() {
		linesChan <- countLines(bytes.NewReader(data))
	}()
	go func() {
		charactersChan <- countCharacters(bytes.NewReader(data))
	}()

	// Collect results from the channels
	results.bytesCount = <-bytesChan
	results.wordsCount = <-wordsChan
	results.linesCount = <-linesChan
	results.charactersCount = <-charactersChan

	// Send the combined results back to the main goroutine
	resultsChan <- results
}

func countLines(reader io.Reader) int {
	// Create a new Scanner
	scanner := bufio.NewScanner(reader)
	lineCounter := 0
	for scanner.Scan() {
		lineCounter++
	}

	return lineCounter

}

func countWords(reader io.Reader) int {
	// Create a new Scanner
	scanner := bufio.NewScanner(reader)

	wordCounter := 0
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Fields(line)
		wordCounter += len(words)

	}

	return wordCounter
}

func countCharacters(reader io.Reader) int {
	// Mimicking wc -m
	anotherScanner := bufio.NewScanner(reader)
	anotherScanner.Split(bufio.ScanRunes)
	runesCounter := 0
	for anotherScanner.Scan() {
		runesCounter++
	}
	return runesCounter
}

func countBytes(reader io.Reader) int64 {
	var bc ByteCounter
	byteCounter := &bc
	_, err := io.Copy(byteCounter, reader)
	if err != nil {
		fmt.Printf("lyswc: byteCounter error: %s", err.Error())
		os.Exit(1)
	}
	return byteCounter.count
}
