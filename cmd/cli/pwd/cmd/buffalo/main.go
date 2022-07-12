package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gobuffalo/cli/cmd/cli"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}

	err = cli.DefaultApp.Main(context.Background(), pwd, os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}
