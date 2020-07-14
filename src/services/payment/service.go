package payment

import (
	"errors"
	"fmt"

	"github.com/linshenqi/payground/src/services/base"
	"github.com/linshenqi/payground/src/services/wechat"
	"github.com/linshenqi/sptty"
)

const (
	ServiceName = "payment"
)

type Service struct {
	cfg       Config
	providers map[string]base.IPaymentProvider
}

func (s *Service) Init(app sptty.Sptty) error {
	if err := app.GetConfig(s.ServiceName(), &s.cfg); err != nil {
		return err
	}

	if err := s.initProviders(app); err != nil {
		return err
	}

	app.AddRoute("POST", "/v1/payments", s.postPayment)

	return nil
}

func (s *Service) Release() {
	for _, v := range s.providers {
		v.Release()
	}
}

func (s *Service) Enable() bool {
	return true
}

func (s *Service) ServiceName() string {
	return ServiceName
}

func (s *Service) initProviders(app sptty.Sptty) error {
	s.providers = map[string]base.IPaymentProvider{}

	var provider base.IPaymentProvider
	for k, v := range s.cfg.Endpoints {
		switch v.Provider {
		case base.PaymentWechat:
			provider = &wechat.PaymentProvider{}

		default:
			return fmt.Errorf("Provider Error: %s", v.Provider)
		}

		controller := provider.GetNotifyController()
		app.AddRoute(controller.Method, controller.Endpoint, controller.Handler)

		endpoint := s.cfg.Endpoints[k]
		if err := provider.Init(s.cfg.PaymentUrl, &endpoint); err != nil {
			return err
		}

		s.providers[k] = provider
	}

	return nil
}

func (s *Service) getProvider(endpoint string) (base.IPaymentProvider, error) {
	provider, exist := s.providers[endpoint]
	if !exist {
		return nil, errors.New("Provider Not Found ")
	}

	return provider, nil
}

func (s *Service) CreatePayment(endpoint string, payment *base.Payment) (*base.CreatePaymentResp, error) {
	provider, err := s.getProvider(endpoint)
	if err != nil {
		return nil, err
	}

	return provider.CreatePayment(payment)
}

func (s *Service) GetPayment(endpoint string, query *base.PaymentQuery) (*base.PaymentNotify, error) {
	provider, err := s.getProvider(endpoint)
	if err != nil {
		return nil, err
	}

	return provider.GetPayment(query)
}

func (s *Service) Transfer(endpoint string, payment *base.Payment) error {
	provider, err := s.getProvider(endpoint)
	if err != nil {
		return err
	}

	return provider.Transfer(payment)
}

func (s *Service) Refund(endpoint string, payment *base.Payment) error {
	provider, err := s.getProvider(endpoint)
	if err != nil {
		return err
	}

	return provider.Refund(payment)
}

func (s *Service) SetupNotify(endpoint string, handler base.PaymentNotifyHandler) {
	provider, err := s.getProvider(endpoint)
	if err != nil {
		return
	}

	provider.SetupNotify(handler)
}
