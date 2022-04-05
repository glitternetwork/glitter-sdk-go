package glittersdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"

	tmjson "github.com/tendermint/tendermint/libs/json"
)

const (
	urlGetDocs          = "/v1/get_docs"
	urlSearch           = "/v1/search"
	urlPutDoc           = "/v1/put_doc"
	urlListSchema       = "/v1/list_schema"
	urlCreateSchema     = "/v1/create_schema"
	urlUpdateValidator  = "/v1/admin/update_validator"
	urlChainValidator   = "/v1/chain/validators"
	urlChainStatus      = "/v1/chain/status"
	urlChainTxSearch    = "/v1/chain/tx_search"
	urlChainBlockSearch = "/v1/chain/block_search"
	urlChainBlock       = "/v1/chain/block"
	urlChainNetInfo     = "/v1/chain/net_info"
	urlChainBlockChain  = "/v1/chain/blockchain"
	urlChainHealth      = "/v1/chain/health"
)

// Client is a HTTP API client to the glitter service.
type Client struct {
	option  *clientOption
	db      *Database
	cluster *Cluster
	chain   *Chain

	client    *http.Client
	addrIndex uint32
}

// New create a new glitter client
func New(opts ...ClientOption) *Client {
	opt := defaultClientOption()
	for _, o := range opts {
		o.apply(opt)
	}
	client := &http.Client{
		Timeout: opt.timeout,
	}
	c := &Client{option: opt, client: client}
	c.db = &Database{c: c}
	c.cluster = &Cluster{c: c}
	c.chain = &Chain{c: c}
	return c
}

// DB provide database api to put/search document or manager schema
func (c *Client) DB() *Database {
	return c.db
}

// Cluster provide cluster api to manager validator or nodes
func (c *Client) Cluster() *Cluster {
	return c.cluster
}

// Chain provide chain api to search tx and block info
func (c *Client) Chain() *Chain {
	return c.chain
}

const tmPrefix = "/v1/chain"

func unmarshalTMResp(body io.Reader, v interface{}) error {
	r := &tmResponse{}
	err := unmarshalJSON(body, r)
	if err != nil {
		return err
	}
	if r.Err != nil {
		return errors.New(r.Err.Message)
	}
	return tmjson.Unmarshal(r.Result, v)
}

func unmarshalGlitterResp(body io.Reader, v interface{}) error {
	var r *response
	if rr, ok := v.(*response); ok {
		r = rr
	} else {
		r = &response{Data: v}
	}
	err := unmarshalJSON(body, r)
	if err != nil {
		return err
	}
	if r.Code != 0 {
		return errors.New(r.Message)
	}
	return nil
}

func unmarshalJSON(body io.Reader, v interface{}) error {
	return json.NewDecoder(body).Decode(v)
}

func (c *Client) post(path string, request, dest interface{}) error {
	data, err := json.Marshal(request)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, c.joinURL(path), bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("access_token", c.option.accessToken)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if strings.HasPrefix(path, tmPrefix) {
		return unmarshalTMResp(resp.Body, dest)
	}
	return unmarshalGlitterResp(resp.Body, dest)
}

func (c *Client) get(path string, request map[string]string, dest interface{}) error {
	uv := url.Values{}
	for k, v := range request {
		uv.Set(k, v)
	}
	url := fmt.Sprintf("%s?%s", c.joinURL(path), uv.Encode())
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("access_token", c.option.accessToken)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if strings.HasPrefix(path, tmPrefix) {
		return unmarshalTMResp(resp.Body, dest)
	}
	return unmarshalGlitterResp(resp.Body, dest)
}

func (c *Client) joinURL(path string) string {
	return c.selectAddr() + path
}

func (c *Client) selectAddr() string {
	addrs := c.option.addrs
	if len(addrs) == 1 {
		return addrs[0]
	}
	i := int(atomic.AddUint32(&c.addrIndex, 1))
	if i < 0 {
		i = -i
	}
	i = i % len(c.option.addrs)
	return c.option.addrs[i]
}
