package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/glitternetwork/chain-dep/core"
	chaindepindextype "github.com/glitternetwork/chain-dep/glitter_proto/glitterchain/index/types"
	"github.com/glitternetwork/glitter-sdk-go/msg"
	"github.com/glitternetwork/glitter-sdk-go/tx"
	"github.com/glitternetwork/glitter-sdk-go/utils"
	"golang.org/x/net/context/ctxhttp"
	"io/ioutil"
)

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

type ResultSet struct {
	Row map[string]*RowValue `json:"row,omitempty"`
}

type RowValue struct {
	Value           string `json:"value,omitempty"`
	ColumnValueType string `json:"column_value_type,omitempty"`
}

type Response struct {
	Result          []*ResultSet `json:"result"`
	Code            int32        `json:"code"`
	EngineTookTimes float32      ` json:"engine_took_times,omitempty"`
	FullTookTimes   float32      ` json:"full_took_times,omitempty"`
	Msg             string       ` json:"msg,omitempty"`
	TraceID         string       ` json:"trace_id,omitempty"`
}

type EngineResponse struct {
	Result    []*ResultSet `json:"result"`
	TookTimes float32      ` json:"took_times,omitempty"`
}

type Argument struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (lcd *LCDClient) QuerySql(ctx context.Context, datasetName string, sql string, argument []Argument) (resp string, err error) {
	var engineHost = ""
	var engineParam = make(map[string]interface{})
	url := lcd.URL + "/glitterchain/index/dataset/" + datasetName
	data, err := utils.CurlGet(url)
	if err != nil {
		return "", err
	}
	c := core.MakeEncodingConfig(core.ModuleBasics)
	r := &chaindepindextype.QueryDatesetResponse{}
	err = c.Marshaler.UnmarshalJSON([]byte(data), r)
	if err != nil {
		return "", err
	}
	if r.Dateset.Hosts == "" {
		return "", errors.New("obsent host")
	}
	host := r.Dateset.Hosts
	engineHost = host + "/api/v1/simple_sql_query"
	engineParam["sql"] = sql
	engineParam["arguments"] = argument
	resp, err = utils.CurlPost(engineHost, engineParam)
	return resp, err
}
