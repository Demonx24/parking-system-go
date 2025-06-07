package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// generateSign 生成微信支付签名
// 微信规则：1. 排序；2. 过滤空值和 sign 字段；3. 拼接成 query string；4. 追加 key；5. MD5 加密后转大写

func GenerateSign(params map[string]string, apiKey string) string {
	var keys []string
	for k := range params {
		if params[k] != "" && k != "sign" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys) // 按照 ASCII 排序 key

	var signParts []string
	for _, k := range keys {
		signParts = append(signParts, fmt.Sprintf("%s=%s", k, params[k]))
	}

	// 拼接字符串并附加 key
	signString := strings.Join(signParts, "&") + "&key=" + apiKey

	// MD5 加密
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signString))
	sign := md5Ctx.Sum(nil)

	return strings.ToUpper(hex.EncodeToString(sign)) // 返回大写的 MD5 值

}
