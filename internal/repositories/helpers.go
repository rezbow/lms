package repositories

import "math"

func totalPage(total, limit float64) int {
	return int(math.Ceil(total / limit))
}
