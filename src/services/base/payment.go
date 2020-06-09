package base

const (
	PaymentWechat = "wechat"
	PaymentAlipay = "alipay"
)

type PaymentEndpoint struct {
	Provider  string `yaml:"provider"`
	AppKey    string `yaml:"app_key"`
	AppSecret string `yaml:"app_secret"`
}

type Payment struct {
}

type IPaymentProvider interface {
	Init() error

	// 创建支付(下单)
	CreatePayment(payment *Payment) error

	// 订单查询
	GetPayment() (*Payment, error)

	// 转账给个人
	Transfer(payment *Payment) error

	// 退款
	Refund(payment *Payment) error
}

type BasePaymentProvider struct {
	IPaymentProvider
}

func (s *BasePaymentProvider) Init() error {
	return nil
}

func (s *BasePaymentProvider) CreatePayment(payment *Payment) error {
	return nil
}

func (s *BasePaymentProvider) GetPayment() (*Payment, error) {
	return nil
}

func (s *BasePaymentProvider) Transfer(payment *Payment) error {
	return nil
}

func (s *BasePaymentProvider) Refund(payment *Payment) error {
	return nil
}
