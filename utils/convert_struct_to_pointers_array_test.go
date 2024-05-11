package utils

import "testing"

func TestConvertStructToPointersArray(t *testing.T) {
	one := 1
	two := 2
	three := 3
	data := []int{one, two, three}
	result := ConvertStructToPointersArray(data)

	for i := range data {
		if *result[i] != data[i] {
			t.Fatalf(`Expected pointer %d at index %d but found %d instead.`, &data[i], i, result[i])
		}
	}
}
