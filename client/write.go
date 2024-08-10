package client

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	chaindepconsumertype "github.com/glitternetwork/chain-dep/glitter_proto/glitterchain/consumer/types"
	chaindepindextype "github.com/glitternetwork/chain-dep/glitter_proto/glitterchain/index/types"
	"github.com/glitternetwork/glitter-sdk-go/msg"
	"github.com/glitternetwork/glitter-sdk-go/tx"
)

func (lcd *LCDClient) CreateDataset(ctx context.Context, datasetName string, workStatus chaindepindextype.ServiceStatus, hosts string, manageAddresses string, meta string, description string, duration int64) (*sdk.TxResponse, error) {
	_msg := chaindepindextype.NewCreateDatasetRequest(lcd.GetAddress(), datasetName, workStatus, hosts, manageAddresses, meta, description, duration)
	options := CreateTxOptions{SignMode: tx.SignModeDirect}
	options.Msgs = []msg.Msg{_msg}
	return lcd.SignAndBroadcastTX(ctx, options)
}

func (lcd *LCDClient) EditDatasetRequest(ctx context.Context, datasetName string, workStatus chaindepindextype.ServiceStatus, hosts string, manageAddresses string, meta string, description string) (*sdk.TxResponse, error) {
	_msg := chaindepindextype.NewEditDatasetRequest(lcd.GetAddress(), datasetName, workStatus, hosts, manageAddresses, meta, description)
	options := CreateTxOptions{SignMode: tx.SignModeDirect}
	options.Msgs = []msg.Msg{_msg}
	return lcd.SignAndBroadcastTX(ctx, options)
}

func (lcd *LCDClient) RenewalDataset(ctx context.Context, datasetName string, duration int64) (*sdk.TxResponse, error) {
	_msg := chaindepindextype.NewRenewalDatasetRequest(lcd.GetAddress(), datasetName, duration)
	options := CreateTxOptions{SignMode: tx.SignModeDirect}
	options.Msgs = []msg.Msg{_msg}
	return lcd.SignAndBroadcastTX(ctx, options)
}

func (lcd *LCDClient) Pledge(ctx context.Context, datasetName string, amount sdk.Int) (*sdk.TxResponse, error) {
	_msg := chaindepconsumertype.NewPledgeRequest(lcd.GetAddress(), datasetName, amount)
	options := CreateTxOptions{SignMode: tx.SignModeDirect}
	options.Msgs = []msg.Msg{_msg}
	return lcd.SignAndBroadcastTX(ctx, options)
}

func (lcd *LCDClient) ReleasePledge(ctx context.Context, datasetName string, amount sdk.Int) (*sdk.TxResponse, error) {
	_msg := chaindepconsumertype.NewReleasePledgeRequest(lcd.GetAddress(), datasetName, amount)
	options := CreateTxOptions{SignMode: tx.SignModeDirect}
	options.Msgs = []msg.Msg{_msg}
	return lcd.SignAndBroadcastTX(ctx, options)
}
