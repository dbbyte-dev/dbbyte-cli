module github.com/dbbyte-dev/dbbyte-cli

go 1.21

require github.com/dbbyte-dev/dbbyte-go v0.1.0

// Remove this once dbbyte-go is published/tagged - points at a sibling
// checkout for local development.
replace github.com/dbbyte-dev/dbbyte-go => ../dbbyte-go
