# glitter-sdk-go

> glitter sdk for golang

## Quick Start

```go
import "github.com/glitternetwork/glitter-sdk-go"
```

```go
	// create sdk
	client := glittersdk.New()

	db:=client.DB()

	// put document
	doc := glittersdk.Document(`{
		"doi": "10.1003/(sci)1099-1697(199803/04)7:2<65::aid-jsc357>3.0.c",
		"title": "British Steel Corporation: probably the biggest turnaround story in UK industrial history",
		"ipfs_cid": "bafybeibxvp6bawmr4u24vuza2vyretip4n7sfvivg7hdbyolxrvbodwlte"
		}`)

	txID, err := db.PutDoc("demo", doc)
	checkerr(err)
	fmt.Printf("tx id=%s\n", txID)

	// search document
	cond := glittersdk.
		NewSearchCond().
		Schema("demo").
		Select("doi", "title").
		Query("British Steel").
		Page(1).
		Limit(10)
	sr, err := db.Search(cond)
	checkerr(err)
	fmt.Printf("%+v\n", sr)
```

## SDK
### Options

|Option|Description|
|----|----|
|WithAddrs(address ...string)|set glitter address|
|WithAccessToken(token string)|set glitter token|
|WithTimeout(timeout time.Duration)|set client timeout|

## API

See [doc.md](doc.md)

## More useage

See [client_test.go](client_test.go)
