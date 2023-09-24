package main

import (
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

	counterPtr := mainFs.Bool("c", false, "The number of bytes in each input file is written to the standard output.")

	// Parse the command-line arguments with the custom FlagSet
	err := mainFs.Parse(os.Args[1:])
	fmt.Println("tail:", mainFs.Args())

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

	fmt.Println("You entered:", args)

	// Process the input (example: reverse it)
	input := args[0]
	reversed_input := reverse(input)

	fmt.Println(reversed_input)
}

func reverse(s string) string {
	reversed_s := []rune(s)
	for i, j := 0, len(reversed_s)-1; i < j; i, j = i+1, j-1 {
		reversed_s[i], reversed_s[j] = reversed_s[j], reversed_s[i]
	}
	return string(reversed_s)
}
