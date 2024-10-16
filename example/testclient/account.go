package testclient

import (
	"github.com/glitternetwork/glitter-sdk-go/client"
	"github.com/glitternetwork/glitter-sdk-go/key"
)

func New() *client.LCDClient {
	const chainID = "glitter_12001-4"
	mnemonicKey := "drip feed dish dirt hold mushroom neutral vessel permit cost palace direct access piano attract crystal winner august sail amused cabin test glad prize"
	pk, err := key.DerivePrivKeyBz(mnemonicKey, key.CreateHDPath(0, 0))
	if err != nil {
		panic(err)
	}

	privKey, err := key.PrivKeyGen(pk)
	if err != nil {
		panic(err)
	}
	return client.New(chainID, privKey)
}
