package base

import (
	"github.com/kataras/iris/v12/context"
)

const (
	PaymentMiniProgram = "miniprogram"
	PaymentWechat      = "wechat"
	PaymentAlipay      = "alipay"
)

const (
	PaymentNotifyLen = 4096
)

const (
	// 支付成功
	PaymentStatusSuccess = "success"

	// 转入退款
	PaymentStatusRefund = "refund"

	// 未支付
	PaymentStatusNotPay = "notpay"

	// 已关闭
	PaymentStatusClosed = "closed"

	// 已撤销
	PaymentStatusRevoked = "revoked"

	// 支付中
	PaymentStatusPaying = "paying"

	// 其他错误
	PaymentStatusUnKnown = "unknown"
)

const (
	// 转账成功
	TransferStatusSuccess = "success"

	// 转账失败
	TransferStatusFailed = "failed"

	// 转账处理中
	TransferStatusProcessing = "processing"
)

const (
	// 退款成功
	RefundStatusSuccess = "success"

	// 退款处理中
	RefundStatusProcessing = "processing"

	// 退款关闭
	RefundStatusClosed = "closed"

	// 退款异常
	RefundStatusException = "exception"
)

const (
	ErrorUnknown = "ErrorUnknown"
)

type PaymentNotifyHandler func(notify *PaymentNotify)

type PaymentNotifyController struct {
	Method   string
	Endpoint string
	Handler  context.Handler
}

type PaymentEndpoint struct {
	Provider  string `yaml:"provider"`
	MchKey    string `yaml:"mch_key"`
	MchSecret string `yaml:"mch_secret"`
	CertFile  string `yaml:"cert_file"`
	KeyFile   string `yaml:"key_file"`
}

type Payment struct {
	AppKey   string `json:"app_key"`
	Type     string `json:"type"`
	TradeNo  string `json:"trade_no"`
	Desc     string `json:"desc"`
	TotalFee int    `json:"total_fee"`

	// 退款时使用，如果为0则表示全额退款
	RefundFee int `json:"refund_fee"`

	OpenID string `json:"openid"`
}

type PaymentNotify struct {
	TradeNo string
	Status  string
	Msg     string
}

type PaymentQuery struct {
	TradeNo string `json:"trade_no"`
	AppKey  string `json:"app_key"`
}

type CreatePaymentResp struct {
	Type      string `json:"type"`
	PrePayID  string `json:"prepay_id"`
	TimeStamp int64  `json:"timestamp"`
	NonceStr  string `json:"nonce_str"`
	SignType  string `json:"sign_type"`
	Sign      string `json:"sign"`
}

type PaymentTransfer struct {
	AppKey   string `json:"app_key"`
	OpenID   string `json:"openid"`
	TotalFee int    `json:"total_fee"`
	Desc     string `json:"desc"`
	TradeNo  string `json:"trade_no"`
}

type QueryTransfer struct {
	AppKey  string `json:"app_key"`
	TradeNo string `json:"trade_no"`
}

type QueryTransferResp struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

type QueryRefund struct {
	AppKey  string `json:"app_key"`
	TradeNo string `json:"trade_no"`
}

type QueryRefundResp struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

type IPaymentProvider interface {
	Init(paymentUrl string, endpoint *PaymentEndpoint) error

	Release()

	// 创建支付(下单)
	CreatePayment(payment *Payment) (*CreatePaymentResp, error)

	// 订单查询
	GetPayment(query *PaymentQuery) (*PaymentNotify, error)

	// 转账给个人
	Transfer(transfer *PaymentTransfer) error

	// 转账查询
	QueryTransfer(query *QueryTransfer) (*QueryTransferResp, error)

	// 退款
	Refund(payment *Payment) error

	// 退款查询
	QueryRefund(query *QueryRefund) (*QueryRefundResp, error)

	// 关闭订单
	Close(payment *Payment) error

	// 通知回调
	SetupNotify(handler PaymentNotifyHandler)

	GetNotifyController() *PaymentNotifyController
}

type BasePaymentProvider struct {
	IPaymentProvider
	Endpoint      *PaymentEndpoint
	PaymentUrl    string
	notifyHandler PaymentNotifyHandler
	notifyBuf     chan *PaymentNotify
	closing       chan struct{}
}

func (s *BasePaymentProvider) Init(paymentUrl string, endpoint *PaymentEndpoint) error {
	s.PaymentUrl = paymentUrl
	s.Endpoint = endpoint
	s.notifyBuf = make(chan *PaymentNotify, PaymentNotifyLen)
	s.closing = make(chan struct{}, 1)

	go s.asyncNotify()
	return nil
}

func (s *BasePaymentProvider) Release() {
	s.closing <- struct{}{}
}

func (s *BasePaymentProvider) CreatePayment(payment *Payment) (*CreatePaymentResp, error) {
	return nil, nil
}

func (s *BasePaymentProvider) GetPayment(query *PaymentQuery) (*PaymentNotify, error) {
	return nil, nil
}

func (s *BasePaymentProvider) Transfer(transfer *PaymentTransfer) error {
	return nil
}

func (s *BasePaymentProvider) QueryTransfer(query *QueryTransfer) (*QueryTransferResp, error) {
	return nil, nil
}

func (s *BasePaymentProvider) Refund(payment *Payment) error {
	return nil
}

func (s *BasePaymentProvider) QueryRefund(query *QueryRefund) (*QueryRefundResp, error) {
	return nil, nil
}

func (s *BasePaymentProvider) Close(payment *Payment) error {
	return nil
}

func (s *BasePaymentProvider) SetupNotify(handler PaymentNotifyHandler) {
	s.notifyHandler = handler
}

func (s *BasePaymentProvider) GetNotifyController() *PaymentNotifyController {
	return nil
}

func (s *BasePaymentProvider) GetReqIP() string {
	return GetIPByHost(s.PaymentUrl, "127.0.0.1")
}

func (s *BasePaymentProvider) PostNotify(notify *PaymentNotify) {
	s.notifyBuf <- notify
}

func (s *BasePaymentProvider) asyncNotify() {
	for {
		select {
		case notify := <-s.notifyBuf:
			if s.notifyHandler != nil {
				s.notifyHandler(notify)
			}

		case <-s.closing:
			return
		}
	}
}
