package server

import (
	"gonum.org/v1/gonum/stat"
	"strconv"
)

func identifyOutlierTimeUnits(unit string, values []float64, c chan string) {
	mean := stat.Mean(values, nil)
	stdDev := stat.StdDev(values, nil)
	// Set a threshold for identifying outliers
	threshold := 2.0

	// Identify outliers based on z-score
	for _, value := range values {
		zScore := (value - mean) / stdDev
		if zScore > threshold {
			c <- strconv.Itoa(int(value)) + " " + unit
		}
	}
}
