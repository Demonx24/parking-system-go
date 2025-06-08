package utils

import (
	"encoding/json"
	"golang.org/x/net/context"
	"parking-system-go/global"
	"time"
)

var ctx = context.Background()

// GetOrSetStruct 从 Redis 获取结构体，如果不存在则执行 fetchFunc 获取并缓存
func GetOrSetStruct[T any](key string, ttl time.Duration, fetchFunc func() (T, error)) (T, error) {
	var result T
	client := global.Redis
	// 尝试从 Redis 获取
	val, err := client.Get(ctx, key).Result()
	if err == nil && val != "" {
		err = json.Unmarshal([]byte(val), &result)
		if err == nil {
			return result, nil
		}
	}

	// 如果未命中或反序列化失败，则通过回调获取
	result, err = fetchFunc()
	if err != nil {
		return result, err
	}

	// 序列化并写入 Redis
	data, _ := json.Marshal(result)
	_ = client.Set(ctx, key, data, ttl).Err()

	return result, nil
}
