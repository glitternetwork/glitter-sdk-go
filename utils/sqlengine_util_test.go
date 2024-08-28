package utils

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestScanRows_2(t *testing.T) {
	str := `
{"result":[{"row":{"_id":{"value":"2","column_value_type":"String"},"ch_name":{"value":"李四","column_value_type":"String"},"gender":{"value":"male","column_value_type":"String"},"name":{"value":"bob","column_value_type":"String"},"size":{"value":"2234","column_value_type":"Float"}}},{"row":{"_id":{"value":"2","column_value_type":"String"},"ch_name":{"value":"王五","column_value_type":"String"},"gender":{"value":"male","column_value_type":"String"},"name":{"value":"mark","column_value_type":"String"},"size":{"value":"2234","column_value_type":"Float"}}}],"code":0,"full_took_times":0.005606036,"trace_id":"90dd2dac-d2db-48dc-9493-26f34e9bcc64"}`

	gatewayResponse := &GatewayResponse{}
	err := json.Unmarshal([]byte(str), gatewayResponse)
	fmt.Println(err)

	type DestType struct {
		ID     int     `db:"_id"`
		CHName string  `db:"ch_name"`
		Gender string  `db:"gender"`
		Name   string  `db:"name"`
		Size   float64 `db:"size"`
	}
	var dest []*DestType = make([]*DestType, 0)
	ScanRows(gatewayResponse.Result, &dest)
	fmt.Println(ConvToJSON(dest))

	str = `{"result":[],"code":0,"full_took_times":0.005606036,"trace_id":"90dd2dac-d2db-48dc-9493-26f34e9bcc64"}`
	gatewayResponse = &GatewayResponse{}
	err = json.Unmarshal([]byte(str), gatewayResponse)
	fmt.Println(err)

	var dest2 []*DestType = make([]*DestType, 0)
	ScanRows(gatewayResponse.Result, &dest2)
	fmt.Println(ConvToJSON(dest))
}
