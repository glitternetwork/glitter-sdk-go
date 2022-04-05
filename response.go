package glittersdk

import (
	"encoding/json"

	"github.com/tendermint/tendermint/libs/bytes"
	ttypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
)

type response struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	TX      bytes.HexBytes `json:"tx,omitempty"`
	Data    interface{}    `json:"data,omitempty"`
}

// tendermint response
type tmResponse struct {
	JSONRPC string           `json:"jsonrpc"`
	Result  json.RawMessage  `json:"result,omitempty"`
	Err     *ttypes.RPCError `json:"error,omitempty"`
}
