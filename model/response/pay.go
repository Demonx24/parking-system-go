package response

import "encoding/xml"

// 统一下单响应结构体
type UnifiedOrderResponse struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
	AppID      string   `xml:"appid"`
	MchID      string   `xml:"mch_id"`
	NonceStr   string   `xml:"nonce_str"`
	Sign       string   `xml:"sign"`
	ResultCode string   `xml:"result_code"`
	PrepayID   string   `xml:"prepay_id"`
	CodeURL    string   `xml:"code_url"` // 二维码链接（NATIVE支付时）
}

// 微信支付回调请求结构体
type WechatNotifyRequest struct {
	XMLName       xml.Name `xml:"xml"`
	ReturnCode    string   `xml:"return_code"`
	ReturnMsg     string   `xml:"return_msg"`
	AppID         string   `xml:"appid"`
	MchID         string   `xml:"mch_id"`
	NonceStr      string   `xml:"nonce_str"`
	Sign          string   `xml:"sign"`
	ResultCode    string   `xml:"result_code"`
	OpenID        string   `xml:"openid"`
	TradeType     string   `xml:"trade_type"`
	BankType      string   `xml:"bank_type"`
	TotalFee      int      `xml:"total_fee"`
	CashFee       int      `xml:"cash_fee"`
	TransactionID string   `xml:"transaction_id"`
	OutTradeNo    string   `xml:"out_trade_no"`
	TimeEnd       string   `xml:"time_end"`
}

// 支付成功响应结构体（返回给微信）
type WxPaySuccessResponse struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
}

//type WechatUnifiedOrderResponse struct {
//	ReturnCode string `xml:"return_code"`
//	ReturnMsg  string `xml:"return_msg"`
//
//	AppID      string `xml:"appid"`
//	MchID      string `xml:"mch_id"`
//	NonceStr   string `xml:"nonce_str"`
//	Sign       string `xml:"sign"`
//	ResultCode string `xml:"result_code"`
//
//	PrepayID string `xml:"prepay_id"`
//	CodeURL  string `xml:"code_url"`
//}
