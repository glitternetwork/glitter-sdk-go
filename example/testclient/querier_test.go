package testclient

import (
	"context"
	"github.com/glitternetwork/glitter-sdk-go/client"
	"testing"
)

func Test_Query(t *testing.T) {
	cli := New()
	ctx := context.Background()
	datasetName := "vec"
	sql := "SELECT md5,title,VECTOR_L2_DISTANCE(vector,TEXT_TO_VEC(?)) AS distance FROM vec.ebook ORDER BY distance LIMIT 100"
	arg := []client.Argument{
		{
			Type:  "STRING",
			Value: "mathematischer vorkus",
		},
	}

	r, e := cli.QuerySql(ctx, datasetName, sql, arg)
	t.Log(e)
	t.Log(r)
}
