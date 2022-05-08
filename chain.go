package glittersdk

import (
	"strconv"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type Chain struct {
	client *Client // let's use cl for client and ch for chain, be consistent within the sdk.
}

// Status of the node
func (c *Chain) Status() (*ctypes.ResultStatus, error) {
	r := new(ctypes.ResultStatus)
	req := make(map[string]string)
	err := c.client.get(urlChainStatus, req, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// TxSearch searches transactions with the given query.
//
// prove: indicating whether the return value is included in the block's transaction proof
//
// page and perPage:  limit on the number of returned results
func (c *Chain) TxSearch(query string, prove bool, page, perPage *int, orderBy string) (*ctypes.ResultTxSearch, error) {
	r := new(ctypes.ResultTxSearch)
	req := map[string]string{
		"query":    query,
		"prove":    strconv.FormatBool(prove),
		"order_by": orderBy,
	}

	if page != nil {
		req["page"] = strconv.Itoa(*page)
	}
	if perPage != nil {
		req["per_page"] = strconv.Itoa(*perPage)
	}
	err := c.client.get(urlChainTxSearch, req, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// BlockSearch search for blocks by BeginBlock and EndBlock events.
//
// query: query condition example: `"block.height<10"`
//
// more detail about query https://docs.tendermint.com/v0.35/rpc/#/Websocket/subscribe
//
// page and perPage:  limit on the number of returned results
//
// orderBy: order in which blocks are sorted ("asc" or "desc"), by height.
// if empty, default sorting will be still applied.
func (c *Chain) BlockSearch(query string,
	page, perPage *int,
	orderBy string,
) (*ctypes.ResultBlockSearch, error) {
	r := new(ctypes.ResultBlockSearch)
	req := map[string]string{
		"query":    query,
		"order_by": orderBy,
	}

	if page != nil {
		req["page"] = strconv.Itoa(*page)
	}
	if perPage != nil {
		req["per_page"] = strconv.Itoa(*perPage)
	}
	err := c.client.get(urlChainBlockSearch, req, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Block get block by height
func (c *Chain) Block(height *int64) (*ctypes.ResultBlock, error) {
	r := new(ctypes.ResultBlock)
	req := map[string]string{}

	if height != nil {
		req["height"] = strconv.FormatInt(*height, 10)
	}
	err := c.client.get(urlChainBlock, req, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Health get node health.
// returns empty result (200 OK) on success, no response - in case of an error.
func (c *Chain) Health() (*ctypes.ResultHealth, error) {
	r := new(ctypes.ResultHealth)
	req := map[string]string{}
	err := c.client.get(urlChainHealth, req, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// NetInfo get network information
func (c *Chain) NetInfo() (*ctypes.ResultNetInfo, error) {
	r := new(ctypes.ResultNetInfo)
	req := map[string]string{}
	err := c.client.get(urlChainNetInfo, req, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// BlockChainInfo get blockchain info
func (c *Chain) BlockChainInfo(minHeight, maxHeight int64) (*ctypes.ResultBlockchainInfo, error) {
	r := new(ctypes.ResultBlockchainInfo)
	req := map[string]string{}
	req["minHeight"] = strconv.FormatInt(minHeight, 10)
	req["maxHeight"] = strconv.FormatInt(maxHeight, 10)

	err := c.client.get(urlChainBlockChain, req, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
