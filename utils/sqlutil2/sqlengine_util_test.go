package sqlutil2

import (
	"encoding/json"
	"fmt"
	"github.com/glitternetwork/glitter-sdk-go/utils"
	"testing"
)

func TestScanRows_2(t *testing.T) {
	str := `
{"result":[{"row":{"_id":{"value":"2","column_value_type":"String"},"ch_name":{"value":"李四","column_value_type":"String"},"gender":{"value":"male","column_value_type":"String"},"name":{"value":"bob","column_value_type":"String"},"size":{"value":"2234","column_value_type":"Float"}}},{"row":{"_id":{"value":"2","column_value_type":"String"},"ch_name":{"value":"王五","column_value_type":"String"},"gender":{"value":"male","column_value_type":"String"},"name":{"value":"mark","column_value_type":"String"},"size":{"value":"2234","column_value_type":"Float"}}}],"code":0,"full_took_times":0.005606036,"trace_id":"90dd2dac-d2db-48dc-9493-26f34e9bcc64"}`

	a := &GatewayResponse{}
	err := json.Unmarshal([]byte(str), a)
	fmt.Println(err)
	fmt.Println(a.Result)

	type DestType struct {
		ID     int     `db:"_id"`
		CHName string  `db:"ch_name"`
		Gender string  `db:"gender"`
		Name   string  `db:"name"`
		Size   float64 `db:"size"`
	}
	var dest []*DestType

	ScanRows(a, &dest)
	fmt.Println(utils.ConvToJSON(dest))

	//col := make([]string, 0, len(rs.ColumnDefs))
	//for _, cd := range rs.ColumnDefs {
	//	col = append(col, cd.ColumnName)
	//}

	//rows := Rows{
	//	rs:   rs,
	//	cols: col,
	//	idx:  -1,
	//	err:  nil,
	//}
	//return sqlx.StructScan(&rows, dest)
}
