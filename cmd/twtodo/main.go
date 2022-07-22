package main

import (
	"os"

	"github.com/simachri/taskwarrior-ms-todo/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
