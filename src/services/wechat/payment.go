package wechat

import (
	"github.com/linshenqi/sptty"
	"gopkg.in/resty.v1"
)

type Payment struct {
	BasePaymentProvider

	http *resty.Client
}

func (s *Payment) Init() error {
	s.http = sptty.CreateHttpClient(sptty.DefaultHttpClientConfig())

	return nil
}

func (s *Payment) CreatePayment(payment *Payment) error {
	return nil
}

func (s *Payment) GetPayment() (*Payment, error) {
	return nil
}

func (s *Payment) Transfer(payment *Payment) error {
	return nil
}

func (s *Payment) Refund(payment *Payment) error {
	return nil
}
