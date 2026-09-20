package main

import (
	"fmt"
	"os"

	"github.com/edheltzel/iris/internal/cli"
)

var version = "dev"

func main() {
	if err := cli.Run(os.Args[1:], version); err != nil {
		if exit, ok := err.(interface{ ExitCode() int }); ok {
			os.Exit(exit.ExitCode())
		}
		fmt.Fprintln(os.Stderr, "iris:", err)
		os.Exit(1)
	}
}
