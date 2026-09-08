package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
)

// Download implements: dbbyte download <bucket> <key> <output-file> [flags]
func Download(args []string) error {
	fs := flag.NewFlagSet("download", flag.ExitOnError)
	getClient := newSharedFlags(fs)
	version := fs.String("version", "", `Specific version to fetch, or "all" to download a zip of every version`)

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: dbbyte download <bucket> <key> <output-file> [flags]")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	positional := fs.Args()
	if len(positional) != 3 {
		fs.Usage()
		return fmt.Errorf("expected exactly 3 positional arguments: <bucket> <key> <output-file>")
	}
	bucket, key, outputPath := positional[0], positional[1], positional[2]

	dbb, err := getClient()
	if err != nil {
		return err
	}

	dl, err := dbb.Objects.Download(context.Background(), bucket, key, *version)
	if err != nil {
		return err
	}

	if err := os.WriteFile(outputPath, dl.Data, 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", outputPath, err)
	}

	fmt.Printf("wrote %d bytes to %s (content-type: %s)\n", len(dl.Data), outputPath, dl.ContentType)
	return nil
}
