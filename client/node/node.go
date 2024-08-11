package node

import (
	"context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	types3 "github.com/cosmos/cosmos-sdk/x/auth/types"
	types2 "github.com/cosmos/cosmos-sdk/x/slashing/types"
	types4 "github.com/cosmos/cosmos-sdk/x/staking/types"
	chaindepconsumertype "github.com/glitternetwork/chain-dep/glitter_proto/glitterchain/consumer/types"
	chaindepindextype "github.com/glitternetwork/chain-dep/glitter_proto/glitterchain/index/types"
	"github.com/glitternetwork/glitter-sdk-go/utils"
	constypes "github.com/tendermint/tendermint/consensus/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"google.golang.org/grpc"
)

const (
	LocalKeeper  = "local"
	RemoteKeeper = "remote"
)

type Source interface {
	Type() string
}

type Node interface {
	// Genesis returns the genesis state
	Genesis() (*tmctypes.ResultGenesis, error)

	// ConsensusState returns the consensus state of the chain
	ConsensusState() (*constypes.RoundStateSimple, error)

	// LatestHeight returns the latest block height on the active chain. An error
	// is returned if the query fails.
	LatestHeight() (int64, error)

	// ChainID returns the network ID
	ChainID() (string, error)

	// Validators returns all the known Tendermint validators for a given block
	// height. An error is returned if the query fails.
	Validators(height int64) (*tmctypes.ResultValidators, error)

	// Block queries for a block by height. An error is returned if the query fails.
	Block(height int64) (*tmctypes.ResultBlock, error)

	// BlockResults queries the results of a block by height. An error is returnes if the query fails
	BlockResults(height int64) (*tmctypes.ResultBlockResults, error)

	// Tx queries for a transaction from the REST client and decodes it into a sdk.Tx
	// if the transaction exists. An error is returned if the tx doesn't exist or
	// decoding fails.
	Tx(hash string) (*utils.Tx, error)

	// Txs queries for all the transactions in a block. Transactions are returned
	// in the sdk.TxResponse format which internally contains an sdk.Tx. An error is
	// returned if any query fails.
	Txs(block *tmctypes.ResultBlock) ([]*utils.Tx, error)

	// TxSearch defines a method to search for a paginated set of transactions by DeliverTx event search criteria.
	TxSearch(query string, page *int, perPage *int, orderBy string) (*tmctypes.ResultTxSearch, error)

	// SubscribeEvents subscribes to new events with the given query through the RPCConfig
	// client with the given subscriber name. A receiving only channel, context
	// cancel function and an error is returned. It is up to the caller to cancel
	// the context and handle any errors appropriately.
	SubscribeEvents(subscriber, query string) (<-chan tmctypes.ResultEvent, context.CancelFunc, error)

	// SubscribeNewBlocks subscribes to the new block event handler through the RPCConfig
	// client with the given subscriber name. An receiving only channel, context
	// cancel function and an error is returned. It is up to the caller to cancel
	// the context and handle any errors appropriately.
	SubscribeNewBlocks(subscriber string) (<-chan tmctypes.ResultEvent, context.CancelFunc, error)

	// Stop defers the node stop execution to the client.
	Stop()

	GetEvmAccount(ctx context.Context, addr string)

	GetAccount(ctx context.Context, addr string) (types3.AccountI, error)

	GetCodec(ctx context.Context) codec.Codec

	GetAllValidator(ctx context.Context) ([]types4.Validator, error)

	GetAllSigningInfo(ctx context.Context) ([]types2.ValidatorSigningInfo, error)

	GetGrpcConn(ctx context.Context) *grpc.ClientConn

	GetTokenSupply(ctx context.Context, tokenName string) (sdk.Dec, error)

	GetStakingPool(ctx context.Context) (types4.Pool, error)

	QueryDateset(ctx context.Context, request *chaindepindextype.QueryDatesetRequest) (*chaindepindextype.QueryDatesetResponse, error)

	QueryDatasetExpirations(ctx context.Context, request *chaindepindextype.QueryDatasetExpirationsRequest) (*chaindepindextype.QueryDatasetExpirationsResponse, error)

	QueryDatesets(ctx context.Context, request *chaindepindextype.QueryDatesetsRequest) (*chaindepindextype.QueryDatesetsResponse, error)

	QueryCPDT(ctx context.Context, request *chaindepindextype.QueryCPDTRequest) (*chaindepindextype.QueryCPDTResponse, error)

	QueryCPDTs(ctx context.Context, request *chaindepindextype.QueryCPDTsRequest) (*chaindepindextype.QueryCPDTsResponse, error)

	QueryConsumer(ctx context.Context, request *chaindepconsumertype.QueryConsumerRequest) (*chaindepconsumertype.QueryConsumerResponse, error)

	QueryConsumers(ctx context.Context, request *chaindepconsumertype.QueryConsumersRequest) (*chaindepconsumertype.QueryConsumersResponse, error)

	QueryReleasingCPDT(ctx context.Context, request *chaindepconsumertype.QueryReleasingCPDTRequest) (*chaindepconsumertype.QueryReleasingCPDTResponse, error)

	QueryReleasingCPDTs(ctx context.Context, request *chaindepconsumertype.QueryReleasingCPDTsRequest) (*chaindepconsumertype.QueryReleasingCPDTsResponse, error)
}
