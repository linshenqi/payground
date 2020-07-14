package wechat

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/linshenqi/payground/src/services/base"
	"github.com/linshenqi/sptty"
	"gopkg.in/resty.v1"
)

type PaymentProvider struct {
	base.BasePaymentProvider
	http *resty.Client
}

func (s *PaymentProvider) Init(paymentUrl string, endpoint *base.PaymentEndpoint) error {
	s.http = sptty.CreateHttpClient(sptty.DefaultHttpClientConfig())
	return s.BasePaymentProvider.Init(paymentUrl, endpoint)
}

func (s *PaymentProvider) CreatePayment(payment *base.Payment) (*base.CreatePaymentResp, error) {
	reqOrder := ReqOrder{
		MchKey:         s.BasePaymentProvider.Endpoint.MchKey,
		NonceStr:       sptty.GenerateUID(),
		SpbillCreateIP: s.BasePaymentProvider.GetReqIP(),
		NotifyUrl:      s.getNotifyUrl(),
		OpenID:         payment.OpenID,
	}

	reqOrder.FromPayment(payment)
	reqOrder.GenerateSign(s.Endpoint.MchSecret)

	url := fmt.Sprintf("https://api.mch.weixin.qq.com/pay/unifiedorder")
	body, _ := xml.Marshal(reqOrder)
	resp, err := s.http.R().SetBody(body).Post(url)
	if err != nil {
		return nil, err
	}

	respBody := ResqOrder{}
	_ = xml.Unmarshal(resp.Body(), &respBody)
	if respBody.ResultCode != ResultSuccess || respBody.ReturnCode != ResultSuccess {
		errBody, _ := json.Marshal(respBody.RespErr)
		return nil, fmt.Errorf(string(errBody))
	}

	return s.generatePaymentResp(payment.Type, respBody.PrepayID), nil
}

func (s *PaymentProvider) generatePaymentResp(paymentType string, prepayID string) *base.CreatePaymentResp {
	resp := base.CreatePaymentResp{
		Type:      paymentType,
		PrePayID:  prepayID,
		TimeStamp: time.Now().Unix(),
		NonceStr:  sptty.GenerateUID(),
		SignType:  "MD5",
	}

	signBoby := map[string]interface{}{
		"appId":     s.Endpoint.MchKey,
		"nonceStr":  resp.NonceStr,
		"package":   fmt.Sprintf("prepay_id=%s", resp.PrePayID),
		"signType":  resp.SignType,
		"timeStamp": resp.TimeStamp,
	}

	resp.Sign = generateSign(signBoby, s.Endpoint.MchSecret)
	return &resp
}

func (s *PaymentProvider) GetPayment(query *base.PaymentQuery) (*base.PaymentNotify, error) {
	req := ReqOrderQuery{
		AppKey:   query.AppKey,
		MchKey:   s.BasePaymentProvider.Endpoint.MchKey,
		NonceStr: sptty.GenerateUID(),
		TradeNo:  query.TradeNo,
	}

	req.GenerateSign(s.Endpoint.MchSecret)

	url := fmt.Sprintf("https://api.mch.weixin.qq.com/pay/orderquery")
	resp, err := s.http.R().SetBody(req).Post(url)
	if err != nil {
		return nil, err
	}

	respBody := ReqNotify{}
	_ = xml.Unmarshal(resp.Body(), &respBody)
	notify := base.PaymentNotify{
		TradeNo: query.TradeNo,
		Success: false,
	}

	if respBody.ResultCode == ResultSuccess && respBody.ReturnCode == ResultSuccess {
		notify.Success = true
	}

	return &notify, nil
}

func (s *PaymentProvider) Transfer(payment *base.Payment) error {
	return nil
}

func (s *PaymentProvider) Refund(payment *base.Payment) error {
	return nil
}

func (s *PaymentProvider) GetNotifyController() *base.PaymentNotifyController {
	return &base.PaymentNotifyController{
		Method:   "POST",
		Endpoint: NotifyEndpoint,
		Handler:  s.notifyController,
	}
}

func (s *PaymentProvider) getNotifyUrl() string {
	url := fmt.Sprintf("%s%s", s.BasePaymentProvider.PaymentUrl, NotifyEndpoint)
	return url
}

func (s *PaymentProvider) notifyController(ctx iris.Context) {
	req := ReqNotify{}
	if err := ctx.ReadXML(&req); err != nil {

		body, _ := xml.Marshal(RespNotify{
			RespReturn: RespReturn{
				ReturnCode: ResultFail,
				ReturnMsg:  "Body Format Error",
			},
		})

		_, _ = ctx.Write(body)
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	sptty.Log(sptty.DebugLevel, fmt.Sprintf("Raw Payment Notify: %+v", req))

	notify := base.PaymentNotify{
		TradeNo: req.TradeNo,
		Success: false,
	}

	if req.ResultCode == ResultSuccess && req.ReturnCode == ResultSuccess {
		notify.Success = true
	}

	s.BasePaymentProvider.PostNotify(&notify)

	body, _ := xml.Marshal(RespNotify{
		RespReturn: RespReturn{
			ReturnCode: ResultSuccess,
			ReturnMsg:  "OK",
		},
	})

	_, _ = ctx.Write(body)
}
