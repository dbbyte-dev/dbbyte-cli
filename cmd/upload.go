package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/dbbyte-dev/dbbyte-go"
)

// Upload implements: dbbyte upload <bucket> <key> <file> [flags]
func Upload(args []string) error {
	fs := flag.NewFlagSet("upload", flag.ExitOnError)
	getClient := newSharedFlags(fs)
	contentType := fs.String("content-type", "", "MIME type of the file (auto-detected by the server if omitted)")
	overwrite := fs.Bool("overwrite", false, "Overwrite an existing object with the same key")

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: dbbyte upload <bucket> <key> <file> [flags]")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	positional := fs.Args()
	if len(positional) != 3 {
		fs.Usage()
		return fmt.Errorf("expected exactly 3 positional arguments: <bucket> <key> <file>")
	}
	bucket, key, filePath := positional[0], positional[1], positional[2]

	dbb, err := getClient()
	if err != nil {
		return err
	}

	body, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", filePath, err)
	}

	result, err := dbb.Objects.Upload(context.Background(), bucket, dbbyte.UploadOptions{
		Key:         key,
		Body:        body,
		ContentType: *contentType,
		Overwrite:   *overwrite,
	})
	if err != nil {
		return err
	}
	return printJSON(result)
}
