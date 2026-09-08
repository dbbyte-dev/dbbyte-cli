package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
)

// Delete implements: dbbyte delete <bucket> <key> [flags]
func Delete(args []string) error {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)
	getClient := newSharedFlags(fs)
	version := fs.String("version", "", "Specific version to delete (defaults to the latest)")
	allVersions := fs.Bool("all-versions", false, "Permanently delete every version of this key")

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: dbbyte delete <bucket> <key> [flags]")
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

	ctx := context.Background()
	var result interface{}
	if *allVersions {
		result, err = dbb.Objects.DeleteAllVersions(ctx, bucket, key)
	} else {
		result, err = dbb.Objects.Delete(ctx, bucket, key, *version)
	}
	if err != nil {
		return err
	}
	return printJSON(result)
}
