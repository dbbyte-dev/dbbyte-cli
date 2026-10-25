package cmd

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dbbyte-dev/dbbyte-go"
)

const dataUsage = `Usage:
  dbbyte data get <slug> [--fields a,b] [--limit N] [--cursor C] [--filter field=value]...
  dbbyte data create <slug> (--json '{...}' | --file record.json)
  dbbyte data append-text <slug> (--text 'words' | --file notes.txt)
  dbbyte data update <slug> <record-id> (--json '{...}' | --file patch.json)
  dbbyte data delete <slug> <record-id>

Flags go before positional arguments.`

// Data implements: dbbyte data <get|create|append-text|update|delete>
func Data(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, dataUsage)
		return fmt.Errorf("missing data subcommand")
	}
	sub, rest := args[0], args[1:]
	switch sub {
	case "get":
		return dataGet(rest)
	case "create":
		return dataCreate(rest)
	case "append-text":
		return dataAppendText(rest)
	case "update":
		return dataUpdate(rest)
	case "delete":
		return dataDelete(rest)
	case "-h", "--help", "help":
		fmt.Fprintln(os.Stderr, dataUsage)
		return nil
	}
	fmt.Fprintln(os.Stderr, dataUsage)
	return fmt.Errorf("unknown data subcommand %q", sub)
}

type filterList []string

func (f *filterList) String() string     { return strings.Join(*f, ",") }
func (f *filterList) Set(v string) error { *f = append(*f, v); return nil }

func dataGet(args []string) error {
	fs := flag.NewFlagSet("data get", flag.ExitOnError)
	getClient := newSharedFlags(fs)
	fields := fs.String("fields", "", "Comma-separated fields to return")
	limit := fs.Int("limit", 0, "Page size")
	cursor := fs.String("cursor", "", "Cursor from a previous page")
	var filters filterList
	fs.Var(&filters, "filter", "Exact-match filter as field=value (repeatable)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("expected exactly 1 positional argument: <slug>")
	}
	opts := dbbyte.DataQueryOptions{Cursor: *cursor, Limit: *limit}
	if *fields != "" {
		opts.Fields = strings.Split(*fields, ",")
	}
	if len(filters) > 0 {
		opts.Filters = map[string]string{}
		for _, f := range filters {
			name, value, ok := strings.Cut(f, "=")
			if !ok || name == "" {
				return fmt.Errorf("--filter %q must look like field=value", f)
			}
			opts.Filters[name] = value
		}
	}
	dbb, err := getClient()
	if err != nil {
		return err
	}
	result, err := dbb.Data.Get(context.Background(), fs.Arg(0), opts)
	if err != nil {
		return err
	}
	return printJSON(result)
}

// readRecord loads a JSON object from --json or --file (exactly one).
func readRecord(jsonStr, file string) (dbbyte.DataRecord, error) {
	if (jsonStr == "") == (file == "") {
		return nil, fmt.Errorf("pass exactly one of --json or --file")
	}
	raw := []byte(jsonStr)
	if file != "" {
		b, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", file, err)
		}
		raw = b
	}
	var rec dbbyte.DataRecord
	if err := json.Unmarshal(raw, &rec); err != nil || rec == nil {
		return nil, fmt.Errorf("record must be a JSON object")
	}
	return rec, nil
}

func dataCreate(args []string) error {
	fs := flag.NewFlagSet("data create", flag.ExitOnError)
	getClient := newSharedFlags(fs)
	jsonStr := fs.String("json", "", "Record as a JSON object")
	file := fs.String("file", "", "Read the record from a JSON file")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("expected exactly 1 positional argument: <slug>")
	}
	rec, err := readRecord(*jsonStr, *file)
	if err != nil {
		return err
	}
	dbb, err := getClient()
	if err != nil {
		return err
	}
	result, err := dbb.Data.Create(context.Background(), fs.Arg(0), rec)
	if err != nil {
		return err
	}
	return printJSON(result)
}

func dataAppendText(args []string) error {
	fs := flag.NewFlagSet("data append-text", flag.ExitOnError)
	getClient := newSharedFlags(fs)
	text := fs.String("text", "", "Text to append")
	file := fs.String("file", "", "Read the text to append from a file")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("expected exactly 1 positional argument: <slug>")
	}
	if (*text == "") == (*file == "") {
		return fmt.Errorf("pass exactly one of --text or --file")
	}
	body := *text
	if *file != "" {
		b, err := os.ReadFile(*file)
		if err != nil {
			return fmt.Errorf("reading %s: %w", *file, err)
		}
		body = string(b)
	}
	dbb, err := getClient()
	if err != nil {
		return err
	}
	result, err := dbb.Data.AppendText(context.Background(), fs.Arg(0), body)
	if err != nil {
		return err
	}
	return printJSON(result)
}

func dataUpdate(args []string) error {
	fs := flag.NewFlagSet("data update", flag.ExitOnError)
	getClient := newSharedFlags(fs)
	jsonStr := fs.String("json", "", "Fields to change, as a JSON object")
	file := fs.String("file", "", "Read the patch from a JSON file")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 2 {
		return fmt.Errorf("expected exactly 2 positional arguments: <slug> <record-id>")
	}
	patch, err := readRecord(*jsonStr, *file)
	if err != nil {
		return err
	}
	dbb, err := getClient()
	if err != nil {
		return err
	}
	result, err := dbb.Data.Update(context.Background(), fs.Arg(0), fs.Arg(1), patch)
	if err != nil {
		return err
	}
	return printJSON(result)
}

func dataDelete(args []string) error {
	fs := flag.NewFlagSet("data delete", flag.ExitOnError)
	getClient := newSharedFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 2 {
		return fmt.Errorf("expected exactly 2 positional arguments: <slug> <record-id>")
	}
	dbb, err := getClient()
	if err != nil {
		return err
	}
	result, err := dbb.Data.Delete(context.Background(), fs.Arg(0), fs.Arg(1))
	if err != nil {
		return err
	}
	return printJSON(result)
}
