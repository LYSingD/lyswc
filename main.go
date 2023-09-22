package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello, World!")

	// flag.*() returns a Pointer
	wordPtr := flag.String("word", "foo", "a string")
	numbPtr := flag.Int("numb", 42, "an int")
	forkPtr := flag.Bool("fork", false, "a bool")

	var svar string
	flag.StringVar(&svar, "svar", "bar", "a string var")

	flag.Parse()

	fmt.Println("word:", *wordPtr)
	fmt.Println("numb:", *numbPtr)
	fmt.Println("fork:", *forkPtr)
	fmt.Println("svar:", svar)
	fmt.Println("tail:", flag.Args())

	// os.Args is a slice of strings that contains the command-line arguments.
	// os.Args[0] is always the program's name/path itself.
	args := flag.Args()
	fmt.Println(args)

	if len(args) < 1 {
		fmt.Println("Usage: lyswc <filepath>")
		os.Exit(1)
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
