package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/glitternetwork/glitter-sdk-go/client"
	"github.com/glitternetwork/glitter-sdk-go/example/testclient"
	"github.com/glitternetwork/glitter-sdk-go/example/testdata"
)

func main() {
	cli := testclient.New()
	ctx := context.TODO()
	createDatabase(ctx, cli, testdata.TestDBName)
	createFulltextEngineTable(ctx, cli, testdata.TestDBName, testdata.TestTableNameBook)
	createStandardEngineTable(ctx, cli, testdata.TestDBName, testdata.TestTableNameUser)
	listDBTables(ctx, cli, testdata.TestDBName)

}

func createDatabase(ctx context.Context, cli *client.LCDClient, db string) {
	r, err := cli.CreateDatabase(ctx, db)
	if err != nil {
		panic(errors.Wrap(err, "failed to create database"))
	}
	fmt.Printf("create database result: %+v", r)
	time.Sleep(time.Second * 10)
}

func createFulltextEngineTable(ctx context.Context, cli *client.LCDClient, db, table string) {
	ddlTpl := `CREATE TABLE IF NOT EXISTS %s.%s (
        _id VARCHAR(255) PRIMARY KEY COMMENT 'md5',
        title VARCHAR(2000) COMMENT 'title',
        series VARCHAR(512) COMMENT 'series',
        author VARCHAR(512) COMMENT 'author',
        publisher VARCHAR(512) COMMENT 'publisher',
        language VARCHAR(128) COMMENT 'language',
        tags VARCHAR(512) COMMENT 'tags',
        issn VARCHAR(32) COMMENT 'issn',
        ipfs_cid VARCHAR(512) COMMENT 'ipfs cid',
        extension VARCHAR(512) COMMENT 'extension',
        year VARCHAR(14) COMMENT 'year',
        filesize INT(11),
        _tx_id VARCHAR(255) COMMENT 'transaction id auto generate',
        FULLTEXT INDEX(title) WITH PARSER standard,
        FULLTEXT INDEX(series) WITH PARSER keyword,
        FULLTEXT INDEX(author) WITH PARSER standard,
        FULLTEXT INDEX(publisher) WITH PARSER standard,
        FULLTEXT INDEX(language) WITH PARSER standard,
        FULLTEXT INDEX(tags) WITH PARSER standard,
        FULLTEXT INDEX(ipfs_cid) WITH PARSER keyword,
        FULLTEXT INDEX(extension) WITH PARSER keyword,
        FULLTEXT INDEX(year) WITH PARSER keyword
    ) ENGINE = full_text COMMENT 'book records'`
	ddl := fmt.Sprintf(ddlTpl, db, table)
	r, err := cli.CreateTable(ctx, ddl)
	if err != nil {
		panic(errors.Wrap(err, "failed to create fulltext engine table"))
	}
	fmt.Printf("create table result: %+v", r)
	time.Sleep(time.Second * 10)
}

func createStandardEngineTable(ctx context.Context, cli *client.LCDClient, db, table string) {
	ddlTpl := `CREATE TABLE  IF NOT EXISTS %s.%s (
     _id VARCHAR(500) PRIMARY KEY COMMENT 'document id',
     author VARCHAR(255) NOT NULL Default '' COMMENT 'ens address or lens address',
     handle VARCHAR(128) NOT NULL Default '' COMMENT 'ens or lens handler', 
     display_name VARCHAR(128) NOT NULL Default '' COMMENT 'nickname',
     avatar_url VARCHAR(255) NOT NULL Default '' COMMENT 'the url of avatar',
     entry_num int(11) NOT NULL Default 0 COMMENT 'the article numbers' ,
     status int(11) NOT NULL Default 0 ,         
     source VARCHAR(64) NOT NULL Default '' COMMENT 'enum: mirror, lens, eip1577',
     domain VARCHAR(128) NOT NULL Default '' COMMENT 'mirror second domain',
     _tx_id VARCHAR(255) COMMENT 'transaction id auto generate',
     KEY ` + "`author_idx` (`author`)," + `
     KEY ` + "`handle_idx` (`handle`)," + `
     KEY ` + "`display_name_idx` (`display_name`)," + `
     KEY ` + "`domain_idx` (`domain`)" + `
     ) ENGINE=standard COMMENT 'all user info:mirror,lens,eip1577 and so on'`
	ddl := fmt.Sprintf(ddlTpl, db, table)
	r, err := cli.CreateTable(ctx, ddl)
	if err != nil {
		panic(errors.Wrap(err, "failed to create standard engine table"))
	}
	fmt.Printf("create table result: %+v", r)
	time.Sleep(time.Second * 10)
}

func listDBTables(ctx context.Context, cli *client.LCDClient, db string) {
	r, err := cli.ListTables(ctx, "", "", db, nil, nil)
	if err != nil {
		panic(errors.Wrap(err, "failed to list db tables"))
	}
	fmt.Printf("%s\t%s\n", "Table", "Creator")
	for _, v := range r.Tables {
		fmt.Printf("%s\t%s\n", v.TableName, v.Creator)
	}
}
