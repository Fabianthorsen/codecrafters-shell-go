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

const SINGLE_QUOTE = byte('\'')
const DOUBLE_QUOTE = byte('"')
const BACKSLASH = byte('\\')
const DOLLAR = byte('$')
const SPACE = byte(' ')
const REDIRECT = byte('>')
const STDOUT = byte('1')

var BUILTINS = [5]string{"exit", "echo", "type", "pwd", "cd"}
var EMPTY_ARGV = []string{}

func IsBuiltin(command string) bool {
	for _, value := range BUILTINS {
		if command == value {
			return true
		}
	}
	return false
}

func HasArgsN(nargs int, n int) bool {
	return nargs == n
}

func HandleWrongNumberOfArgs() {
	fmt.Println("Wrong number of arguments passed.")
}

func RunExecutableCmd(argc string, argv []string, redirect bool, output string) {
	_, err := exec.LookPath(argc)
	if err != nil {
		fmt.Println(argc + ": command not found")
		return
	}
	cmd := exec.Command(argc, argv...)
	stdout, err := cmd.Output()
	if err != nil {
		for _, arg := range argv {
			if _, err := os.Stat(arg); os.IsNotExist(err) {
				fmt.Printf("%s: %s: No such file or directory\n", argc, arg)
			}
		}
	}
	if redirect {
		os.WriteFile(output, stdout, 0666)
	} else {
		fmt.Printf("%s", stdout)
	}
}

func CheckFileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func ChangeDirectory(arg string) {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting user home directory:", err)
	}
	arg = strings.Replace(arg, "~", home, 1)
	if stat, err := os.Stat(arg); err == nil && stat.IsDir() {
		os.Chdir(arg)
	} else if err == nil && !stat.IsDir() {
		fmt.Printf("%s is not a directory", arg)
	} else {
		fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", arg)
	}
}

func GetProgramType(argument string) {
	path, err := exec.LookPath(argument)
	if err != nil {
		fmt.Println(argument + ": not found")
		return
	}
	fmt.Printf("%s is %s\n", argument, path)
}

func HandleCommand(argc string, argv []string, redirect bool, output string) {
	nargs := len(argv)

	switch argc {

	// Exits with supplied code
	case "exit":
		if !HasArgsN(nargs, 1) {
			HandleWrongNumberOfArgs()
			return
		}

		exitcode, err := strconv.Atoi(argv[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error converting error code:", err)
			os.Exit(1)
		}
		os.Exit(exitcode)

	// Prints every argument to StdOut as they are supplied
	case "echo":
		if redirect {
			os.WriteFile(output, []byte(strings.Join(argv, " ")+"\n"), 0666)
		} else {
			fmt.Printf("%s\n", strings.Join(argv, " "))
		}

	// pwd return the present working directory
	case "pwd":
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error getting pwd:", err)
			return
		}
		fmt.Printf("%s\n", pwd)

	// cd will change directory if argument is a directory and exists
	case "cd":
		if !HasArgsN(nargs, 1) {
			HandleWrongNumberOfArgs()
			return
		}
		ChangeDirectory(argv[0])

	// type command that returns either if its a builtin command or where
	// the external executable is found if it is in PATH
	case "type":
		if !HasArgsN(nargs, 1) {
			HandleWrongNumberOfArgs()
			return
		}

		arg := argv[0]
		if IsBuiltin(arg) {
			fmt.Println(arg + " is a shell builtin")
		} else {
			GetProgramType(arg)
		}

	// If not builtin command, run command as executable with arguments
	default:
		RunExecutableCmd(argc, argv, redirect, output)
	}
}

func ExtractQuotedSingle(input string) string {
	var sb strings.Builder
	for _, ch := range input {
		if ch == rune(SINGLE_QUOTE) {
			break
		}
		sb.WriteRune(ch)
	}
	return sb.String()
}

func ExtractQuotedDouble(input string) (string, int) {
	i := 0
	var sb strings.Builder
	for i < len(input) {
		switch input[i] {
		case DOUBLE_QUOTE:
			return sb.String(), i
		case BACKSLASH:
			next := input[i+1]
			if next == BACKSLASH || next == DOLLAR || next == DOUBLE_QUOTE {
				next := input[i+1]
				sb.WriteByte(next)
				i += 2
				continue
			}
			fallthrough
		default:
			sb.WriteByte(input[i])
			i++
		}
	}
	return sb.String(), i
}

func ExtractOutput(input string) string {
	trimmed := strings.TrimSpace(input)
	var sb strings.Builder
	for j := 0; j < len(trimmed); j++ {
		sb.WriteByte(trimmed[j])
	}
	return sb.String()
}

func addWordToBuilder(sb *strings.Builder, word string) (int, error) {
	len, err := sb.WriteString(word)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to write string to builder:", err)
	}
	return len, err
}

func ParseArgs(input string) (string, []string, bool, string) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", EMPTY_ARGV, false, ""
	}

	args := []string{}
	redirect := false
	var output strings.Builder
	var sb strings.Builder
	i := 0
	for i < len(trimmed) {
		current := trimmed[i]
		switch current {
		case SINGLE_QUOTE:
			word := ExtractQuotedSingle(trimmed[i+1:])
			l, _ := addWordToBuilder(&sb, word)
			i += l + 1
		case DOUBLE_QUOTE:
			word, l := ExtractQuotedDouble(trimmed[i+1:])
			addWordToBuilder(&sb, word)
			i += l + 1
		case SPACE:
			if sb.Len() > 0 {
				args = append(args, sb.String())
				sb.Reset()
			}
		case BACKSLASH:
			sb.WriteByte(trimmed[i+1])
			i++
		case REDIRECT:
			redirect = true
			path := ExtractOutput(trimmed[i+1:])
			addWordToBuilder(&output, path)
			i = len(trimmed)
		case STDOUT:
			if i < len(trimmed)-1 && trimmed[i+1] == REDIRECT {
				redirect = true
				path := ExtractOutput(trimmed[i+2:])
				addWordToBuilder(&output, path)
				i = len(trimmed)
				continue
			}
			fallthrough
		default:
			sb.WriteByte(current)
			if i == len(trimmed)-1 {
				args = append(args, sb.String())
				sb.Reset()
			}
		}
		i++
	}

	if sb.Len() > 0 {
		args = append(args, sb.String())
	}

	return args[0], args[1:], redirect, output.String()
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		argc, argv, redirect, output := ParseArgs(input)
		if argc == "" {
			continue
		}

		HandleCommand(argc, argv, redirect, output)
	}
}
