package glittersdk

import (
	"strings"
	"time"
)

type ClientOption interface {
	apply(o *clientOption)
}

type clientOption struct {
	addrs       []string
	accessToken string
	privateKey  string
	timeout     time.Duration
}

type clientOptionFn func(o *clientOption)

func (fn clientOptionFn) apply(o *clientOption) {
	fn(o)
}

func defaultClientOption() *clientOption {
	return &clientOption{
		addrs: []string{
			"http://sg1.testnet.glitter.link:26659",
			"http://sg2.testnet.glitter.link:26659",
			"http://sg3.testnet.glitter.link:26659",
			"http://sg4.testnet.glitter.link:26659",
			"http://sg5.testnet.glitter.link:26659",
		},
	}
}

// WithAddrs create client with access token (your public key)
func WithAddrs(address ...string) ClientOption {
	return clientOptionFn(func(o *clientOption) {
		for i := 0; i < len(address); i++ {
			address[i] = strings.TrimSuffix(address[i], "/")
		}
		o.addrs = address
	})
}

// WithAccessToken create client with access token (your public key)
func WithAccessToken(token string) ClientOption {
	return clientOptionFn(func(o *clientOption) {
		o.accessToken = token
	})
}

// WithTimeout create client with request timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return clientOptionFn(func(o *clientOption) {
		o.timeout = timeout
	})
}

// WithPrivateKey create client with your private key
// some APIs that need to be signed need to use it
func WithPrivateKey(key string) ClientOption {
	return clientOptionFn(func(o *clientOption) {
		o.privateKey = key
	})
}
