package watchup

import (
	"math"

	"github.com/fatih/color"
)

// calculateAverage calculates the average response time
func calculateAverage(responseTimes []int64) int64 {
	total := int64(0)
	for _, rt := range responseTimes {
		total += rt
	}
	return total / int64(len(responseTimes))
}

// calculateJitter calculates the jitter given a list of response times and the average response time
func calculateJitter(responseTimes []int64, avgResponseTime int64) int64 {
	sumSquaredDifferences := 0.0
	for _, rt := range responseTimes {
		difference := rt - avgResponseTime
		sumSquaredDifferences += math.Pow(float64(difference), 2)
	}
	meanSquaredDifference := sumSquaredDifferences / float64(len(responseTimes))
	jitter := math.Sqrt(meanSquaredDifference)
	return int64(jitter)
}

func FormatStatusCode(code int) string {
	switch {
	case code < 100:
		return color.HiRedString("%d", code)
	case code >= 200 && code < 300:
		return color.GreenString("%d", code)
	case code >= 300 && code < 400:
		return color.YellowString("%d", code)
	case code >= 400 && code < 500:
		return color.CyanString("%d", code)
	case code >= 500:
		return color.HiRedString("%d", code)
	default:
		return color.HiRedString("%d", code)
	}
}
