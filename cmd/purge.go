package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/dbbyte-dev/dbbyte-go"
)

// Purge implements: dbbyte purge [flags] <bucket>
func Purge(args []string) error {
	fs := flag.NewFlagSet("purge", flag.ExitOnError)
	getClient := newSharedFlags(fs)
	scope := fs.String("scope", "bucket", "What to purge: file or bucket")
	key := fs.String("key", "", "File name to purge (required with --scope file)")

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: dbbyte purge [flags] <bucket>")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return fmt.Errorf("expected exactly 1 positional argument: <bucket>")
	}
	bucket := fs.Arg(0)

	if *scope != "file" && *scope != "bucket" {
		return fmt.Errorf("--scope must be \"file\" or \"bucket\"")
	}
	if *scope == "file" && *key == "" {
		return fmt.Errorf("--key is required with --scope file")
	}
	if *scope == "bucket" && *key != "" {
		return fmt.Errorf("--key can only be used with --scope file")
	}

	dbb, err := getClient()
	if err != nil {
		return err
	}
	result, err := dbb.CDN.Purge(context.Background(), dbbyte.PurgeOptions{
		Slug:  bucket,
		Scope: dbbyte.PurgeScope(*scope),
		Path:  *key,
	})
	if err != nil {
		return err
	}
	return printJSON(result)
}
