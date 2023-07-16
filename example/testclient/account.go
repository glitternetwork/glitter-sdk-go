package testclient

import (
	"github.com/glitternetwork/glitter-sdk-go/client"
	"github.com/glitternetwork/glitter-sdk-go/key"
)

func New() *client.LCDClient {
	const chainID = "glitter_12000-2"
	mnemonicKey := "lesson police usual earth embrace someone opera season urban produce jealous canyon shrug usage subject cigar imitate hollow route inhale vocal special sun fuel"
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
