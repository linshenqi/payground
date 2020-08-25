package wechat

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/linshenqi/payground/src/services/base"
	"github.com/linshenqi/sptty"
	"time"
)

func (s *PaymentProvider) generatePaymentResp(paymentType string, prepayID string, nonceStr string, appID string) *base.CreatePaymentResp {
	resp := base.CreatePaymentResp{
		Type:      paymentType,
		PrePayID:  prepayID,
		TimeStamp: time.Now().Unix(),
		NonceStr:  nonceStr,
		SignType:  "MD5",
	}

	signBoby := map[string]interface{}{
		"appId":     appID,
		"nonceStr":  resp.NonceStr,
		"package":   fmt.Sprintf("prepay_id=%s", resp.PrePayID),
		"signType":  resp.SignType,
		"timeStamp": fmt.Sprintf("%d", resp.TimeStamp),
	}

	resp.Sign = generateSign(signBoby, s.Endpoint.MchSecret)
	return &resp
}

func (s *PaymentProvider) loadCert() (*tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(s.Endpoint.CertFile, s.Endpoint.KeyFile)
	if err != nil {
		return nil, err
	}

	return &cert, nil
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

	s.BasePaymentProvider.PostNotify(req.ToPaymentNotify())

	body, _ := xml.Marshal(RespNotify{
		RespReturn: RespReturn{
			ReturnCode: ResultSuccess,
			ReturnMsg:  "OK",
		},
	})

	_, _ = ctx.Write(body)
}
