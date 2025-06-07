package service

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"parking-system-go/global"
	"parking-system-go/model/database"
	"parking-system-go/model/request"
	"parking-system-go/model/response"
	"parking-system-go/utils"
	"time"
)

type PayWeChatService struct{}

func getNonceStr() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// 用于模拟调用微信“统一下单”接口的流程，生成二维码支付链接
func (p *PayWeChatService) CreatePayment(order *database.Order, clientIP string) (codeURL string, prepayID string, err error) {
	WeChat := global.Config.WeChat

	// 金额转为分
	totalFee := int(order.Amount * 100)

	// 创建模拟订单
	if err := ServiceGroupApp.CreateMockOrder(order); err != nil {
		return "", "", err
	}

	nonceStr := getNonceStr()

	params := map[string]string{
		"appid":            WeChat.WeChatAppID,
		"mch_id":           WeChat.WeChatMchID,
		"nonce_str":        nonceStr,
		"body":             "停车场支付订单:" + order.OrderID,
		"out_trade_no":     order.OrderID,
		"total_fee":        fmt.Sprintf("%d", totalFee),
		"spbill_create_ip": clientIP,
		"notify_url":       WeChat.WeChatNotifyURL,
		"trade_type":       "NATIVE",
	}

	// 生成签名
	sign := utils.GenerateSign(params, WeChat.WeChatAPIKey)
	params["sign"] = sign

	// 组装 XML 请求体
	reqXML, _ := xml.Marshal(request.UnifiedOrderRequest{
		AppID:          params["appid"],
		MchID:          params["mch_id"],
		NonceStr:       params["nonce_str"],
		Sign:           params["sign"],
		Body:           params["body"],
		OutTradeNo:     params["out_trade_no"],
		TotalFee:       totalFee,
		SpbillCreateIP: params["spbill_create_ip"],
		NotifyURL:      params["notify_url"],
		TradeType:      params["trade_type"],
	})

	resp, err := http.Post("http://localhost:9090/pay/unifiedorder", "text/xml", bytes.NewReader(reqXML))
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var wxResp response.UnifiedOrderResponse
	if err := xml.Unmarshal(respBody, &wxResp); err != nil {
		return "", "", err
	}

	return wxResp.CodeURL, wxResp.PrepayID, nil
}
