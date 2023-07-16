package client

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/glitternetwork/glitter-sdk-go/utils"
)

// CreateDatabase Create a new database with the specified name
// Args:
// - database: The name of the database to create
//
// Returns:
// The result of executing the SQL CREATE DATABASE statement
func (lcd *LCDClient) CreateDatabase(ctx context.Context, database string) (*sdk.TxResponse, error) {
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", database)
	return lcd.SQLExec(ctx, sql, nil)
}

// CreateTable Creates a new table in the database using the provided SQL DDL statement.
// table name must be a full path format <database>.<table>
// Args:
// - sql: The SQL statement to create a new table
//
// Returns:
// The result of executing the SQL statement
func (lcd *LCDClient) CreateTable(ctx context.Context, sql string) (*sdk.TxResponse, error) {
	return lcd.SQLExec(ctx, sql, nil)
}

// DropTable Drop (deletes) a table from the specified database
// Args:
//   - database: The database name
//   - table: The table name
//
// Returns:
// The result of executing the SQL DROP TABLE statement
func (lcd *LCDClient) DropTable(ctx context.Context, db, table string) (*sdk.TxResponse, error) {
	sql := fmt.Sprintf("DROP TABLE IF EXISTS %s.%s", db, table)
	return lcd.SQLExec(ctx, sql, nil)
}

// DropDatabase Drop (deletes) the specified database
// Args:
//   - database: The database name
//
// Returns:
// The result of executing the SQL DROP DATABASE statement
func (lcd *LCDClient) DropDatabase(ctx context.Context, db, table string) (*sdk.TxResponse, error) {
	sql := fmt.Sprintf("DROP DATABASE IF EXISTS %s", db)
	return lcd.SQLExec(ctx, sql, nil)
}

// Delete rows from the specified table based on the provided conditions
// Args:
//   - db: The database name
//   - table: The table name
//   - where: Condition to match rows to delete
//   - order_by: Column to order deletion
//   - asc: Sort order ascending if True
//   - limit: Max number of rows to delete
//
// Returns:
// The result of executing the DELETE statement
func (lcd *LCDClient) Delete(ctx context.Context, db, table string, where map[string]interface{}, orderBy string, asc bool, limit int) (*sdk.TxResponse, error) {
	sql, args, err := utils.BuildDeleteStatement(utils.FullTableName(db, table), where, orderBy, asc, limit)
	if err != nil {
		return nil, err
	}
	return lcd.SQLExec(ctx, sql, args)
}

// Insert a new row into the specified table with the provided column-value pairs
// Args:
//   - db: The database name
//   - table: The table name
//   - columns: A dictionary of column names and values to insert
//
// Returns:
// The result of executing the INSERT statement
func (lcd *LCDClient) Insert(ctx context.Context, db, table string, columns map[string]interface{}) (*sdk.TxResponse, error) {
	insertSql, args, err := utils.BuildInsertStatement(utils.FullTableName(db, table), columns)
	if err != nil {
		return nil, err
	}
	return lcd.SQLExec(ctx, insertSql, args)
}

// BatchInsert insert multiple rows into the specified table using the provided column names and row values
// Args:
//   - db: The database name
//   - table: The table name
//   - rows: A list of rows to insert, each row is a dict of column names to values
//
// Returns:
// The result of executing the batch INSERT statement
func (lcd *LCDClient) BatchInsert(ctx context.Context, db, table string, columns []string, rowValues [][]interface{}) (*sdk.TxResponse, error) {
	batchInsertSql, args, err := utils.BuildBatchInsertStatement(utils.FullTableName(db, table), columns, rowValues)
	if err != nil {
		return nil, err
	}
	return lcd.SQLExec(ctx, batchInsertSql, args)
}

// Update rows in the specified table with the provided column-value pairs based on the specified conditions
// Args:
//   - db: The database name
//   - table: The table name
//   - columnsValue: A dictionary of column names and updated values
//   - where: A dictionary of column names and values to match
//
// Returns:
// The result of executing the UPDATE statement
func (lcd *LCDClient) Update(ctx context.Context, db, table string, columns map[string]interface{}, where map[string]interface{}) (*sdk.TxResponse, error) {
	updateSql, args, err := utils.BuildUpdateStatement(utils.FullTableName(db, table), columns, where)
	if err != nil {
		return nil, err
	}
	return lcd.SQLExec(ctx, updateSql, args)
}
