package utils

import (
	"math/rand"
	"strconv"
	"time"
)

// RandomStringNumber 随机范围数字内的数字字符串
func RandomStringNumber(low, high int) string {
	rand.Seed(time.Now().UnixNano())

	return strconv.Itoa(rand.Intn(high-low+1) + low)
}
