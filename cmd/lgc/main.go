// Command lgc submits source code to Luogu and observes the judge result.
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/haomingbai/lgc/cmd/internal/arguments"
	"github.com/haomingbai/lgc/cmd/internal/configuration"
	"github.com/haomingbai/lgc/cmd/internal/executor"
	"github.com/haomingbai/lgc/cmd/internal/resources"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	os.Exit(run(ctx, os.Args[1:], os.Stdout, os.Stderr))
}

func run(
	ctx context.Context,
	argv []string,
	stdout io.Writer,
	stderr io.Writer,
) int {
	invocation, err := arguments.NewReader().Read(argv)
	if err != nil {
		fmt.Fprintf(stderr, "lgc: %v\n", err)
		return int(executor.ExitCodeForError(err))
	}

	configurationValue, err := configuration.NewReader().Read(ctx, invocation)
	if err != nil {
		fmt.Fprintf(stderr, "lgc: %v\n", err)
		return int(executor.ExitCodeForError(err))
	}

	registry := resources.NewRegistry()
	bundle, err := resources.NewGenerator(registry).Generate(
		ctx,
		invocation,
		configurationValue,
	)
	if err != nil {
		fmt.Fprintf(stderr, "lgc: %v\n", err)
		return int(executor.ExitCodeForError(err))
	}
	defer bundle.Close()

	_, exitCode := executor.NewExecutor(stdout).Execute(ctx, bundle)
	return int(exitCode)
}
