package main

import (
	"fmt"
	"github.com/minectl/cmd/minectl"
	"os"
)

var (
	version   string
	gitCommit string
)

func main() {
	if err := minectl.Execute(version, gitCommit); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
