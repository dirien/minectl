package main

import (
	"fmt"
	"github.com/minectl/cmd/minectl"
	"os"
)

// These values will be injected into these variables at the build time.
var (
	Version   string
	GitCommit string
)

func main() {
	if err := minectl.Execute(Version, GitCommit); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
