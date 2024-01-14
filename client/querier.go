package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	glittertypes "github.com/glitternetwork/chain-dep/glitter_proto/blockved/glitterchain/index/types"
	glittercommon "github.com/glitternetwork/chain-dep/glitter_proto/common"

	"io/ioutil"
	"net/url"
	"reflect"
	"strconv"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/glitternetwork/glitter-sdk-go/msg"
	"github.com/glitternetwork/glitter-sdk-go/tx"
	"github.com/glitternetwork/glitter-sdk-go/utils/sqlutil"
	"github.com/pkg/errors"
	"golang.org/x/net/context/ctxhttp"
)

// QueryAccountResData response
type QueryAccountResData struct {
	Address       msg.AccAddress `json:"address"`
	AccountNumber msg.Int        `json:"account_number"`
	Sequence      msg.Int        `json:"sequence"`
}

// QueryAccountRes response
type QueryAccountRes struct {
	Account QueryAccountResData `json:"account"`
}

// LoadAccount simulates gas and fee for a transaction
func (lcd *LCDClient) LoadAccount(ctx context.Context, address msg.AccAddress) (res authtypes.AccountI, err error) {
	resp, err := ctxhttp.Get(ctx, lcd.c, lcd.URL+fmt.Sprintf("/cosmos/auth/v1beta1/accounts/%s", address))
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to estimate")
	}
	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to read response")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 response code %d: %s", resp.StatusCode, string(out))
	}

	var response authtypes.QueryAccountResponse
	//err = json.Unmarshal(out, &response)
	err = lcd.GetMarshaler().UnmarshalJSON(out, &response)
	//err = lcd.EncodingConfig.Marshaller.UnmarshalJSON(out, &response)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to unmarshal response")
	}

	//fmt.Printf("LoadAccount address=%s resp=%+v\n", address.String(), response)
	return response.Account.GetCachedValue().(authtypes.AccountI), nil
}

// Simulate tx and get response
func (lcd *LCDClient) Simulate(ctx context.Context, txbuilder tx.Builder, options CreateTxOptions) (*sdktx.SimulateResponse, error) {
	// Create an empty signature literal as the ante handler will populate with a
	// sentinel pubkey.
	sig := signing.SignatureV2{
		PubKey: &secp256k1.PubKey{},
		Data: &signing.SingleSignatureData{
			SignMode: options.SignMode,
		},
		Sequence: options.Sequence,
	}
	if err := txbuilder.SetSignatures(sig); err != nil {
		return nil, err
	}

	bz, err := txbuilder.GetTxBytes()
	if err != nil {
		return nil, err
	}
	reqBytes, err := lcd.GetMarshaler().MarshalJSON(&sdktx.SimulateRequest{
		Tx:      nil,
		TxBytes: bz,
	})
	if err != nil {
		return nil, err
	}

	resp, err := ctxhttp.Post(ctx, lcd.c, lcd.URL+"/cosmos/tx/v1beta1/simulate", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to estimate")
	}
	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to read response")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 response code %d: %s", resp.StatusCode, string(out))
	}

	var response sdktx.SimulateResponse
	err = lcd.GetMarshaler().UnmarshalJSON(out, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// QueryScan execute a SQL query statement and scan result to target
// Args:
//   - target: Target to scan result values
//   - sql: The SQL query string
//   - args: Optional list of arguments to substitute into the query
//
// Returns:
// A list of rows where each row is a dict mapping column name to value
func (lcd *LCDClient) QueryScan(ctx context.Context, target interface{}, sql string, args ...*glittercommon.Argument) error {
	rt := reflect.TypeOf(target)
	if rt.Kind() != reflect.Ptr {
		return fmt.Errorf("result must be ptr slice")
	}
	if rt.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("result must be ptr slice")
	}
	resp, err := lcd.Query(ctx, sql, args...)
	if err != nil {
		return err
	}
	if len(resp.Results) < 1 {
		return errors.New("invalid query result")
	}
	err = sqlutil.ScanRows(resp.Results[0], target)
	if err != nil {
		return errors.Errorf("failed to convert query result: err=%v", err)
	}
	return nil
}

// Query execute a SQL query statement
// Args:
//   - sql: The SQL query string
//   - args: Optional list of arguments to substitute into the query
//
// Returns:
// A list of rows where each row is a dict mapping column name to value
func (lcd *LCDClient) Query(ctx context.Context, sql string, args ...*glittercommon.Argument) (res *glittertypes.SQLQueryResponse, err error) {
	req := glittertypes.SQLQueryRequest{
		Sql:       sql,
		Arguments: args,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to gen request")
	}

	resp, err := ctxhttp.Post(ctx, lcd.c, lcd.URL+"/blockved/glitterchain/index/sql/query", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to get schema")
	}
	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to read response")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 response code %d: %s", resp.StatusCode, string(out))
	}

	var response glittertypes.SQLQueryResponse
	err = lcd.GetMarshaler().UnmarshalJSON(out, &response)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to unmarshal response")
	}
	return &response, nil
}

// ListTables List tables in glitter, filtering by various criteria
// Args:
//   - tableKeyword: Filter tables by keyword
//   - uid: Filter tables by creator uid
//   - database: Filter tables by database name
//   - page: Page number for pagination
//   - pageSize: Number of results per page
//
// Returns:
// ListTablesResponse containing matching tables
func (lcd *LCDClient) ListTables(ctx context.Context, tableKeyword, uid, database string, page, pageSize *int) (res *glittertypes.SQLListTablesResponse, err error) {
	uv := url.Values{}
	if len(tableKeyword) > 0 {
		uv.Add("keyword", tableKeyword)
	}
	if len(uid) > 0 {
		uv.Add("uid", uid)
	}
	if len(database) > 0 {
		uv.Add("database", database)
	}
	if page != nil {
		uv.Add("page", strconv.Itoa(*page))
	}
	if pageSize != nil {
		uv.Add("page_size", strconv.Itoa(*pageSize))
	}

	resp, err := ctxhttp.Get(ctx, lcd.c, lcd.URL+fmt.Sprintf("/blockved/glitterchain/index/sql/list_tables?%s", uv.Encode()))
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to get doc")
	}
	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to read response")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 response code %d: %s", resp.StatusCode, string(out))
	}

	var response glittertypes.SQLListTablesResponse
	err = lcd.GetMarshaler().UnmarshalJSON(out, &response)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to unmarshal response")
	}

	return &response, nil
}

// ListDatabases List all databases or filter by creator in glitter
// Args:
//   - creator (optional): Only return databases created by this creator
//
// Returns:
// ListDatabasesResponse containing matching databases
func (lcd *LCDClient) ListDatabases(ctx context.Context, creator string) (res *glittertypes.SQLListDatabasesResponse, err error) {
	resp, err := ctxhttp.Get(ctx, lcd.c, lcd.URL+"/blockved/glitterchain/index/sql/list_databases")
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to get doc")
	}
	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to read response")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 response code %d: %s", resp.StatusCode, string(out))
	}

	var response glittertypes.SQLListDatabasesResponse
	err = lcd.GetMarshaler().UnmarshalJSON(out, &response)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to unmarshal response")
	}

	if len(creator) > 0 {
		filtered := make([]*glittercommon.DatabaseInfo, 0, len(response.Databases))
		for _, e := range response.Databases {
			if e.Creator == creator {
				filtered = append(filtered, e)
			}
		}
		response.Databases = filtered
	}
	return &response, nil
}

// ShowCreateTable Show the CREATE TABLE statement for an existing table
// Args:
// - database: The database name
// - table: The table name
//
// Returns:
// The result containing the CREATE TABLE statement
func (lcd *LCDClient) ShowCreateTable(ctx context.Context, database string, table string) (res *glittertypes.ShowCreateTableResponse, err error) {
	url := fmt.Sprintf("%s/blockved/glitterchain/index/sql/show_create_table?databaseName=%s&tableName=%s", lcd.URL, database, table)
	resp, err := ctxhttp.Get(ctx, lcd.c, url)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to get doc")
	}
	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to read response")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 response code %d: %s", resp.StatusCode, string(out))
	}

	var response glittertypes.ShowCreateTableResponse
	err = lcd.EncodingConfig.Marshaller.UnmarshalJSON(out, &response)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to unmarshal response")
	}

	return &response, nil
}
