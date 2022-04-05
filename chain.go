package glittersdk

import (
	"strconv"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type Chain struct {
	c *Client
}

// Status of the node
func (c *Chain) Status() (*ctypes.ResultStatus, error) {
	r := new(ctypes.ResultStatus)
	req := make(map[string]string)
	err := c.c.get(urlChainStatus, req, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// TxSearch Search for transactions
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
	err := c.c.get(urlChainTxSearch, req, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// BlockSearch search for blocks by BeginBlock and EndBlock events.
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
	err := c.c.get(urlChainBlockSearch, req, r)
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
	err := c.c.get(urlChainBlock, req, r)
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
	err := c.c.get(urlChainHealth, req, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// NetInfo get network information
func (c *Chain) NetInfo() (*ctypes.ResultNetInfo, error) {
	r := new(ctypes.ResultNetInfo)
	req := map[string]string{}
	err := c.c.get(urlChainNetInfo, req, r)
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

	err := c.c.get(urlChainBlockChain, req, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
