package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unicode"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint
var BUILTINS = [...]string{"exit", "echo", "type", "pwd", "cd"}
var EMPTY_ARGV = []string{}
var SINGLE_QUOTE = '\''

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

func RunExecutableCmd(argc string, argv []string) {
	_, err := exec.LookPath(argc)
	if err != nil {
		fmt.Println(argc + ": command not found")
		return
	}
	fmt.Printf("%s called with %v\n. %d", argc, strings.Join(argv, " "), len(argv))
	cmd := exec.Command(argc, argv...)
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

func HandleCommand(argc string, argv []string) {
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
		fmt.Printf("%s\n", strings.Join(argv, " "))

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
		RunExecutableCmd(argc, argv)
	}
}

func CleanInput(input string) (string, string) {
	trimmed := strings.TrimSpace(input)
	split := strings.SplitN(trimmed, " ", 2)
	if len(split) == 2 {
		return split[0], split[1]
	}
	return split[0], ""
}

func MakeArgv(argstr string) []string {
	if argstr == "" {
		return EMPTY_ARGV
	}

	argv := []string{}
	quoted := false
	skipLetter := false
	spaceCount := 0
	var sb strings.Builder
	for i, ch := range argstr {

		if skipLetter {
			skipLetter = false
			continue
		}

		// TODO: echo test  world -> testworld
		if ch != SINGLE_QUOTE {
			if quoted {
				sb.WriteRune(ch)
				continue
			} else {
				if unicode.IsSpace(ch) {
					spaceCount++
					if spaceCount < 2 {
						sb.WriteRune(ch)
						argv = append(argv, strings.Trim(sb.String(), " "))
						sb.Reset()
						continue
					}
				} else {
					spaceCount = 0
					sb.WriteRune(ch)
				}
			}
			// TODO: Unquoted works. Need to split quoted on space too between quoted words.
		} else {
			if !quoted {
				if sb.Len() > 0 {
					argv = append(argv, strings.Trim(sb.String(), " "))
					sb.Reset()
				}
				quoted = true
				continue
			} else {
				if i < len(argstr)-1 {
					if argstr[i+1] != byte(SINGLE_QUOTE) {
						argv = append(argv, sb.String())
						sb.Reset()
					} else {
						skipLetter = true
					}
				}
			}
		}

		if i == len(argstr)-1 {
			argv = append(argv, strings.Trim(sb.String(), " "))
		}
	}
	return argv
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		argc, argstr := CleanInput(input)
		if argc == "" {
			continue
		}

		argv := MakeArgv(argstr)

		HandleCommand(argc, argv)
	}
}
