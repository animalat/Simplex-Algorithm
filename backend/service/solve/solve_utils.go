package solve

import (
	"math"
	"strconv"
)

// float to string
func ftos(num float64) string {
	return strconv.FormatFloat(num, 'g', -1, 64)
}

func floatsEqual(a float64, b float64) bool {
	return math.Abs(a-b) < EPSILON
}

func floatsEqualWithError(a float64, b float64, precisionError float64) bool {
	return math.Abs(a-b) < (EPSILON + precisionError)
}

func getTableInverse(table map[string]int) map[int]string {
	inverse := make(map[int]string)
	for key := range table {
		inverse[table[key]] = key
	}

	return inverse
}
