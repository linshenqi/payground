package payment

import "github.com/linshenqi/payground/src/base"

type Config struct {
	PaymentUrl string                          `yaml:"payment_url"`
	Endpoints  map[string]base.PaymentEndpoint `yaml:"endpoints"`
}

func (s *Config) ConfigName() string {
	return ServiceName
}

func (s *Config) Validate() error {
	return nil
}

func (s *Config) Default() interface{} {
	return &Config{}
}
