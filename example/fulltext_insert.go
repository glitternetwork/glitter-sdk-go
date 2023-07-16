package main

import (
	"context"
	"fmt"

	"github.com/glitternetwork/glitter-sdk-go/example/testclient"
	"github.com/glitternetwork/glitter-sdk-go/example/testdata"
)

func main() {
	cli := testclient.New()
	ctx := context.TODO()
	columnValues := map[string]interface{}{
		"_id":       "7f2b6638ab9ec6bfeb5924bf8e7f17e1",
		"_tx_id":    "", // The _tx_id is filled in automatically
		"author":    "J. K. Rowling",
		"extension": "pdf",
		"filesize":  743406,
		"ipfs_cid":  "bafykbzaceah6cdfb3syzrntpuuxycsfp55rtmby4oxzli2wodajgtea3ghafg",
		"issn":      "",
		"language":  "English",
		"publisher": "",
		"series":    "",
		"tags":      "'",
		"title":     "Harry Potter and the Sorcerers Stone",
		"year":      "1999",
	}
	fmt.Println("insert one row:")
	resp, err := cli.Insert(ctx, testdata.TestDBName, testdata.TestTableNameBook, columnValues)
	fmt.Printf("response=%+v,err=%+v\n", resp, err)

	fmt.Println("insert multi rows:")
	columns := []string{"_id", "_tx_id", "author", "extension", "filesize", "ipfs_cid", "issn", "language", "publisher", "series", "tags", "title", "year"}
	rowValues := [][]interface{}{
		{
			"1532675066c4913e5d0f44b82014ca9e",
			"",
			"J. K. Rowling",
			"pdf",
			3475199,
			"bafykbzaceasltcubwipjpirdmxklcwdznq4mkdx4zrey5xradmoaif34a5bn2",
			"",
			"English",
			"",
			"Harry Potter 2",
			"",
			"Harry Potter and the Chamber of Secrets (Book 2)",
			"2000",
		},
		{
			"50740153c2bf4a5db99f8b807b4a4b60",
			"",
			"J.K. Rowling, Mary GrandPré",
			"pdf",
			4478241,
			"bafykbzaceaaxtdouipt5managw2creovvg6pscsjkyqfqtocpaqg3zsmbndtm",
			"",
			"English",
			"Scholastic",
			"Harry Potter 3",
			"",
			"Harry Potter and the Prisoner of Azkaban",
			"1999",
		},
	}
	resp, err = cli.BatchInsert(ctx, testdata.TestDBName, testdata.TestTableNameBook, columns, rowValues)
	fmt.Printf("response=%+v,err=%+v\n", resp, err)
}
