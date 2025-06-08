package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

type UnifiedOrderRequest struct {
	XMLName        xml.Name `xml:"xml"`
	AppID          string   `xml:"appid"`
	MchID          string   `xml:"mch_id"`
	NonceStr       string   `xml:"nonce_str"`
	Sign           string   `xml:"sign"`
	Body           string   `xml:"body"`
	OutTradeNo     string   `xml:"out_trade_no"`
	TotalFee       int      `xml:"total_fee"`
	SpbillCreateIP string   `xml:"spbill_create_ip"`
	NotifyURL      string   `xml:"notify_url"`
	TradeType      string   `xml:"trade_type"`
}

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
	CodeURL    string   `xml:"code_url"`
}

type PayNotifyRequest struct {
	XMLName       xml.Name `xml:"xml"`
	ReturnCode    string   `xml:"return_code"`
	ReturnMsg     string   `xml:"return_msg"`
	AppID         string   `xml:"appid"`
	MchID         string   `xml:"mch_id"`
	NonceStr      string   `xml:"nonce_str"`
	Sign          string   `xml:"sign"`
	ResultCode    string   `xml:"result_code"`
	OutTradeNo    string   `xml:"out_trade_no"`
	TotalFee      int      `xml:"total_fee"`
	TransactionID string   `xml:"transaction_id"`
	TimeEnd       string   `xml:"time_end"`
}

func generateSign(params map[string]string, apiKey string) string {
	var keys []string
	for k := range params {
		if params[k] != "" && k != "sign" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var signParts []string
	for _, k := range keys {
		signParts = append(signParts, fmt.Sprintf("%s=%s", k, params[k]))
	}

	signString := strings.Join(signParts, "&") + "&key=" + apiKey
	hash := md5.Sum([]byte(signString))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

func getNonceStr() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
func generatePrepayID() string {
	timestamp := time.Now().Format("20060102150405")
	randomStr := fmt.Sprintf("%010d", time.Now().UnixNano()%1e10)
	return "wx" + timestamp + randomStr
}

var prepayID = generatePrepayID()

func main() {
	apiKey := "wlc1224"
	appID := "wx1234567890abcdef"
	mchID := "1234567890"
	notifyURL := "http://localhost:8081/api/pay/payment_notify"

	r := gin.Default()

	// 模拟统一下单接口
	r.POST("/pay/unifiedorder", func(c *gin.Context) {
		var req UnifiedOrderRequest
		if err := c.ShouldBindXML(&req); err != nil {
			c.XML(http.StatusBadRequest, gin.H{"error": "参数错误"})
			return
		}

		params := map[string]string{
			"appid":            req.AppID,
			"body":             req.Body,
			"mch_id":           req.MchID,
			"nonce_str":        req.NonceStr,
			"notify_url":       req.NotifyURL,
			"out_trade_no":     req.OutTradeNo,
			"spbill_create_ip": req.SpbillCreateIP,
			"total_fee":        fmt.Sprintf("%d", req.TotalFee),
			"trade_type":       req.TradeType,
		}
		sign := generateSign(params, apiKey)
		if sign != req.Sign {
			c.XML(http.StatusOK, UnifiedOrderResponse{ReturnCode: "FAIL", ReturnMsg: "签名错误"})
			return
		}

		nonceStr := getNonceStr()

		codeURL := "http://swwuam7b1.hn-bkt.clouddn.com/pay.jpg?e=1749205770&token=xWO3ioDnx_B5AhAsQBWwcFNAoqbvEBezEUSIJOea:FJ1GcpLXotconZUNT23YBKRMakw="

		respParams := map[string]string{
			"appid":       appID,
			"mch_id":      mchID,
			"nonce_str":   nonceStr,
			"prepay_id":   prepayID,
			"result_code": "failure",
		}
		respSign := generateSign(respParams, apiKey)

		c.XML(http.StatusOK, UnifiedOrderResponse{
			ReturnCode: "failure",
			ReturnMsg:  "OK",
			AppID:      appID,
			MchID:      mchID,
			NonceStr:   nonceStr,
			Sign:       respSign,
			ResultCode: "failure",
			PrepayID:   prepayID,
			CodeURL:    codeURL,
		})

		go func() {
			time.Sleep(300 * time.Second)

			notify := PayNotifyRequest{
				AppID:         appID,
				MchID:         mchID,
				NonceStr:      getNonceStr(),
				OutTradeNo:    req.OutTradeNo,
				ResultCode:    "SUCCESS",
				ReturnCode:    "SUCCESS",
				TotalFee:      req.TotalFee,
				TransactionID: prepayID,
				TimeEnd:       time.Now().Format("20060102150405"),
			}

			notifyParams := map[string]string{
				"return_code":    notify.ReturnCode,
				"return_msg":     notify.ReturnMsg,
				"appid":          notify.AppID,
				"mch_id":         notify.MchID,
				"nonce_str":      notify.NonceStr,
				"out_trade_no":   notify.OutTradeNo,
				"result_code":    notify.ResultCode,
				"total_fee":      fmt.Sprintf("%d", notify.TotalFee),
				"transaction_id": notify.TransactionID,
				"time_end":       notify.TimeEnd,
			}
			notify.Sign = generateSign(notifyParams, apiKey)

			xmlData, _ := xml.Marshal(notify)
			resp, err := http.Post(notifyURL, "text/xml", bytes.NewReader(xmlData))
			if err != nil {
				fmt.Println("回调失败:", err)
				return
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			fmt.Println("商户响应:", string(body))
		}()
	})
	r.POST("/pay/mock_notify", func(c *gin.Context) {
		var req struct {
			OutTradeNo string `json:"out_trade_no"`
			TotalFee   int    `json:"total_fee"`
			PrepayID   string `json:"prepay_id"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
			return
		}

		notify := PayNotifyRequest{
			AppID:         appID,
			MchID:         mchID,
			NonceStr:      getNonceStr(),
			OutTradeNo:    req.OutTradeNo,
			ResultCode:    "failure",
			ReturnCode:    "failure",
			TotalFee:      req.TotalFee,
			TransactionID: req.PrepayID,
			TimeEnd:       time.Now().Format("20060102150405"),
		}

		notifyParams := map[string]string{
			"return_code":    notify.ReturnCode,
			"return_msg":     notify.ReturnMsg,
			"appid":          notify.AppID,
			"mch_id":         notify.MchID,
			"nonce_str":      notify.NonceStr,
			"out_trade_no":   notify.OutTradeNo,
			"result_code":    notify.ResultCode,
			"total_fee":      fmt.Sprintf("%d", notify.TotalFee),
			"transaction_id": notify.TransactionID,
			"time_end":       notify.TimeEnd,
		}
		notify.Sign = generateSign(notifyParams, apiKey)

		xmlData, _ := xml.Marshal(notify)
		resp, err := http.Post(notifyURL, "text/xml", bytes.NewReader(xmlData))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "回调失败"})
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		c.JSON(http.StatusOK, gin.H{"msg": "回调成功", "response": string(body)})
	})

	// 接收支付回调接口（注意路径和真实业务统一）
	r.POST("/api/pay/payment_notify", func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		fmt.Println("收到支付回调：", string(body))
		c.String(http.StatusOK, `<xml><return_code>SUCCESS</return_code><return_msg>OK</return_msg></xml>`)
	})

	r.Run(":9090")
}
