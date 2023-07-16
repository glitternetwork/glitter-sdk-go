package sqlutil

import (
	"testing"

	glittertypes "github.com/glitternetwork/glitter.proto/golang/glitter_proto/index/types"
)

func TestScanRows(t *testing.T) {
	type args struct {
		rs   *glittertypes.ResultSet
		dest interface{}
	}
	type DestType struct {
		A int     `db:"a"`
		B string  `db:"b"`
		C bool    `db:"c"`
		D float64 `db:"d"`
	}
	var dest []*DestType
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				rs: &glittertypes.ResultSet{
					Id: "",
					ColumnDefs: []*glittertypes.ColumnDef{
						{ColumnName: "a", ColumnType: "int"},
						{ColumnName: "b", ColumnType: "string"},
						{ColumnName: "c", ColumnType: "bool"},
						{ColumnName: "d", ColumnType: "float"},
					},
					Rows: []*glittertypes.RowData{
						{
							Columns: []string{"1", "abcde", "true", "2.33"},
						},
						{
							Columns: []string{"2", "bcdef", "false", "-2.33"},
						},
					},
				},
				dest: &dest,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ScanRows(tt.args.rs, tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("ScanRows() error = %v, wantErr %v", err, tt.wantErr)
			}
			for _, v := range dest {
				t.Logf("%+v", v)
			}
		})
	}
}
