package main

import (
	"fmt"
	"os"

	"prototogozeroapi/internal/cli"
)

func main() {
	if err := cli.Run(os.Args[1:]); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

