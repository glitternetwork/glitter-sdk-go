package main

import (
	"context"
	"fmt"
	glittercommon "github.com/glitternetwork/chain-dep/glitter_proto/common"

	"github.com/glitternetwork/glitter-sdk-go/example/testclient"
	"github.com/glitternetwork/glitter-sdk-go/example/testdata"
	"github.com/glitternetwork/glitter-sdk-go/utils"
)

func main() {
	cli := testclient.New()
	ctx := context.TODO()

	type Book struct {
		ID        string `db:"_id"`
		TxID      string `db:"_tx_id"`
		Author    string `db:"author"`
		Extension string `db:"extension"`
		Filesize  int    `db:"filesize"`
		IPFSCID   string `db:"ipfs_cid"`
		ISSN      string `db:"issn"`
		Language  string `db:"language"`
		Publisher string `db:"publisher"`
		Series    string `db:"series"`
		Tags      string `db:"tags"`
		Title     string `db:"title"`
		Year      string `db:"year"`
	}

	fmt.Println("=====query all:")
	var books []Book
	err := cli.QueryScan(ctx, &books,
		fmt.Sprintf("select * from %s.%s limit 10", testdata.TestDBName, testdata.TestTableNameBook))
	fmt.Printf("books=%+v,err=%+v\n", books, err)

	// full text search
	fmt.Println("=====match query:")
	title := "Harry Potter"
	author := "J.K. Rowling"

	qb := utils.NewQueryString()
	qb.AddMatchQuery("title", title, 1)
	qb.AddMatchQuery("author", author, 0.5)
	hint := utils.HighlightHint([]string{"author", "title"})

	sql := fmt.Sprintf("select %s _score,* from %s.%s where  query_string(?) limit 0,10", hint, testdata.TestDBName, testdata.TestTableNameBook)
	arg := &glittercommon.Argument{
		Type:  glittercommon.Argument_STRING,
		Value: qb.GetQueryString(),
	}
	resp, err := cli.Query(ctx, sql, arg)
	fmt.Printf("resp=%+v,err=%+v\n", resp, err)

	fmt.Println("=====match phrase query:")
	title = "Harry Potter"
	qb = utils.NewQueryString()
	qb.AddMatchPhraseQuery("title", title, 1)
	hint = utils.HighlightHint([]string{"title"})

	sql = fmt.Sprintf("select %s _score,* from %s.%s where  query_string(?) limit 0,10", hint, testdata.TestDBName, testdata.TestTableNameBook)
	arg = &glittercommon.Argument{
		Type:  glittercommon.Argument_STRING,
		Value: qb.GetQueryString(),
	}
	resp, err = cli.Query(ctx, sql, arg)
	fmt.Printf("resp=%+v,err=%+v\n", resp, err)

}
