package client

import (
	"fmt"
	"github.com/glitternetwork/chain-dep/core"
	"github.com/glitternetwork/glitter-sdk-go/client/node"
	nodeconfig "github.com/glitternetwork/glitter-sdk-go/client/node/config"
	"github.com/glitternetwork/glitter-sdk-go/client/node/remote"
)

func BuildNode(cfg nodeconfig.Config) (node.Node, error) {
	encodingConfig := core.MakeEncodingConfig(core.ModuleBasics)
	switch cfg.Type {
	case nodeconfig.TypeRemote:
		return remote.NewNode(cfg.Details.(*remote.Details), encodingConfig.Marshaler)
	case nodeconfig.TypeNone:
		return nil, nil
	default:
		return nil, fmt.Errorf("invalid node type: %s", cfg.Type)
	}
}
