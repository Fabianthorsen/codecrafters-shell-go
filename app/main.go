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
var BUILTINS = [...]string{"exit", "echo", "type"}
var PATH = strings.Split(os.Getenv("PATH"), ":")

func HandleCommandNotFound(command string) {
	fmt.Println(command + ": command not found")
}

func HandleNotFound(argument string) {
	fmt.Println(argument + ": not found")
}

func HandleBuiltin(command string) {
	fmt.Println(command + " is a shell builtin")
}

func HandleInPath(command string, filepath string) {
	fmt.Printf("%s is %s\n", command, filepath)
}

func CheckFileExists(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return false
	}
	return true
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
		for _, value := range BUILTINS {
			if arguments == value {
				HandleBuiltin(arguments)
				return
			}
		}
		for _, value := range PATH {
			filepath := fmt.Sprintf("%s/%s", value, arguments)
			if CheckFileExists(filepath) {
				HandleInPath(arguments, filepath)
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
