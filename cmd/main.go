package main

import (
	"fmt"
	"os"
)

var exitCode int = 0

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	os.Exit(exitCode)
}
