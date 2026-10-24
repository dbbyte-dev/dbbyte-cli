# dbbyte-cli

Official command-line tool for the DBByte storage/CDN API, built on top of
[`dbbyte-go`](https://github.com/dbbyte-dev/dbbyte-go).

## Install

```bash
go install github.com/dbbyte-dev/dbbyte-cli@latest
```


## Auth

Every command needs an API key, via either:

```bash
export DBBYTE_API_KEY=dbk_...
```

or the `--api-key` flag on any command. `--base-url` / `DBBYTE_BASE_URL`
override the API host — useful for pointing at a local dev server.

## Commands

```bash
dbbyte upload <bucket> <key> <file> [--content-type TYPE] [--overwrite]
dbbyte list <bucket> [--prefix PREFIX] [--cursor CURSOR] [--limit N]
dbbyte download <bucket> <key> <output-file> [--version V]
dbbyte delete <bucket> <key> [--version V] [--all-versions]
dbbyte sign <bucket> <key> [--version V] [--expires-in SECONDS]
dbbyte versions <bucket> <key> [--limit N] [--cursor CURSOR]
dbbyte upload-url <bucket> <name> [--content-type TYPE] [--size BYTES] [--overwrite] [--expires-in SECONDS]
dbbyte purge --scope file|bucket [--key NAME] <bucket>
dbbyte buckets list
dbbyte buckets get <slug>
dbbyte buckets create --region REGION [--replica-count N] [--public] [--cdn] <name>
dbbyte buckets update [--versioning=BOOL] [--public=BOOL] [--cdn=BOOL] <slug>
dbbyte buckets delete --yes <slug>
dbbyte data get [--fields a,b] [--limit N] [--cursor C] [--filter field=value] <slug>
dbbyte data create (--json '{...}' | --file record.json) <slug>
dbbyte data append-text (--text 'words' | --file notes.txt) <slug>
dbbyte data update (--json '{...}' | --file patch.json) <slug> <record-id>
dbbyte data delete <slug> <record-id>
dbbyte health
```

Flags must come before positional arguments (Go flag parsing). All commands print JSON to stdout on success. On failure, the error
message from the API (and hint, if present) is printed to stderr and the
process exits non-zero.

### Examples

```bash
dbbyte upload my-bucket images/logo.png ./logo.png --content-type image/png

dbbyte list my-bucket --limit 20

dbbyte download my-bucket images/logo.png ./downloaded-logo.png

dbbyte sign my-bucket images/logo.png --expires-in 3600

dbbyte upload-url --content-type image/png --size 2000000 my-bucket photo.png

dbbyte purge --scope file --key logo.png my-bucket
dbbyte purge --scope bucket my-bucket

dbbyte buckets update --cdn=false my-bucket
```

## Note on `list --prefix`

As of the current API, `--prefix` filters on the storage-generated `key`
field (a timestamp-prefixed, slash-flattened server ID), not the
human-readable `name` you uploaded with. It will rarely match what you
expect — treat as unreliable until the backend filters on `name` instead.
