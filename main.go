// Command dbbyte is the official CLI for the DBByte storage/CDN API.
package main

import (
	"fmt"
	"os"

	"github.com/dbbyte-dev/dbbyte-cli/cmd"
)

func usage() {
	fmt.Fprintln(os.Stderr, `dbbyte - DBByte storage/CDN CLI

Usage:
  dbbyte <command> [arguments]

Commands:
  upload      Upload an object to a bucket
  list        List objects in a bucket
  download    Download an object
  delete      Delete an object (or all its versions)
  sign        Create a signed, time-limited URL for an object
  versions    List the version history of an object
  health      Check API health

Auth (required for every command except help):
  --api-key <key>       or the DBBYTE_API_KEY environment variable
  --base-url <url>      or the DBBYTE_BASE_URL environment variable (optional)

Run "dbbyte <command> -h" for command-specific flags.`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	var err error
	switch command {
	case "upload":
		err = cmd.Upload(args)
	case "list":
		err = cmd.List(args)
	case "download":
		err = cmd.Download(args)
	case "delete":
		err = cmd.Delete(args)
	case "sign":
		err = cmd.Sign(args)
	case "versions":
		err = cmd.Versions(args)
	case "health":
		err = cmd.Health(args)
	case "-h", "--help", "help":
		usage()
		return
	default:
		fmt.Fprintf(os.Stderr, "dbbyte: unknown command %q\n\n", command)
		usage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "dbbyte: %s\n", err)
		os.Exit(1)
	}
}
