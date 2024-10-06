package client

import (
	"context"
	"testing"
)

func Test_QuerySql(t *testing.T) {
	cli := LCDClient{
		URL: "https://orlando-api.glitterprotocol.tech",
	}

	ctx := context.Background()
	datasetName := "vec"
	sql := "SELECT md5,title,VECTOR_L2_DISTANCE(vector,TEXT_TO_VEC(?)) AS distance FROM vec.ebook ORDER BY distance LIMIT 100"
	arg := []Argument{
		{
			Type:  "STRING",
			Value: "mathematischer vorkus",
		},
	}

	r, e := cli.QuerySql(ctx, datasetName, sql, arg)
	t.Log(e)
	t.Log(r)
}
