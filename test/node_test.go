package test

import (
	"context"
	"github.com/cosmos/cosmos-sdk/types/query"
	chaindepindextype "github.com/glitternetwork/chain-dep/glitter_proto/glitterchain/index/types"
	"github.com/glitternetwork/glitter-sdk-go/client"
	nodeconfig "github.com/glitternetwork/glitter-sdk-go/client/node/config"
	"github.com/glitternetwork/glitter-sdk-go/client/node/remote"
	"github.com/tendermint/tendermint/libs/json"
	"testing"
)

func Test_QueryDatesets(t *testing.T) {
	cfg := nodeconfig.Config{
		Type: nodeconfig.TypeRemote,
		Details: remote.NewDetails(
			remote.NewRPCConfig("xxx", "http://sg5.testnet.glitter.link:46657", 200),
			remote.NewGrpcConfig("http://sg5.testnet.glitter.link:49090", true)),
	}
	node, err := client.BuildNode(cfg)
	if err != nil {
		t.Log(err)
		return
	}

	ctx := context.Background()
	resp, err := node.QueryDatesets(ctx, &chaindepindextype.QueryDatesetsRequest{Pagination: &query.PageRequest{Limit: 1000}})
	if err != nil {
		t.Log(err)
		return
	}

	t.Log(resp)

	datasets := resp.Datasets
	b, e := json.Marshal(datasets)
	t.Log(e)
	t.Log(string(b))
}
