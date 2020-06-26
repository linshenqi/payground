package wechat

import (
	"encoding/xml"

	"github.com/linshenqi/payground/src/services/base"
)

const (
	NotifyEndpoint = "/api/v1/payment-wechat"
)

const (
	TradeTypeMiniProgram = "JSAPI"
)

const (
	ResultSuccess = "SUCCESS"
	ResultFail    = "FAIL"
)

type PayloadBase struct {
	XMLName xml.Name `xml:"xml"`
}

type ReqOrder struct {
	PayloadBase

	AppKey         string `xml:"appid"`
	MchKey         string `xml:"mch_id"`
	NonceStr       string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	Boby           string `xml:"body"`
	TradeNo        string `xml:"out_trade_no"`
	TotalFee       int    `xml:"total_fee"`
	SpbillCreateIP string `xml:"spbill_create_ip"`
	NotifyUrl      string `xml:"notify_url"`
	TradeType      string `xml:"trade_type"`
}

func (s *ReqOrder) FromPayment(payment *base.Payment) {
	s.AppKey = payment.AppKey
	s.TotalFee = payment.TotalFee
	s.TradeNo = payment.TradeNo
	s.Boby = payment.Desc

	switch payment.Type {
	case base.PaymentMiniProgram:
		s.TradeType = TradeTypeMiniProgram

	default:
		s.TradeType = TradeTypeMiniProgram
	}
}

func (s *ReqOrder) GenerateSign() {

}

type RespErr struct {
	ErrCode     string `xml:"err_code" json:"err_code"`
	ErrCodeDesc string `xml:"err_code_des" json:"err_code_des"`
}

type RespReturn struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
}

type ResqOrder struct {
	PayloadBase
	RespReturn
	RespErr
	ResultCode string `xml:"result_code"`
	PrepayID   string `xml:"prepay_id"`
}

type ReqNotify struct {
	PayloadBase
	RespReturn
	ResultCode string `xml:"result_code"`
	TradeNo    string `xml:"out_trade_no"`
}

type RespNotify struct {
	PayloadBase
	RespReturn
}

type ReqOrderQuery struct {
	PayloadBase

	AppKey   string `xml:"appid"`
	MchKey   string `xml:"mch_id"`
	TradeNo  string `xml:"out_trade_no"`
	NonceStr string `xml:"nonce_str"`
	Sign     string `xml:"sign"`
}

func (s *ReqOrderQuery) GenerateSign() {

}
