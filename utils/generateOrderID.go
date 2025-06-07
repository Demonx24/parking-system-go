package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateOrderID() string {
	now := time.Now()
	// 格式化时间：YYYYMMDDHHMMSS
	timestamp := now.Format("20060102150405")

	// 生成4位随机数
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(9000) + 1000 // 1000-9999

	return fmt.Sprintf("%s%d", timestamp, randomNum)
}
