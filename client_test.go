package glittersdk_test

import (
	"testing"

	glittersdk "github.com/glitternetwork/glitter-sdk-go"
)

func Test_Client_DB(t *testing.T) {
	assert := genAssert(t)

	c := glittersdk.New()
	db := c.DB()

	// list schema
	schemas, err := db.ListSchema()
	assert(err)
	t.Log(schemas)

	// get schema
	schema, err := db.GetSchema("demo")
	assert(err)
	t.Log(schema)

	// put doc
	doc := glittersdk.Document(`{
		"doi": "10.1003/(sci)1099-1697(199803/04)7:2<65::aid-jsc357>3.0.c",
		"title": "British Steel Corporation: probably the biggest turnaround story in UK industrial history",
		"ipfs_cid": "bafybeibxvp6bawmr4u24vuza2vyretip4n7sfvivg7hdbyolxrvbodwlte"
		}`)

	txID, err := db.PutDoc("demo", doc)
	//assert(err)
	t.Log(txID)

	// get docs
	r0, err := db.GetDocs("demo", []string{"10.1003/(sci)1099-1697(199803/04)7:2<65::aid-jsc357>3.0.c"})
	assert(err)
	t.Logf("%+v\n", r0)

	// simple search
	cond1 := glittersdk.
		NewSearchCond().
		Schema("demo").
		Select("doi", "title").
		Query("British Steel").
		Page(1).
		Limit(10)
	r1, err := db.Search(cond1)
	assert(err)
	t.Logf("%+v\n", r1)

	// complex search
	cond2 := glittersdk.
		NewSearchCond().
		Schema("libgen").
		Select("doc_id", "title", "series", "author", "publisher", "language", "md5", "tags").
		Query("Springer").
		Filter(glittersdk.Filter{
			Type:     "term",
			Field:    "language",
			Value:    "English",
			From:     0.9,
			To:       1,
			DocCount: 100,
		})

	r2, err := db.Search(cond2)
	assert(err)
	t.Logf("%+v\n", r2)
}

func Test_Client_Chain(t *testing.T) {
	assert := genAssert(t)

	c := glittersdk.New()
	chain := c.Chain()

	// get chain status
	r0, err := chain.Status()
	assert(err)
	t.Logf("%+v", r0)

	// query tx
	r1, err := chain.TxSearch(`"tx.height=1"`, true, nil, nil, "")
	assert(err)
	t.Logf("%+v", r1)

	// query block
	r2, err := chain.BlockSearch(`"block.height<10"`, nil, nil, "")
	assert(err)
	t.Logf("%+v", r2)

	// block info by height
	height := int64(1)
	r3, err := chain.Block(&height)
	assert(err)
	t.Logf("%+v", r3)

	// query blockchain info
	r4, err := chain.BlockChainInfo(0, 10)
	assert(err)
	t.Logf("%+v", r4)

	// get network info
	r5, err := chain.NetInfo()
	assert(err)
	t.Logf("%+v", r5)

	// check heath status
	r6, err := chain.Health()
	assert(err)
	t.Logf("%+v", r6)
}

func genAssert(t *testing.T) func(error) {
	return func(err error) {
		if err != nil {
			t.Fatalf("err: %v", err)
			return
		}
	}
}
