package glittersdk

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto/ed25519"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type Cluster struct {
	c *Client
}

// UpdateValidator update validator to glitter cluster
// any suggestion clu for cluster?
func (c *Cluster) UpdateValidator(validatorPubKey string, validatorPubKeyPower int) error {
	if validatorPubKeyPower <= 0 {
		return errors.New("power must > 0")
	}
	return c.updateValidator(validatorPubKey, validatorPubKeyPower)
}

// RemoveValidator remove validator from glitter cluster
func (c *Cluster) RemoveValidator(validatorPubKey string) error {
	return c.updateValidator(validatorPubKey, 0)
}

type updateValidatorReq struct {
	ValidatorPubKey string `json:"validator_pub_key" `
	ValidatorPower  int64  `json:"validator_power"`

	SequenceID int64  `json:"seq"`
	Signature  []byte `json:"signature"`
}

type validatorSignContent struct {
	ValidatorPubKey string `json:"validator_pub_key"`
	ValidatorPower  int64  `json:"validator_power"`
	SequenceID      int64  `json:"seq"`
}

func (s *updateValidatorReq) getBytesForSign() []byte {
	m := validatorSignContent{
		ValidatorPubKey: s.ValidatorPubKey,
		ValidatorPower:  s.ValidatorPower,
		SequenceID:      s.SequenceID,
	}
	b, _ := json.Marshal(m)
	return b
}

func (c *Cluster) updateValidator(pubKey string, power int) error {
	req := &updateValidatorReq{
		ValidatorPubKey: pubKey,
		ValidatorPower:  int64(power),
		SequenceID:      time.Now().UnixMilli(),
	}
	k, err := base64.StdEncoding.DecodeString(c.c.option.privateKey)
	if err != nil {
		return errors.Errorf("invalid private key: %v", err)
	}
	key := ed25519.PrivKey(k)
	sig, err := key.Sign(req.getBytesForSign())
	if err != nil {
		return errors.Errorf("failed to sign: %v", err)
	}
	req.Signature = sig
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
