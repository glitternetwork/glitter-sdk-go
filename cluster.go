package glittersdk

import (
	"errors"
	"strconv"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type Cluster struct {
	c *Client
}

// UpdateValidator update validator to glitter cluster
func (c *Cluster) UpdateValidator(pubKey string, power int) error {
	if power <= 0 {
		return errors.New("power must > 0")
	}
	return c.updateValidator(pubKey, power)
}

// RemoveValidator remove validator from glitter cluster
func (c *Cluster) RemoveValidator(pubKey string) error {
	return c.updateValidator(pubKey, 0)
}

type updateValidatorReq struct {
	PubKey string `json:"pub_key"`
	Power  int    `json:"power"`
}

func (c *Cluster) updateValidator(pubKey string, power int) error {
	req := &updateValidatorReq{PubKey: pubKey, Power: power}
	var resp []byte
	return c.c.post(urlUpdateValidator, req, &resp)
}

// ListValidators list validators info
func (c *Cluster) ListValidators(height *int64, page, perPage *int) (*ctypes.ResultValidators, error) {
	r := new(ctypes.ResultValidators)
	req := make(map[string]string)
	if page != nil {
		req["page"] = strconv.Itoa(*page)
	}
	if perPage != nil {
		req["per_page"] = strconv.Itoa(*perPage)
	}
	if height != nil {
		req["height"] = strconv.FormatInt(*height, 10)
	}
	err := c.c.get(urlChainValidator, req, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
