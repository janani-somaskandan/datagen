package utils

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