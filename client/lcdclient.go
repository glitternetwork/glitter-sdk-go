package client

import (
	"context"
	"net/http"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/glitternetwork/glitter-sdk-go/key"
	"github.com/glitternetwork/glitter-sdk-go/msg"
	"github.com/glitternetwork/glitter-sdk-go/tx"
	glittertypes "github.com/glitternetwork/glitter.proto/golang/glitter_proto/index/types"
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

// SignAndBroadcastTX sign and broadcast transaction
func (lcd *LCDClient) SignAndBroadcastTX(ctx context.Context, options CreateTxOptions) (*sdk.TxResponse, error) {
	builder, err := lcd.CreateAndSignTx(ctx, options)
	if err != nil {
		return nil, err
	}
	return lcd.Broadcast(ctx, builder)
}

// GetAddress get account address
func (lcd *LCDClient) GetAddress() msg.AccAddress {
	return msg.AccAddress(lcd.PrivKey.PubKey().Address())
}

// SQLExecWithOptions Execute a SQL with options
// Args:
// sql: SQL statement to execute
// args: Parameters of the SQL statement, default to None
// Returns: Transaction information of the SQL execution
func (lcd *LCDClient) SQLExecWithOptions(ctx context.Context, options CreateTxOptions, sql string, args []*glittertypes.Argument) (*sdk.TxResponse, error) {
	_msg := glittertypes.NewSQLExecRequest(lcd.GetAddress(), sql, args)
	options.Msgs = []msg.Msg{_msg}
	return lcd.SignAndBroadcastTX(ctx, options)
}

// SQLExec Execute a SQL
// Args:
// sql: SQL statement to execute
// args: Parameters of the SQL statement, default to None
// Returns: Transaction information of the SQL execution
func (lcd *LCDClient) SQLExec(ctx context.Context, sql string, args []*glittertypes.Argument) (*sdk.TxResponse, error) {
	return lcd.SQLExecWithOptions(ctx, CreateTxOptions{SignMode: tx.SignModeDirect}, sql, args)
}

const (
	GrantWriter = "writer"
	GrantReader = "reader"
	GrantOwner  = "admin"
)

func (lcd *LCDClient) sqlGrantWithOptions(ctx context.Context, options CreateTxOptions, onDatabase string, onTable string, toUID string, role string) (*sdk.TxResponse, error) {
	_msg := glittertypes.NewSQLGrantRequest(lcd.GetAddress(), onDatabase, onTable, toUID, role)
	options.Msgs = []msg.Msg{_msg}
	return lcd.SignAndBroadcastTX(ctx, options)
}

// SQLGrant Grant database or table access permission
// Args:
//   - toUID: Address to grant access to
//   - role: SQL role name
//   - onDatabase: SQL database name
//   - onTable: SQL table name, optional (Grant authorization to the table if specified, otherwise grant authorization to the database)
//
// Returns:
// Result of broadcasting grant transaction
func (lcd *LCDClient) SQLGrant(ctx context.Context, onDatabase string, onTable string, toUID string, role string) (*sdk.TxResponse, error) {
	return lcd.sqlGrantWithOptions(ctx, CreateTxOptions{}, onDatabase, onTable, toUID, role)
}

// GrantWriter (insert/update/delete) permissions on the specified table to the specified user
// Args:
//   - toUID: Address to grant access
//   - onDatabase: SQL database name
//   - onTable: SQL table name, optional
//
// Returns:
// Result of grant transaction
func (lcd *LCDClient) GrantWriter(ctx context.Context, onDatabase string, onTable string, toUID string) (*sdk.TxResponse, error) {
	return lcd.SQLGrant(ctx, onDatabase, onTable, toUID, GrantWriter)
}

// GrantReader (select) permissions on the specified table to the specified user
// Args:
//   - toUID: Address to grant access
//   - onDatabase: SQL database name
//   - onTable: SQL table name, optional
//
// Returns:
// Result of grant transaction
func (lcd *LCDClient) GrantReader(ctx context.Context, onDatabase string, onTable string, toUID string) (*sdk.TxResponse, error) {
	return lcd.SQLGrant(ctx, onDatabase, onTable, toUID, GrantReader)
}

// GrantAdmin (admin) permissions on the specified table to the specified user
// Args:
//   - toUID: Address to grant access
//   - onDatabase: SQL database name
//   - onTable: SQL table name, optional
//
// Returns:
// Result of grant transaction
func (lcd *LCDClient) GrantAdmin(ctx context.Context, onDatabase string, onTable string, toUID string) (*sdk.TxResponse, error) {
	return lcd.SQLGrant(ctx, onDatabase, onTable, toUID, GrantOwner)
}
