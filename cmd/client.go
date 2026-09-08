package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/dbbyte-dev/dbbyte-go"
)

// newSharedFlags registers the auth flags common to every subcommand and
// returns a function to build a client after fs.Parse() has run.
func newSharedFlags(fs *flag.FlagSet) func() (*dbbyte.DBByte, error) {
	apiKeyFlag := fs.String("api-key", "", "API key (defaults to DBBYTE_API_KEY env var)")
	baseURLFlag := fs.String("base-url", "", "API base URL (defaults to DBBYTE_BASE_URL env var, or the production API)")

	return func() (*dbbyte.DBByte, error) {
		apiKey := *apiKeyFlag
		if apiKey == "" {
			apiKey = os.Getenv("DBBYTE_API_KEY")
		}
		if apiKey == "" {
			return nil, fmt.Errorf("no API key: pass --api-key or set DBBYTE_API_KEY")
		}

		opts := dbbyte.ClientOptions{APIKey: apiKey}
		baseURL := *baseURLFlag
		if baseURL == "" {
			baseURL = os.Getenv("DBBYTE_BASE_URL")
		}
		if baseURL != "" {
			opts.BaseURL = baseURL
		}

		return dbbyte.New(opts)
	}
}
