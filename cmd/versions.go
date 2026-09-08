package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/dbbyte-dev/dbbyte-go"
)

// Versions implements: dbbyte versions <bucket> <key> [flags]
func Versions(args []string) error {
	fs := flag.NewFlagSet("versions", flag.ExitOnError)
	getClient := newSharedFlags(fs)
	cursor := fs.String("cursor", "", "Pagination cursor from a previous response's nextCursor")
	limit := fs.Int("limit", 0, "Max results to return (server default if omitted)")

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: dbbyte versions <bucket> <key> [flags]")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	positional := fs.Args()
	if len(positional) != 2 {
		fs.Usage()
		return fmt.Errorf("expected exactly 2 positional arguments: <bucket> <key>")
	}
	bucket, key := positional[0], positional[1]

	dbb, err := getClient()
	if err != nil {
		return err
	}

	result, err := dbb.Objects.Versions(context.Background(), bucket, key, dbbyte.ListVersionsOptions{
		Cursor: *cursor,
		Limit:  *limit,
	})
	if err != nil {
		return err
	}
	return printJSON(result)
}
