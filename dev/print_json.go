package dev

import (
	"encoding/json"
	"fmt"
)

func PrintJson(anyStruct ...interface{}) {
	for _, v := range anyStruct {
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(b))
	}
}
