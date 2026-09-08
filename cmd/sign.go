package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/dbbyte-dev/dbbyte-go"
)

// Sign implements: dbbyte sign <bucket> <key> [flags]
func Sign(args []string) error {
	fs := flag.NewFlagSet("sign", flag.ExitOnError)
	getClient := newSharedFlags(fs)
	version := fs.String("version", "", "Specific version to sign (defaults to the latest)")
	expiresIn := fs.Int("expires-in", 0, "Seconds until the signed URL expires (server default if omitted)")

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: dbbyte sign <bucket> <key> [flags]")
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

	result, err := dbb.Objects.Sign(context.Background(), bucket, key, dbbyte.SignOptions{
		Version:   *version,
		ExpiresIn: *expiresIn,
	})
	if err != nil {
		return err
	}
	return printJSON(result)
}
