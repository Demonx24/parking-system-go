package request

import "encoding/xml"

type UnifiedOrderRequest struct {
	XMLName        xml.Name `xml:"xml"`
	AppID          string   `xml:"appid"`
	MchID          string   `xml:"mch_id"`
	NonceStr       string   `xml:"nonce_str"`
	Sign           string   `xml:"sign"`
	Body           string   `xml:"body"`
	OutTradeNo     string   `xml:"out_trade_no"`
	TotalFee       int      `xml:"total_fee"`        // 分为单位，必须是整数
	SpbillCreateIP string   `xml:"spbill_create_ip"` // 用户端ip
	NotifyURL      string   `xml:"notify_url"`
	TradeType      string   `xml:"trade_type"` // JSAPI，NATIVE，APP等
}
type Request struct {
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
}
