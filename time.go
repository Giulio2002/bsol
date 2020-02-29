package main

import (
	"strconv"
	"strings"
)

func convertElapsedToNano(elapsed string) float64 {
	switch {
	case strings.Contains(elapsed, "µs"):
		elapsed = strings.Replace(elapsed, "µs", "", 1)
		result, _ := strconv.ParseFloat(elapsed, 64)
		return result
	case strings.Contains(elapsed, "ms"):
		elapsed = strings.Replace(elapsed, "ms", "", 1)
		result, _ := strconv.ParseFloat(elapsed, 64)
		return result * 1000
	case strings.Contains(elapsed, "s"):
		elapsed = strings.Replace(elapsed, "ms", "", 1)
		result, _ := strconv.ParseFloat(elapsed, 64)
		return result * 1000000
	}
	return 0
}
