package main

import (
	"fmt"
	"os"

	"github.com/dirien/minectl/cmd/minectl"
)

var (
	version string
	commit  string
	date    string
)

func main() {
	if err := minectl.Execute(version, commit, date); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}
