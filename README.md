# glitter-sdk-go

> glitter sdk for golang

## Quick Start

```go
	import "github.com/glitternetwork/glitter-sdk-go"
```

```go
	// get a client
	client := glittersdk.New()
	// list schema
	schema := client.DB().GetSchema("sample")
	// put document
	txID,err := client.DB().PutDoc("sample",glittersdk.Document(`{
		"url"  : "https://glitterprotocol.io/",
		"title": "A Decentralized Content Indexing Network"
	}`))
	// get doc by primary key
	docs, err := db.GetDocs("sample", []string{"https://glitterprotocol.io/"})
	// search document
	searchRes, err := db.Search(
		glittersdk.NewSearchCond().
		Schema("sample").
		Select("title").
		Query("decentralized")
	)
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
