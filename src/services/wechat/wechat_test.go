package wechat

import (
	"fmt"
	"testing"

	"github.com/linshenqi/payground/src/services/base"
)

func getProvider() (*PaymentProvider, error) {
	provider := PaymentProvider{}
	if err := provider.Init("http://wafe.com", &base.PaymentEndpoint{
		Provider:  "wechat",
		MchKey:    "1595767231",
		MchSecret: "JFGxc4TRaYKcwiHi96AhBvxWxx9CUhsO",
		CertFile:  "D:\\self\\ashibro\\微信支付\\WXCertUtil\\cert\\apiclient_cert.pem",
		KeyFile:   "D:\\self\\ashibro\\微信支付\\WXCertUtil\\cert\\apiclient_key.pem",
	}); err != nil {
		return nil, err
	}

	return &provider, nil
}

func TestWechatTransfer(t *testing.T) {
	provider, err := getProvider()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := provider.Transfer(&base.PaymentTransfer{
		AppKey:   "",
		OpenID:   "",
		TotalFee: 1,
		Desc:     "用户提现",
		TradeNo:  "123456",
	}); err != nil {
		fmt.Println(err.Error())
		return
	}
}

func TestWechatTransferQuery(t *testing.T) {
	provider, err := getProvider()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	resp, err := provider.QueryTransfer(&base.QueryTransfer{
		AppKey:  "",
		TradeNo: "123456",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(resp)
}
