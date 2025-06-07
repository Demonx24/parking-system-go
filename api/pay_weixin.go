package api

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"parking-system-go/global"
	"parking-system-go/model/database"
	"parking-system-go/model/request"
	"parking-system-go/model/response"
	"parking-system-go/utils"
	"strconv"
	"strings"
	"time"
)

type PayWeChatApi struct{}

func getNonceStr() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// 用于模拟调用微信“统一下单”接口的流程，生成二维码支付链接
func (p *PayWeChatApi) CreatePayment(c *gin.Context) {
	WeChat := global.Config.WeChat
	req := database.Order{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	// 金额转为分（微信支付单位）
	totalFee := int(req.Amount * 100)
	if err := payService.CreateMockOrder(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	nonceStr := getNonceStr()

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
	// 构建请求参数
	params := map[string]string{
		"appid":            WeChat.WeChatAppID,
		"mch_id":           WeChat.WeChatMchID,
		"nonce_str":        nonceStr,
		"body":             "停车场支付订单:" + req.OrderID,
		"out_trade_no":     req.OrderID,
		"total_fee":        fmt.Sprintf("%d", totalFee),
		"spbill_create_ip": c.ClientIP(),
		"notify_url":       WeChat.WeChatNotifyURL,
		"trade_type":       "NATIVE",
	}

	// 签名
	sign := utils.GenerateSign(params, WeChat.WeChatAPIKey)
	fmt.Println(sign)
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
		response.FailWithMessage("请求微信支付失败", c)
		return
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	var wxResp response.UnifiedOrderResponse
	if err := xml.Unmarshal(respBody, &wxResp); err != nil {
		response.FailWithMessage("解析微信响应失败", c)
		return
	}
	req.PrepayID = wxResp.PrepayID
	if err := orderService.Update(&req); err != nil {
		response.FailWithMessage("更新微信支付订单失败", c)
		return
	}
	response.OkWithData(gin.H{
		"code_url":  wxResp.CodeURL,
		"prepay_id": wxResp.PrepayID,
	}, c)
}

// 支付回调接口，微信支付完成后调用

func (p *PayWeChatApi) PaymentNotify(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "fail")
		return
	}

	// 先把 XML 解析成 map[string]string
	params := make(map[string]string)
	if err := xml.Unmarshal(body, (*MapXML)(&params)); err != nil {
		c.String(http.StatusBadRequest, "fail")
		return
	}

	// 打印所有回调参数，方便调试
	fmt.Println("微信回调参数:", params)

	// 必须先检查 return_code 是否为 SUCCESS
	if params["return_code"] != "SUCCESS" {
		c.String(http.StatusOK, "<xml><return_code>FAIL</return_code><return_msg>支付失败</return_msg></xml>")
		return
	}

	// 生成签名校验，去掉 sign 字段
	sign := utils.GenerateSign(params, global.WeChat.WeChatAPIKey)
	fmt.Println("服务器计算签名:", sign)
	fmt.Println("微信回调原始签名:", params["sign"])

	if sign != params["sign"] {
		c.String(http.StatusOK, "<xml><return_code>FAIL</return_code><return_msg>签名失败</return_msg></xml>")
		return
	}

	// 判断支付结果
	if params["result_code"] == "SUCCESS" {
		// TODO: 处理订单支付成功逻辑

		// 保存更新的停车记录
		order, err := orderService.GetOrder(database.Order{OrderID: params["out_trade_no"]})
		if err != nil {
			// 没找到订单，可能是恶意请求
			global.Log.Error("订单不存在", zap.Error(err))
			response.FailWithMessage(err.Error(), c)
			return
		}
		if order.Status == 0 {
			// 修改订单为已支付
			order.Status = 1
			payService.MockPayment(order.OrderID)

			// 查回原始的 parking record
			parkingRecord, err := parkingrecordService.GetRecord(database.ParkingRecord{ParkingLotID: order.ParkingRecordID})
			if err == nil {
				parkingRecord.PlateNumber = order.PlateNumber
				parkingRecord.Status = 1
				now := time.Now()
				parkingRecord.ExitTime = &now
				parkingRecord.TotalFee = order.Amount
				parkingrecordService.Update(&parkingRecord)

				// 释放车位等其他逻辑
				parkinglotService.IncrementAvailableSlotsWithPessimisticLock(parkingRecord.ParkingLotID)
			}

			// 返回微信处理成功
			c.String(http.StatusOK, "<xml><return_code>SUCCESS</return_code><return_msg>OK</return_msg></xml>")
			return
		}
	}

	c.String(http.StatusOK, "<xml><return_code>FAIL</return_code><return_msg>支付失败</return_msg></xml>")
}

// 签名和生成 XML 回调的逻辑放在后端（Gin）中，由前端触发一个请求，让后端来构造完整的 XML 并发送回调。
func (p *PayWeChatApi) MockCallback(c *gin.Context) {
	var order = database.Order{}
	if err := c.ShouldBindQuery(&order); err != nil {
		response.FailWithMessage("订单号参数错误", c)
		return
	}
	if order.OrderID == "" {
		c.String(http.StatusBadRequest, "缺少订单号")
		return
	}

	// 查询数据库获取订单信息（假设你已经有方法）
	order, err := orderService.GetOrder(order)
	if err != nil {
		c.String(http.StatusNotFound, "订单不存在")
		return
	}

	params := map[string]string{
		"return_code":    "SUCCESS",
		"appid":          "wx1234567890abcdef",
		"mch_id":         "1234567890",
		"nonce_str":      "mock123456789",
		"result_code":    "SUCCESS",
		"out_trade_no":   order.OrderID,
		"total_fee":      strconv.Itoa(int(order.Amount)), // 单位: 分
		"transaction_id": order.PrepayID,
		"time_end":       time.Now().Format("20060102150405"),
	}

	params["sign"] = utils.GenerateSign(params, "wlc1224") // 微信签名封装方法

	xmlData := utils.MapToXML(params)

	// 回调自己接口
	resp, err := http.Post("http://localhost:8081/api/pay/payment_notify", "application/xml", strings.NewReader(xmlData))
	if err != nil {
		c.String(http.StatusInternalServerError, "模拟回调失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.String(http.StatusOK, string(body))
}

// MapXML 是辅助结构，用来解析 XML 到 map[string]string
type MapXML map[string]string

func (m *MapXML) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = MapXML{}
	for {
		t, err := d.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		switch token := t.(type) {
		case xml.StartElement:
			var v string
			if err := d.DecodeElement(&v, &token); err != nil {
				return err
			}
			(*m)[token.Name.Local] = v
		case xml.EndElement:
			if token == start.End() {
				return nil
			}
		}
	}
	return nil
}
