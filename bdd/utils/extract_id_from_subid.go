package utils

import "strings"

func ExtractIdFromSubid(subid string) string {
	id, _, _ := strings.Cut(subid, "-")
	return id
}
