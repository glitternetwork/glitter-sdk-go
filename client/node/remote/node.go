package remote

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/types/tx"
	types3 "github.com/cosmos/cosmos-sdk/x/auth/types"
	types6 "github.com/cosmos/cosmos-sdk/x/bank/types"
	types5 "github.com/cosmos/cosmos-sdk/x/slashing/types"
	types4 "github.com/cosmos/cosmos-sdk/x/staking/types"
	types2 "github.com/evmos/ethermint/x/evm/types"
	chaindepconsumertype "github.com/glitternetwork/chain-dep/glitter_proto/glitterchain/consumer/types"
	chaindepindextype "github.com/glitternetwork/chain-dep/glitter_proto/glitterchain/index/types"
	"github.com/glitternetwork/glitter-sdk-go/client/node"
	"github.com/glitternetwork/glitter-sdk-go/utils"
	constypes "github.com/tendermint/tendermint/consensus/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
	httpclient "github.com/tendermint/tendermint/rpc/client/http"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	jsonrpcclient "github.com/tendermint/tendermint/rpc/jsonrpc/client"
	tmtypes "github.com/tendermint/tendermint/types"
	"google.golang.org/grpc"
	"net/http"
	"strings"
	"time"
)

var (
	_ node.Node = &Node{}
)

// Node implements a wrapper around both a Tendermint RPCConfig client and a
// chain SDK REST client that allows for essential data queries.
type Node struct {
	ctx             context.Context
	codec           codec.Codec
	client          *httpclient.HTTP
	txServiceClient tx.ServiceClient
	grpcConnection  *grpc.ClientConn
}

// NewNode allows to build a new Node instance
func NewNode(cfg *Details, codec codec.Codec) (*Node, error) {
	//log.Infof("....begin newNode....")
	httpClient, err := jsonrpcclient.DefaultHTTPClient(cfg.RPC.Address)
	if err != nil {
		//log.Errorf("DefaultHTTPClient err=%v", err)
		return nil, err
	}
	// Tweak the transport
	httpTransport, ok := (httpClient.Transport).(*http.Transport)

	if !ok {
		//log.Errorf("invalid HTTP Transport:httpTransport=%v", httpTransport)
		return nil, fmt.Errorf("invalid HTTP Transport: %T", httpTransport)
	}

	httpTransport.MaxConnsPerHost = cfg.RPC.MaxConnections
	rpcClient, err := httpclient.NewWithClient(cfg.RPC.Address, "/websocket", httpClient)
	if err != nil {
		//log.Errorf("httpclient.NewWithClient err=%v", httpTransport)
		return nil, err
	}
	//fmt.Println("###NewNode_D")
	//err = rpcClient.Start()
	//if err != nil {
	//fmt.Println("###NewNode33333333333333", err)
	//return nil, err
	//}
	grpcConnection, err := CreateGrpcConnection(cfg.GRPC)
	if err != nil {
		//log.Errorf("CreateGrpcConnection err=%v", httpTransport)
		return nil, err
	}
	return &Node{
		ctx:   context.Background(),
		codec: codec,

		client:          rpcClient,
		txServiceClient: tx.NewServiceClient(grpcConnection),
		grpcConnection:  grpcConnection,
	}, nil
}

// Genesis implements node.Node
func (cp *Node) Genesis() (*tmctypes.ResultGenesis, error) {
	res, err := cp.client.Genesis(cp.ctx)
	if err != nil && strings.Contains(err.Error(), "use the genesis_chunked API instead") {
		return cp.getGenesisChunked()
	}
	return res, err
}

// getGenesisChunked gets the genesis data using the chinked API instead
func (cp *Node) getGenesisChunked() (*tmctypes.ResultGenesis, error) {
	bz, err := cp.getGenesisChunksStartingFrom(0)
	if err != nil {
		return nil, err
	}

	var genDoc *tmtypes.GenesisDoc
	err = tmjson.Unmarshal(bz, &genDoc)
	if err != nil {
		return nil, err
	}

	return &tmctypes.ResultGenesis{Genesis: genDoc}, nil
}

// getGenesisChunksStartingFrom returns all the genesis chunks data starting from the chunk with the given id
func (cp *Node) getGenesisChunksStartingFrom(id uint) ([]byte, error) {
	res, err := cp.client.GenesisChunked(cp.ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error while getting genesis chunk %d out of %d", id, res.TotalChunks)
	}

	bz, err := base64.StdEncoding.DecodeString(res.Data)
	if err != nil {
		return nil, fmt.Errorf("error while decoding genesis chunk %d out of %d", id, res.TotalChunks)
	}

	if id == uint(res.TotalChunks-1) {
		return bz, nil
	}

	nextChunk, err := cp.getGenesisChunksStartingFrom(id + 1)
	if err != nil {
		return nil, err
	}

	return append(bz, nextChunk...), nil
}

// ConsensusState implements node.Node
func (cp *Node) ConsensusState() (*constypes.RoundStateSimple, error) {
	state, err := cp.client.ConsensusState(context.Background())
	if err != nil {
		return nil, err
	}

	var data constypes.RoundStateSimple
	err = tmjson.Unmarshal(state.RoundState, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// LatestHeight implements node.Node
func (cp *Node) LatestHeight() (int64, error) {
	status, err := cp.client.Status(cp.ctx)
	if err != nil {
		return -1, err
	}

	height := status.SyncInfo.LatestBlockHeight
	return height, nil
}

// ChainID implements node.Node
func (cp *Node) ChainID() (string, error) {
	status, err := cp.client.Status(cp.ctx)
	if err != nil {
		return "", err
	}

	chainID := status.NodeInfo.Network
	return chainID, err
}

// Validators implements node.Node
func (cp *Node) Validators(height int64) (*tmctypes.ResultValidators, error) {
	vals := &tmctypes.ResultValidators{
		BlockHeight: height,
	}

	page := 1
	perPage := 100 // maximum 100 entries per page
	stop := false
	for !stop {
		result, err := cp.client.Validators(cp.ctx, &height, &page, &perPage)
		if err != nil {
			return nil, err
		}
		vals.Validators = append(vals.Validators, result.Validators...)
		vals.Count += result.Count
		vals.Total = result.Total
		page++
		stop = vals.Count == vals.Total
	}

	return vals, nil
}

// Block implements node.Node
func (cp *Node) Block(height int64) (*tmctypes.ResultBlock, error) {
	return cp.client.Block(cp.ctx, &height)
}

// BlockResults implements node.Node
func (cp *Node) BlockResults(height int64) (*tmctypes.ResultBlockResults, error) {
	return cp.client.BlockResults(cp.ctx, &height)
}

// Tx implements node.Node
func (cp *Node) Tx(hash string) (*utils.Tx, error) {
	res, err := cp.txServiceClient.GetTx(context.Background(), &tx.GetTxRequest{Hash: hash})
	if err != nil {
		return nil, err
	}

	// Decode messages
	for _, msg := range res.Tx.Body.Messages {
		var stdMsg sdk.Msg
		err = cp.codec.UnpackAny(msg, &stdMsg)
		if err != nil {
			return nil, fmt.Errorf("error while unpacking message: %s", err)
		}
	}

	convTx, err := utils.NewTx(res.TxResponse, res.Tx)
	if err != nil {
		return nil, fmt.Errorf("error converting transaction: %s", err.Error())
	}

	return convTx, nil
}

// Txs implements node.Node
func (cp *Node) Txs(block *tmctypes.ResultBlock) ([]*utils.Tx, error) {
	txResponses := make([]*utils.Tx, len(block.Block.Txs))
	for i, tmTx := range block.Block.Txs {
		txResponse, err := cp.Tx(fmt.Sprintf("%X", tmTx.Hash()))
		if err != nil {
			return nil, err
		}

		txResponses[i] = txResponse
	}

	return txResponses, nil
}

func (cp *Node) TxSearch(query string, page *int, perPage *int, orderBy string) (*tmctypes.ResultTxSearch, error) {
	return cp.client.TxSearch(cp.ctx, query, false, page, perPage, orderBy)
}

func (cp *Node) SubscribeEvents(subscriber, query string) (<-chan tmctypes.ResultEvent, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	eventCh, err := cp.client.Subscribe(ctx, subscriber, query)
	return eventCh, cancel, err
}

func (cp *Node) SubscribeNewBlocks(subscriber string) (<-chan tmctypes.ResultEvent, context.CancelFunc, error) {
	return cp.SubscribeEvents(subscriber, "tm.event = 'NewBlock'")
}

func (cp *Node) Stop() {
	err := cp.client.Stop()
	if err != nil {
		panic(fmt.Errorf("error while stopping proxy: %s", err))
	}

	err = cp.grpcConnection.Close()
	if err != nil {
		panic(fmt.Errorf("error while closing gRPC connection: %s", err))
	}
}

func (cp *Node) GetEvmAccount(ctx context.Context, addr string) {
	out := new(types2.QueryAccountResponse)

	//*types.QueryAccountRequest
	var req = &types2.QueryAccountRequest{
		Address: addr,
	}
	err := cp.grpcConnection.Invoke(ctx, "/cosmos.auth.v1beta1.Query/Account", req, out)
	if err != nil {

	}

	fmt.Println(out.GetNonce())
	fmt.Println(out.GetBalance())
}

func (cp *Node) GetAccount(ctx context.Context, addr string) (types3.AccountI, error) {
	out := new(types3.QueryAccountResponse)
	var req = &types3.QueryAccountRequest{
		Address: addr,
	}
	err := cp.grpcConnection.Invoke(ctx, "/cosmos.auth.v1beta1.Query/Account", req, out)
	if err != nil {
		return nil, err
	}
	var account types3.AccountI
	err = cp.codec.UnpackAny(out.GetAccount(), &account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (cp *Node) GetCodec(ctx context.Context) codec.Codec {
	return cp.codec
}

func (cp *Node) GetAllValidator(ctx context.Context) ([]types4.Validator, error) {
	out := new(types4.QueryValidatorsResponse)
	var req = &types4.QueryValidatorsRequest{}
	req.Pagination = &query.PageRequest{Limit: 10000}

	if cp.grpcConnection == nil {
		return nil, errors.New("grpcConnect_nil")
	}
	err := cp.grpcConnection.Invoke(context.Background(), "/cosmos.staking.v1beta1.Query/Validators", req, out)
	if err != nil {
		return nil, err
	}
	return out.GetValidators(), nil
}

func (cp *Node) GetAllSigningInfo(ctx context.Context) ([]types5.ValidatorSigningInfo, error) {
	out := new(types5.QuerySigningInfosResponse)
	var in = &types5.QuerySigningInfosRequest{}
	in.Pagination = &query.PageRequest{Limit: 10000}
	client := types5.NewQueryClient(cp.grpcConnection)
	out, err := client.SigningInfos(ctx, in)
	if err != nil {
		return nil, err
	}
	return out.GetInfo(), nil
}

func (cp *Node) GetGrpcConn(ctx context.Context) *grpc.ClientConn {
	return cp.grpcConnection
}

func (cp *Node) GetTokenSupply(ctx context.Context, tokenName string) (sdk.Dec, error) {
	out := new(types6.QuerySupplyOfResponse)
	client := types6.NewQueryClient(cp.grpcConnection)
	out, err := client.SupplyOf(ctx, &types6.QuerySupplyOfRequest{Denom: tokenName})
	if err != nil {
		return sdk.NewDec(0), nil
	}
	return sdk.NewDecFromInt(out.Amount.Amount), nil
}

func (cp *Node) GetStakingPool(ctx context.Context) (types4.Pool, error) {
	out := new(types4.QueryPoolResponse)
	client := types4.NewQueryClient(cp.grpcConnection)
	out, err := client.Pool(ctx, &types4.QueryPoolRequest{})
	if err != nil {
		return types4.Pool{}, nil
	}
	return out.GetPool(), nil
}

func (cp *Node) QueryDateset(ctx context.Context, request *chaindepindextype.QueryDatesetRequest) (*chaindepindextype.QueryDatesetResponse, error) {
	c := chaindepindextype.NewQueryClient(cp.grpcConnection)
	return c.QueryDateset(ctx, request)
}

func (cp *Node) QueryDatasetExpirations(ctx context.Context, request *chaindepindextype.QueryDatasetExpirationsRequest) (*chaindepindextype.QueryDatasetExpirationsResponse, error) {
	c := chaindepindextype.NewQueryClient(cp.grpcConnection)
	return c.QueryDatasetExpirations(ctx, request)
}

func (cp *Node) QueryDatesets(ctx context.Context, request *chaindepindextype.QueryDatesetsRequest) (*chaindepindextype.QueryDatesetsResponse, error) {
	c := chaindepindextype.NewQueryClient(cp.grpcConnection)
	return c.QueryDatesets(ctx, request)
}

func (cp *Node) QueryCPDT(ctx context.Context, request *chaindepindextype.QueryCPDTRequest) (*chaindepindextype.QueryCPDTResponse, error) {
	c := chaindepindextype.NewQueryClient(cp.grpcConnection)
	return c.QueryCPDT(ctx, request)
}

func (cp *Node) QueryCPDTs(ctx context.Context, request *chaindepindextype.QueryCPDTsRequest) (*chaindepindextype.QueryCPDTsResponse, error) {
	c := chaindepindextype.NewQueryClient(cp.grpcConnection)
	return c.QueryCPDTs(ctx, request)
}

func (cp *Node) QueryConsumer(ctx context.Context, request *chaindepconsumertype.QueryConsumerRequest) (*chaindepconsumertype.QueryConsumerResponse, error) {
	c := chaindepconsumertype.NewQueryClient(cp.grpcConnection)
	return c.QueryConsumer(ctx, request)
}

func (cp *Node) QueryConsumers(ctx context.Context, request *chaindepconsumertype.QueryConsumersRequest) (*chaindepconsumertype.QueryConsumersResponse, error) {
	c := chaindepconsumertype.NewQueryClient(cp.grpcConnection)
	return c.QueryConsumers(ctx, request)
}

func (cp *Node) QueryReleasingCPDT(ctx context.Context, request *chaindepconsumertype.QueryReleasingCPDTRequest) (*chaindepconsumertype.QueryReleasingCPDTResponse, error) {
	c := chaindepconsumertype.NewQueryClient(cp.grpcConnection)
	return c.QueryReleasingCPDT(ctx, request)
}

func (cp *Node) QueryReleasingCPDTs(ctx context.Context, request *chaindepconsumertype.QueryReleasingCPDTsRequest) (*chaindepconsumertype.QueryReleasingCPDTsResponse, error) {
	c := chaindepconsumertype.NewQueryClient(cp.grpcConnection)
	return c.QueryReleasingCPDTs(ctx, request)
}

func (cp *Node) Balance(ctx context.Context, request *types6.QueryBalanceRequest) (*types6.QueryBalanceResponse, error) {
	c := types6.NewQueryClient(cp.grpcConnection)
	return c.Balance(ctx, request)
}
