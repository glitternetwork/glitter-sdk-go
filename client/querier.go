package client

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"io/ioutil"

	"github.com/glitternetwork/glitter-sdk-go/msg"
	"github.com/glitternetwork/glitter-sdk-go/tx"
	"golang.org/x/net/context/ctxhttp"
)

// QueryAccountResData response
type QueryAccountResData struct {
	Address       msg.AccAddress `json:"address"`
	AccountNumber msg.Int        `json:"account_number"`
	Sequence      msg.Int        `json:"sequence"`
}

// QueryAccountRes response
type QueryAccountRes struct {
	Account QueryAccountResData `json:"account"`
}

// LoadAccount simulates gas and fee for a transaction
func (lcd *LCDClient) LoadAccount(ctx context.Context, address msg.AccAddress) (res authtypes.AccountI, err error) {
	resp, err := ctxhttp.Get(ctx, lcd.c, lcd.URL+fmt.Sprintf("/cosmos/auth/v1beta1/accounts/%s", address))
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to estimate")
	}
	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to read response")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 response code %d: %s", resp.StatusCode, string(out))
	}

	var response authtypes.QueryAccountResponse
	//err = json.Unmarshal(out, &response)
	err = lcd.GetMarshaler().UnmarshalJSON(out, &response)
	//err = lcd.EncodingConfig.Marshaller.UnmarshalJSON(out, &response)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to unmarshal response")
	}

	//fmt.Printf("LoadAccount address=%s resp=%+v\n", address.String(), response)
	return response.Account.GetCachedValue().(authtypes.AccountI), nil
}

// Simulate tx and get response
func (lcd *LCDClient) Simulate(ctx context.Context, txbuilder tx.Builder, options CreateTxOptions) (*sdktx.SimulateResponse, error) {
	// Create an empty signature literal as the ante handler will populate with a
	// sentinel pubkey.
	sig := signing.SignatureV2{
		PubKey: &secp256k1.PubKey{},
		Data: &signing.SingleSignatureData{
			SignMode: options.SignMode,
		},
		Sequence: options.Sequence,
	}
	if err := txbuilder.SetSignatures(sig); err != nil {
		return nil, err
	}

	bz, err := txbuilder.GetTxBytes()
	if err != nil {
		return nil, err
	}
	reqBytes, err := lcd.GetMarshaler().MarshalJSON(&sdktx.SimulateRequest{
		Tx:      nil,
		TxBytes: bz,
	})
	if err != nil {
		return nil, err
	}

	resp, err := ctxhttp.Post(ctx, lcd.c, lcd.URL+"/cosmos/tx/v1beta1/simulate", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to estimate")
	}
	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to read response")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 response code %d: %s", resp.StatusCode, string(out))
	}

	var response sdktx.SimulateResponse
	err = lcd.GetMarshaler().UnmarshalJSON(out, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
