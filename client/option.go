package client

import (
	"time"

	"github.com/glitternetwork/glitter-sdk-go/msg"
)

const DefaultChainEndpoint = "https://orlando-api.glitterprotocol.tech"

type Option interface {
	apply(o *clientOptions)
}

// WithTimeout create client with http timeout
func WithTimeout(duration time.Duration) Option {
	return fnOption(func(o *clientOptions) {
		o.httpTimeout = duration
	})
}

// WithChainEndpoint create client with custom chain endpoint
func WithChainEndpoint(endpoint string) Option {
	return fnOption(func(o *clientOptions) {
		o.endpoint = endpoint
	})
}

// WithGasFeeConfig create client with custom gas fee config
func WithGasFeeConfig(gasPrice msg.DecCoin, gasAdjustment msg.Dec) Option {
	return fnOption(func(o *clientOptions) {
		o.gasPrice = gasPrice
		o.gasAdjustment = gasAdjustment
	})
}

type fnOption func(o *clientOptions)

func (f fnOption) apply(o *clientOptions) {
	f(o)
}

type clientOptions struct {
	endpoint      string
	gasPrice      msg.DecCoin
	gasAdjustment msg.Dec
	httpTimeout   time.Duration
}

var defaultClientOptions = clientOptions{
	endpoint:      DefaultChainEndpoint,
	gasPrice:      msg.NewDecCoinFromDec("agli", mustParseDecFromStr("1")),
	gasAdjustment: mustParseDecFromStr("2.5"),
	httpTimeout:   time.Second * 10,
}

func mustParseDecFromStr(s string) msg.Dec {
	dec, err := msg.NewDecFromStr(s)
	if err != nil {
		panic(err)
	}
	return dec
}
