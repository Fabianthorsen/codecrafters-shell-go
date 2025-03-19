package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint
var builtins = [...]string{"exit", "echo", "type"}

func HandleCommandNotFound(command string) {
	fmt.Println(command + ": command not found")
}

func HandleNotFound(argument string) {
	fmt.Println(argument + ": not found")
}

func HandleInput(input string) {
	split := strings.SplitN(input, " ", 2)
	command := split[0]

	arguments := ""
	if len(split) > 1 {
		arguments = split[1]
	}

	nargs := len(strings.Split(arguments, " "))

	switch command {
	case "exit":
		if nargs > 1 {
			fmt.Println("Too many arguments supplied.")
			return
		}
		exitcode, err := strconv.Atoi(arguments)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error converting error code:", err)
			os.Exit(1)
		}
		os.Exit(exitcode)
	case "echo":
		fmt.Printf("%s\n", arguments)
	case "type":
		if nargs > 1 {
			fmt.Println("Too many arguments supplied.")
			return
		}
		for _, value := range builtins {
			if arguments == value {
				fmt.Printf("%s is a shell builtin\n", arguments)
				return
			}
		}
		HandleNotFound(arguments)
	default:
		HandleCommandNotFound(command)
	}
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		cleanedInput := strings.TrimSpace(command)
		if cleanedInput == "" {
			continue
		}

		HandleInput(cleanedInput)
	}
}
