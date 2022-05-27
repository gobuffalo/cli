package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/gobuffalo/cli/cmd/cli"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	// get the present working directory. (PWD)
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = cli.DefaultApp.Run(ctx, pwd, os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}

	<-ctx.Done()
}
