package cmd

import (
	"encoding/json"
	"fmt"
)

func printJSON(v interface{}) error {
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}
