package utils

func ConvertStructToPointersArray[E any](array []E) []*E {
	result := make([]*E, 0, len(array))
	for i := range array {
		result = append(result, &array[i])
	}
	return result
}
