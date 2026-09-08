package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/dbbyte-dev/dbbyte-go"
)

// List implements: dbbyte list <bucket> [flags]
func List(args []string) error {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	getClient := newSharedFlags(fs)
	prefix := fs.String("prefix", "", "Filter by key prefix (matches the server-generated key, not the uploaded name - see README)")
	cursor := fs.String("cursor", "", "Pagination cursor from a previous response's nextCursor")
	limit := fs.Int("limit", 0, "Max results to return (server default if omitted)")

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: dbbyte list <bucket> [flags]")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	positional := fs.Args()
	if len(positional) != 1 {
		fs.Usage()
		return fmt.Errorf("expected exactly 1 positional argument: <bucket>")
	}
	bucket := positional[0]

	dbb, err := getClient()
	if err != nil {
		return err
	}

	result, err := dbb.Objects.List(context.Background(), bucket, dbbyte.ListObjectsOptions{
		Prefix: *prefix,
		Cursor: *cursor,
		Limit:  *limit,
	})
	if err != nil {
		return err
	}
	return printJSON(result)
}
