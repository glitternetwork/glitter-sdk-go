package testclient

import (
	"context"
	"fmt"
	chaindepindextype "github.com/glitternetwork/chain-dep/glitter_proto/glitterchain/index/types"
	"github.com/glitternetwork/glitter-sdk-go/client"
	nodeconfig "github.com/glitternetwork/glitter-sdk-go/client/node/config"
	"github.com/glitternetwork/glitter-sdk-go/client/node/remote"
	"github.com/tendermint/tendermint/libs/json"
	"testing"
)

type Desc struct {
	Description string `db:"description" json:"description" gorm:"description"`
	Github      string `db:"github" json:"github" gorm:"github"`
}

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
	resp, err := node.QueryDateset(ctx, &chaindepindextype.QueryDatesetRequest{DatasetName: "library"})
	if err != nil {
		t.Log(err)
		return
	}

	t.Log(resp)

	dataset := resp.Dateset

	desc := Desc{}
	err = json.Unmarshal([]byte(dataset.Description), &desc)
	fmt.Println(err)
	fmt.Println(desc)

}
