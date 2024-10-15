package testclient

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	chaindepindextype "github.com/glitternetwork/chain-dep/glitter_proto/glitterchain/index/types"
	"github.com/glitternetwork/glitter-sdk-go/utils"
	"testing"
)

func Test_CreateDataset(t *testing.T) {
	cli := New()
	ctx := context.Background()
	datasetName := "library_test"
	workStatus := chaindepindextype.ServiceStatus(2)
	hosts := "https://anybt.glitterprotocol.xyz"
	managageAddr := "glitter178uwquz93vwkc292s5n56kwq4cc2hnxfvlujmt"
	description := "{}"
	duration := int64(600)
	resp, err := cli.CreateDataset(ctx, datasetName, workStatus, hosts, managageAddr, description, duration)
	t.Log(utils.ConvToJSON(resp))
	t.Log(err)
}

func Test_EditDataset(t *testing.T) {
	cli := New()
	ctx := context.Background()
	datasetName := "library_test"
	workStatus := chaindepindextype.ServiceStatus(2)
	hosts := "https://anybt.glitterprotocol.xyz"
	managageAddr := "glitter178uwquz93vwkc292s5n56kwq4cc2hnxfvlujmt"
	description := "{\"description\":\"Dataset containing magnet links\"}"
	resp, err := cli.EditDataset(ctx, datasetName, workStatus, hosts, managageAddr, description)
	t.Log(utils.ConvToJSON(resp))
	t.Log(err)
}

func Test_EditTable(t *testing.T) {
	cli := New()
	ctx := context.Background()
	datasetName := "library_test"
	tableName := "table_test"
	description := "{\"tableName\": \"ebook_v3\", \"engine\": \"bleve\", \"comment\": \"ebook entries\", \"columns\": [{\"name\": \"_id\", \"type\": \"text\", \"comment\": \"Unique identifier for the dataset\"}, {\"name\": \"author\", \"type\": \"text\", \"comment\": \"Author of the document\"}, {\"name\": \"extension\", \"type\": \"text\", \"comment\": \"File extension of the document\"}, {\"name\": \"filesize\", \"type\": \"number\", \"comment\": \"Size of the file\"}, {\"name\": \"ipfs_cid\", \"type\": \"text\", \"comment\": \"IPFS CID for the document\"}, {\"name\": \"issn\", \"type\": \"text\", \"comment\": \"International Standard Serial Number\"}, {\"name\": \"language\", \"type\": \"text\", \"comment\": \"Language of the document\"}, {\"name\": \"publisher\", \"type\": \"text\", \"comment\": \"Publisher of the document\"}, {\"name\": \"series\", \"type\": \"text\", \"comment\": \"Series name of the document\"}, {\"name\": \"tags\", \"type\": \"text\", \"comment\": \"Tags associated with the document\"}, {\"name\": \"title\", \"type\": \"text\", \"comment\": \"Title of the document\"}, {\"name\": \"descr\", \"type\": \"text\", \"comment\": \"Description of the document\"}, {\"name\": \"coverurl\", \"type\": \"text\", \"comment\": \"URL of the cover image\"}, {\"name\": \"year\", \"type\": \"text\", \"comment\": \"Year of publication\"}]}"
	resp, err := cli.EditTable(ctx, datasetName, tableName, description)
	t.Log(utils.ConvToJSON(resp))
	t.Log(err)
}

func Test_RenewalDataset(t *testing.T) {
	cli := New()
	ctx := context.Background()
	datasetName := "library_test"
	duration := int64(10000)
	resp, err := cli.RenewalDataset(ctx, datasetName, duration)
	t.Log(utils.ConvToJSON(resp))
	t.Log(err)
}

func Test_Plede(t *testing.T) {
	cli := New()
	ctx := context.Background()
	datasetName := "library_test"
	amount := sdk.NewInt(1)
	resp, err := cli.Pledge(ctx, datasetName, amount)
	t.Log(utils.ConvToJSON(resp))
	t.Log(err)
}

func Test_ReleasePledge(t *testing.T) {
	cli := New()
	ctx := context.Background()
	datasetName := "library_test"
	amount := sdk.NewInt(1)
	resp, err := cli.ReleasePledge(ctx, datasetName, amount)
	t.Log(utils.ConvToJSON(resp))
	t.Log(err)
}
