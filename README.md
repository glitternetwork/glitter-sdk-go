# glitter-sdk-go

> glitter sdk for golang

## Quick Start

```go
	import "github.com/glitternetwork/glitter-sdk-go"
```

```go
	// Create a client.
	client := glittersdk.New()
	// Get detailed informatino of the schema.
	schema := client.DB().GetSchema("sample")
	// Put a document into Glitter.
	txID,err := client.DB().PutDoc("sample",glittersdk.Document(`{
		"url"  : "https://glitterprotocol.io/",
		"title": "A Decentralized Content Indexing Network"
	}`))
	// Get the document by primary key in the schema.
	docs, err := db.GetDocs("sample", []string{"https://glitterprotocol.io/"})
	// Search the document which matches the words in the query.
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

See [example_test.go](example_test.go)
