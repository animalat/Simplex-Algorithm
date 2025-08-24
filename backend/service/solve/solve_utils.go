package solve

import "strconv"

// float to string
func ftos(num float64) string {
	return strconv.FormatFloat(num, 'g', -1, 64)
}

func getTableInverse(table map[string]int) map[int]string {
	inverse := make(map[int]string)
	for key := range table {
		inverse[table[key]] = key
	}

	return inverse
}
