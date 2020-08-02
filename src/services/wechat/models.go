package wechat

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"sort"
	"strings"

	"github.com/linshenqi/payground/src/services/base"
)

const (
	NotifyEndpoint = "/v1/payment-wechat"
)

const (
	TradeTypeMiniProgram = "JSAPI"
)

const (
	ResultSuccess = "SUCCESS"
	ResultFail    = "FAIL"
)

func generateSign(payload interface{}, secret string) string {
	vals := map[string]interface{}{}
	body, _ := json.Marshal(payload)
	_ = json.Unmarshal(body, &vals)

	keys := []string{}
	for k := range vals {
		if k == "sign" || k == "XMLName" {
			continue
		}

		keys = append(keys, k)
	}

	sort.Strings(keys)

	str := ""
	for _, v := range keys {
		str += fmt.Sprintf("%s=%v&", v, vals[v])
	}

	str += fmt.Sprintf("key=%s", secret)

	h := md5.New()
	h.Write([]byte(str))

	sign := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
	return sign
}

type PayloadBase struct {
	XMLName xml.Name `xml:"xml"`
}

type ReqOrder struct {
	PayloadBase

	AppKey         string `xml:"appid" json:"appid"`
	MchKey         string `xml:"mch_id" json:"mch_id"`
	NonceStr       string `xml:"nonce_str" json:"nonce_str"`
	Sign           string `xml:"sign" json:"sign"`
	Boby           string `xml:"body" json:"body"`
	TradeNo        string `xml:"out_trade_no" json:"out_trade_no"`
	TotalFee       int    `xml:"total_fee" json:"total_fee"`
	SpbillCreateIP string `xml:"spbill_create_ip" json:"spbill_create_ip"`
	NotifyUrl      string `xml:"notify_url" json:"notify_url"`
	TradeType      string `xml:"trade_type" json:"trade_type"`
	OpenID         string `xml:"openid" json:"openid"`
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

func (s *ReqOrder) GenerateSign(secret string) {
	s.Sign = generateSign(s, secret)
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

func (s *ReqOrderQuery) GenerateSign(secret string) {
	s.Sign = generateSign(s, secret)
}

type ReqTransferQuery struct {
	PayloadBase

	AppKey   string `xml:"appid" json:"appid"`
	MchKey   string `xml:"mch_id" json:"mch_id"`
	NonceStr string `xml:"nonce_str" json:"nonce_str"`
	Sign     string `xml:"sign" json:"sign"`
	TradeNo  string `xml:"partner_trade_no" json:"partner_trade_no"`
}

func (s *ReqTransferQuery) FromTransferQuery(transfer *base.QueryTransfer) {
	s.AppKey = transfer.AppKey
	s.TradeNo = transfer.TradeNo
}

func (s *ReqTransferQuery) GenerateSign(secret string) {
	s.Sign = generateSign(s, secret)
}

type RespQueryTransfer struct {
	PayloadBase

	RespReturn
	RespErr

	ResultCode string `xml:"result_code"`
	Status     string `xml:"status"`
	Reason     string `xml:"reason"`
}

type ReqTransfer struct {
	ReqTransferQuery

	OpenID    string `xml:"openid" json:"openid"`
	CheckName string `xml:"check_name" json:"check_name"`
	TotalFee  int    `xml:"amount" json:"amount"`
	Desc      string `xml:"desc" json:"desc"`
}

func (s *ReqTransfer) FromPaymentTransfer(transfer *base.PaymentTransfer) {
	s.AppKey = transfer.AppKey
	s.Desc = transfer.Desc
	s.TotalFee = transfer.TotalFee
	s.TradeNo = transfer.TradeNo
	s.OpenID = transfer.OpenID
}

func (s *ReqTransfer) GenerateSign(secret string) {
	s.Sign = generateSign(s, secret)
}

type RespTransfer struct {
	PayloadBase

	RespReturn
	RespErr

	ResultCode string `xml:"result_code"`
	TradeNo    string `xml:"partner_trade_no"`
	PaymentNo  string `xml:"payment_no"`
}
