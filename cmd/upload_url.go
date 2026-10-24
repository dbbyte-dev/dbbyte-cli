package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/dbbyte-dev/dbbyte-go"
)

// UploadURL implements: dbbyte upload-url [flags] <bucket> <name>
func UploadURL(args []string) error {
	fs := flag.NewFlagSet("upload-url", flag.ExitOnError)
	getClient := newSharedFlags(fs)
	contentType := fs.String("content-type", "", "Content-Type the browser must send on the PUT")
	size := fs.Int64("size", 0, "Maximum size in bytes the link accepts (up to 50 MB)")
	overwrite := fs.Bool("overwrite", false, "Allow replacing an existing object with the same name")
	expiresIn := fs.Int("expires-in", 0, "Seconds until the link expires (server default if omitted)")

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: dbbyte upload-url [flags] <bucket> <name>")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 2 {
		fs.Usage()
		return fmt.Errorf("expected exactly 2 positional arguments: <bucket> <name>")
	}
	if *size < 0 || *expiresIn < 0 {
		return fmt.Errorf("--size and --expires-in must not be negative")
	}

	dbb, err := getClient()
	if err != nil {
		return err
	}
	result, err := dbb.Objects.UploadURL(context.Background(), fs.Arg(0), dbbyte.UploadURLOptions{
		Name:        fs.Arg(1),
		ContentType: *contentType,
		Size:        *size,
		Overwrite:   *overwrite,
		ExpiresIn:   *expiresIn,
	})
	if err != nil {
		return err
	}
	return printJSON(result)
}
