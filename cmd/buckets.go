package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/dbbyte-dev/dbbyte-go"
)

const bucketsUsage = `Usage:
  dbbyte buckets list
  dbbyte buckets get <slug>
  dbbyte buckets create <name> --region REGION [--replica-count N] [--public] [--cdn]
  dbbyte buckets update <slug> [--versioning=BOOL] [--public=BOOL] [--cdn=BOOL]
  dbbyte buckets delete <slug>

Flags go before positional arguments.`

// Buckets implements: dbbyte buckets <list|get|create|update|delete>
func Buckets(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, bucketsUsage)
		return fmt.Errorf("missing buckets subcommand")
	}
	sub, rest := args[0], args[1:]
	switch sub {
	case "list":
		return bucketsList(rest)
	case "get":
		return bucketsGet(rest)
	case "create":
		return bucketsCreate(rest)
	case "update":
		return bucketsUpdate(rest)
	case "delete":
		return bucketsDelete(rest)
	case "-h", "--help", "help":
		fmt.Fprintln(os.Stderr, bucketsUsage)
		return nil
	}
	fmt.Fprintln(os.Stderr, bucketsUsage)
	return fmt.Errorf("unknown buckets subcommand %q", sub)
}

func bucketsList(args []string) error {
	fs := flag.NewFlagSet("buckets list", flag.ExitOnError)
	getClient := newSharedFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("buckets list takes no positional arguments")
	}
	dbb, err := getClient()
	if err != nil {
		return err
	}
	result, err := dbb.Buckets.List(context.Background())
	if err != nil {
		return err
	}
	return printJSON(result)
}

func bucketsGet(args []string) error {
	fs := flag.NewFlagSet("buckets get", flag.ExitOnError)
	getClient := newSharedFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("expected exactly 1 positional argument: <slug>")
	}
	dbb, err := getClient()
	if err != nil {
		return err
	}
	result, err := dbb.Buckets.Get(context.Background(), fs.Arg(0))
	if err != nil {
		return err
	}
	return printJSON(result)
}

func bucketsCreate(args []string) error {
	fs := flag.NewFlagSet("buckets create", flag.ExitOnError)
	getClient := newSharedFlags(fs)
	region := fs.String("region", "", "Region for the bucket (required)")
	replicas := fs.Int("replica-count", 0, "Number of replicas (plan-dependent)")
	public := fs.Bool("public", false, "Make the bucket public")
	cdn := fs.Bool("cdn", false, "Enable CDN caching")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("expected exactly 1 positional argument: <name>")
	}
	if *region == "" {
		return fmt.Errorf("--region is required")
	}
	if *replicas < 0 {
		return fmt.Errorf("--replica-count must not be negative")
	}
	dbb, err := getClient()
	if err != nil {
		return err
	}
	result, err := dbb.Buckets.Create(context.Background(), dbbyte.CreateBucketOptions{
		Name:         fs.Arg(0),
		Region:       *region,
		ReplicaCount: *replicas,
		IsPublic:     *public,
		CDN:          *cdn,
	})
	if err != nil {
		return err
	}
	return printJSON(result)
}

func bucketsUpdate(args []string) error {
	fs := flag.NewFlagSet("buckets update", flag.ExitOnError)
	getClient := newSharedFlags(fs)
	versioning := fs.Bool("versioning", false, "Turn versioning on or off (--versioning=false)")
	public := fs.Bool("public", false, "Make the bucket public or private (--public=false)")
	cdn := fs.Bool("cdn", false, "Turn CDN caching on or off (--cdn=false)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("expected exactly 1 positional argument: <slug>")
	}

	var opts dbbyte.UpdateBucketOptions
	changed := 0
	fs.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "versioning":
			opts.VersioningEnabled = versioning
			changed++
		case "public":
			opts.IsPublic = public
			changed++
		case "cdn":
			opts.CDN = cdn
			changed++
		}
	})
	if changed == 0 {
		return fmt.Errorf("nothing to update: pass --versioning, --public, or --cdn")
	}

	dbb, err := getClient()
	if err != nil {
		return err
	}
	result, err := dbb.Buckets.Update(context.Background(), fs.Arg(0), opts)
	if err != nil {
		return err
	}
	return printJSON(result)
}

func bucketsDelete(args []string) error {
	fs := flag.NewFlagSet("buckets delete", flag.ExitOnError)
	getClient := newSharedFlags(fs)
	yes := fs.Bool("yes", false, "Confirm deleting the bucket and everything in it")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("expected exactly 1 positional argument: <slug>")
	}
	if !*yes {
		return fmt.Errorf("deleting a bucket removes its contents permanently: re-run with --yes to confirm")
	}
	dbb, err := getClient()
	if err != nil {
		return err
	}
	result, err := dbb.Buckets.Delete(context.Background(), fs.Arg(0))
	if err != nil {
		return err
	}
	return printJSON(result)
}
