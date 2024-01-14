package key

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateMnemonic(t *testing.T) {
	str, err := CreateMnemonic()
	assert.NoError(t, err)
	fmt.Println(str)
}

func Test_DrivePrivKey(t *testing.T) {
	mnemonic, err := CreateMnemonic()
	assert.NoError(t, err)

	// Only Secp256k1 is supported
	_, err = DerivePrivKeyBz(mnemonic, CreateHDPath(1, 1))
	assert.NoError(t, err)
}
