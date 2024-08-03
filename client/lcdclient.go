package client

import (
	"context"
	"net/http"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	chaindepindextype "github.com/glitternetwork/chain-dep/glitter_proto/blockved/glitterchain/index/types"
	"github.com/glitternetwork/glitter-sdk-go/key"
	"github.com/glitternetwork/glitter-sdk-go/msg"
	"github.com/glitternetwork/glitter-sdk-go/tx"
)

// LCDClient outer interface for building & signing & broadcasting tx
type LCDClient struct {
	URL           string
	ChainID       string
	GasPrice      msg.DecCoin
	GasAdjustment msg.Dec

	PrivKey        key.PrivKey
	EncodingConfig EncodingConfig

	c *http.Client
}

func (lcd *LCDClient) GetMarshaler() codec.Codec {
	return lcd.EncodingConfig.Marshaller
}

func (lcd *LCDClient) GetTxConfig() client.TxConfig {
	return lcd.EncodingConfig.TxConfig
}

// New create new Glitter client
func New(chainID string, privateKey key.PrivKey, options ...Option) *LCDClient {
	opt := defaultClientOptions
	for _, o := range options {
		o.apply(&opt)
	}
	return &LCDClient{
		URL:            opt.endpoint,
		ChainID:        chainID,
		GasPrice:       opt.gasPrice,
		GasAdjustment:  opt.gasAdjustment,
		PrivKey:        privateKey,
		EncodingConfig: MakeEncodingConfig(ModuleBasics),
		c:              &http.Client{Timeout: opt.httpTimeout},
	}
}

// CreateTxOptions tx creation options
type CreateTxOptions struct {
	Msgs []msg.Msg
	Memo string

	// Optional parameters
	AccountNumber uint64
	Sequence      uint64
	GasLimit      uint64
	FeeAmount     msg.Coins

	SignMode      tx.SignMode
	FeeGranter    msg.AccAddress
	TimeoutHeight uint64
}

// CreateAndSignTx build and sign tx
func (lcd *LCDClient) CreateAndSignTx(ctx context.Context, options CreateTxOptions) (*tx.Builder, error) {
	txbuilder := tx.NewTxBuilder(lcd.GetTxConfig())
	txbuilder.SetFeeAmount(options.FeeAmount)
	txbuilder.SetFeeGranter(options.FeeGranter)
	txbuilder.SetGasLimit(options.GasLimit)
	txbuilder.SetMemo(options.Memo)
	txbuilder.SetTimeoutHeight(options.TimeoutHeight)
	err := txbuilder.SetMsgs(options.Msgs...)
	if err != nil {
		return &txbuilder, err
	}

	// use direct sign mode as default
	if tx.SignModeUnspecified == options.SignMode {
		options.SignMode = tx.SignModeDirect
	}

	if options.AccountNumber == 0 || options.Sequence == 0 {
		account, err := lcd.LoadAccount(ctx, msg.AccAddress(lcd.PrivKey.PubKey().Address()))
		if err != nil {
			return nil, sdkerrors.Wrap(err, "failed to load account")
		}

		options.AccountNumber = account.GetAccountNumber()
		options.Sequence = account.GetSequence()
		time.Sleep(time.Second)
	}

	gasLimit := int64(options.GasLimit)
	if options.GasLimit == 0 {
		simulateRes, err := lcd.Simulate(ctx, txbuilder, options)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "failed to simulate")
		}

		gasLimit = lcd.GasAdjustment.MulInt64(int64(simulateRes.GasInfo.GasUsed)).TruncateInt64()
		txbuilder.SetGasLimit(uint64(gasLimit))
	}

	if options.FeeAmount.IsZero() {
		gasFee := msg.NewCoin(lcd.GasPrice.Denom, lcd.GasPrice.Amount.MulInt64(gasLimit).TruncateInt())
		coins := msg.Coins{}.Add(gasFee)
		txbuilder.SetFeeAmount(coins)
	} else {
		txbuilder.SetFeeAmount(options.FeeAmount)
	}
	err = txbuilder.Sign(options.SignMode, tx.SignerData{
		AccountNumber: options.AccountNumber,
		ChainID:       lcd.ChainID,
		Sequence:      options.Sequence,
	}, lcd.PrivKey, true)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to sign tx")
	}

	return &txbuilder, nil
}

func (lcd *LCDClient) SignAndBroadcastTX(ctx context.Context, options CreateTxOptions) (*sdk.TxResponse, error) {
	builder, err := lcd.CreateAndSignTx(ctx, options)
	if err != nil {
		return nil, err
	}
	return lcd.Broadcast(ctx, builder)
}

func (lcd *LCDClient) GetAddress() msg.AccAddress {
	return msg.AccAddress(lcd.PrivKey.PubKey().Address())
}

func (lcd *LCDClient) CreateDataset(ctx context.Context, options CreateTxOptions, datasetName string, workStatus chaindepindextype.ServiceStatus, hosts string, manageAddresses string, meta string, description string, duration int64) (*sdk.TxResponse, error) {
	_msg := chaindepindextype.NewCreateDatasetRequest(lcd.GetAddress(), datasetName, workStatus, hosts, manageAddresses, meta, description, duration)
	options.Msgs = []msg.Msg{_msg}
	return lcd.SignAndBroadcastTX(ctx, options)
}
