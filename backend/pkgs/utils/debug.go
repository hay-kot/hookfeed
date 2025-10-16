package utils

import (
	"encoding/json"
	"fmt"
)

// Dump dumps the struct/contents via JSON to the console. Only to be used for debugging.
func Dump(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
