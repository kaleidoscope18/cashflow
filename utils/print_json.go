package utils

import (
	"encoding/json"
	"fmt"
)

func PrintJson(anyStruct interface{}) {
	b, err := json.MarshalIndent(anyStruct, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(string(b))
}
