package wechat

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

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

func (s *PaymentProvider) CreatePayment(payment *base.Payment) (string, error) {
	reqOrder := ReqOrder{
		MchKey:         s.BasePaymentProvider.Endpoint.MchKey,
		NonceStr:       sptty.GenerateUID(),
		SpbillCreateIP: s.BasePaymentProvider.GetReqIP(),
		NotifyUrl:      s.getNotifyUrl(),
	}

	reqOrder.FromPayment(payment)
	reqOrder.GenerateSign()

	url := fmt.Sprintf("https://api.mch.weixin.qq.com/pay/unifiedorder")
	resp, err := s.http.R().SetBody(reqOrder).Post(url)
	if err != nil {
		return "", err
	}

	respBody := ResqOrder{}
	_ = xml.Unmarshal(resp.Body(), &respBody)
	if respBody.ResultCode != ResultSuccess || respBody.ReturnCode != ResultSuccess {
		errBody, _ := json.Marshal(respBody.RespErr)
		return "", fmt.Errorf(string(errBody))
	}

	return respBody.PrepayID, nil
}

func (s *PaymentProvider) GetPayment(query *base.PaymentQuery) (*base.PaymentNotify, error) {
	req := ReqOrderQuery{
		AppKey:   query.AppKey,
		MchKey:   s.BasePaymentProvider.Endpoint.MchKey,
		NonceStr: sptty.GenerateUID(),
		TradeNo:  query.TradeNo,
	}

	req.GenerateSign()

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
