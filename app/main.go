package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}
		cmd := strings.SplitN(command[:len(command)-1], " ", 2)
		switch cmd[0] {
		case "exit":
			os.Exit(0)
		case "echo":
			fmt.Printf("%s\n", cmd[1])
		default:
			fmt.Println(cmd[0] + ": command not found")
		}
	}
}
