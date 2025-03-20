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

func IsBuiltin(command string) bool {
	for _, value := range BUILTINS {
		if command == value {
			return true
		}
	}
	return false
}

func IsTooManyArgs(nargs int, max int) bool {
	return nargs > max
}

func HandleTooManyArgs() {
	fmt.Println("Too many arguments supplied.")
}

func RunExecutableCommand(command string, args []string) {
	_, err := exec.LookPath(command)
	if err != nil {
		fmt.Println(command + ": command not found")
		return
	}
	cmd := exec.Command(command, args...)
	stdout, err := cmd.Output()
	if err != nil {
		return
	}
	fmt.Printf("%s", stdout)
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

	// Exits with supplied code
	case "exit":
		if IsTooManyArgs(nargs, 1) {
			HandleTooManyArgs()
			return
		}

		exitcode, err := strconv.Atoi(arguments)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error converting error code:", err)
			os.Exit(1)
		}
		os.Exit(exitcode)

	// Prints every argument to StdOut as they are supplied
	case "echo":
		fmt.Printf("%s\n", arguments)

	// pwd return the present working directory
	case "pwd":
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error getting pwd:", err)
		}
		fmt.Printf("%s\n", pwd)

	case "cd":
		if IsTooManyArgs(nargs, 2) {
			HandleTooManyArgs()
			return
		}
		if stat, err := os.Stat(arguments); err == nil && stat.IsDir() {
			os.Chdir(arguments)
		} else if err == nil && !stat.IsDir() {
			fmt.Printf("%s is not a directory", arguments)
		} else {
			fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", arguments)
		}

	// type command that returns either if its a builtin command or where
	// the external executable is found if it is in PATH
	case "type":
		if IsTooManyArgs(nargs, 1) {
			HandleTooManyArgs()
			return
		}

		if IsBuiltin(arguments) {
			fmt.Println(arguments + " is a shell builtin")
			return
		}

		path, err := exec.LookPath(arguments)
		if err != nil {
			fmt.Println(arguments + ": not found")
			return
		}
		fmt.Printf("%s is %s\n", arguments, path)

	// If not builtin command, run command as executable with arguments
	default:
		RunExecutableCommand(command, strings.Split(arguments, " "))
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
