package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
)

// Health implements: dbbyte health [flags]
func Health(args []string) error {
	fs := flag.NewFlagSet("health", flag.ExitOnError)
	getClient := newSharedFlags(fs)

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: dbbyte health [flags]")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	dbb, err := getClient()
	if err != nil {
		return err
	}

	result, err := dbb.Health(context.Background())
	if err != nil {
		return err
	}
	return printJSON(result)
}
