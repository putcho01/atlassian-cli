package main

import (
	"os"

	"github.com/putcho01/atlassian-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
