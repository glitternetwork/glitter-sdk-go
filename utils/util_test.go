package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGetUpdateStatement(t *testing.T) {
	table := "demo_table"
	columns := map[string]interface{}{"name": "demo_name", "age": 18}
	where := map[string]interface{}{"author": "james", "tag": "aaa"}
	result, args, err := BuildUpdateStatement(table, columns, where)
	if err != nil {
		fmt.Printf("err=%+v\n", err)
	}
	fmt.Printf("TestGetUpdateStatement||sql=%s||args=%+v\n", result, args)
}

func TestGetInsertStatement(t *testing.T) {
	table := "demo_table"
	columns := map[string]interface{}{"name": "demo_name", "age": 18}
	result, args, err := BuildInsertStatement(table, columns)
	if err != nil {
		fmt.Printf("err=%+v\n", err)
	}
	fmt.Printf("TestGetUpdateStatement||sql=%s||args=%+v\n", result, args)
}

func TestGetBatchInsertStatement(t *testing.T) {
	table := "demo_table"
	columns := []string{"name", "age"}
	rowValues := [][]interface{}{
		{"name1", 3},
		{"name2", 4},
		{"name3", 5},
	}
	result, args, err := BuildBatchInsertStatement(table, columns, rowValues)
	if err != nil {
		fmt.Printf("err=%+v\n", err)
	}
	fmt.Printf("TestGetUpdateStatement||sql=%s||agrs=%+v\n", result, args)
}

func Test_GetEvmAddrFromGlitterAddr(t *testing.T) {
	evmAddr, err := GetEvmAddrFromGlitterAddr("glitter1q5t6prazp4wlzegvaz5sls25phyvvmq6aqpxf7")
	assert.Equal(t, strings.ToLower(evmAddr), strings.ToLower("0x0517a08fa20d5Df1650ce8a90fc1540dc8c66c1a"))
	fmt.Println(err)
	assert.Nil(t, err)
}

func Test_GetGlitterAddrFromEvmAddr(t *testing.T) {
	glitterAddr, err := GetGlitterAddrFromEvmAddr("0x0517a08fa20d5Df1650ce8a90fc1540dc8c66c1a")
	assert.Equal(t, glitterAddr, "glitter1q5t6prazp4wlzegvaz5sls25phyvvmq6aqpxf7")
	assert.Nil(t, err)
}
