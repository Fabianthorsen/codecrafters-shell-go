package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint
var BUILTINS = [...]string{"exit", "echo", "type", "pwd"}
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
	return err == nil
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
	case "pwd":
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error getting pwd:", err)
		}
		fmt.Printf("%s\n", pwd)
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
		path, err := exec.LookPath(arguments)
		if err != nil {
			HandleNotFound(arguments)
			return
		}
		HandleInPath(arguments, path)

	default:
		_, err := exec.LookPath(command)
		if err != nil {
			HandleCommandNotFound(command)
			return
		}
		args := strings.Split(arguments, " ")
		cmd := exec.Command(command, args...)
		stdout, err := cmd.Output()
		if err != nil {
			return
		}
		fmt.Printf("%s", stdout)
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
