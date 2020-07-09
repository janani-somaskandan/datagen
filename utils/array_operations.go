package utils

/*
Util for array operations
*/

import(
	"math"
)

func FindMin(array []float64) float64 {

	min := math.MaxFloat64
	for _, element := range array {
		if element < min {
			min = element
		}
	}
	return min
}

func Contains(a []string, x string) bool {
        for _, n := range a {
                if x == n {
                        return true
                }
        }
        return false
}