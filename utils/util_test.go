package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_GetEvmAddrFromGlitterAddr(t *testing.T) {
	evmAddr, err := GetEvmAddrFromGlitterAddr("glitter1q5t6prazp4wlzegvaz5sls25phyvvmq6aqpxf7")
	assert.Equal(t, strings.ToLower(evmAddr), strings.ToLower("0x0517a08fa20d5Df1650ce8a90fc1540dc8c66c1a"))
	fmt.Println(err)
	assert.Nil(t, err)
}

func Test_GetGlitterAddrFromEvmAddr(t *testing.T) {
	glitterAddr, err := GetGlitterAddrFromEvmAddr("0xb97E46d90EBBA78Fdf8dEbd0D12c312821bFFFFE")
	fmt.Println(glitterAddr)
	fmt.Println(err)
}
