package utils

import (
	"cashflow/models"
	"sort"
)

func SortByDate[A models.WithDate](array []A) []A {
	sort.SliceStable(array, func(i, j int) bool {
		return IsDateBefore(array[i].GetDate(), array[j].GetDate())
	})

	return array
}
