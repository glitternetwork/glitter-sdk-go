package utils

import (
	"encoding/base64"
	"fmt"
	glittertypes "github.com/glitternetwork/glitter.proto/golang/glitter_proto/index/types"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethermint "github.com/evmos/ethermint/types"
)

func toGlitterArgument(columnValue interface{}) (*glittertypes.Argument, error) {
	arg := glittertypes.Argument{
		Type:  0,
		Value: "",
	}
	switch v := columnValue.(type) {
	case int, int8, int16, int32, int64:
		arg.Type = glittertypes.Argument_INT
		arg.Value = fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		arg.Type = glittertypes.Argument_UINT
		arg.Value = fmt.Sprintf("%d", v)
	case float32:
		arg.Type = glittertypes.Argument_FLOAT
		arg.Value = strconv.FormatFloat(float64(v), 'g', -1, 32)
	case float64:
		arg.Type = glittertypes.Argument_FLOAT
		arg.Value = strconv.FormatFloat(v, 'g', -1, 64)
	case string:
		arg.Type = glittertypes.Argument_STRING
		arg.Value = v
	case bool:
		arg.Type = glittertypes.Argument_BOOL
		arg.Value = strconv.FormatBool(v)
	case []byte:
		arg.Type = glittertypes.Argument_BYTES
		arg.Value = base64.StdEncoding.EncodeToString(v)
	default:
		return nil, fmt.Errorf("unsupported value type: %T", v)
	}
	return &arg, nil
}

func toGlitterArguments(rowIndex int, columnValues []interface{}) ([]*glittertypes.Argument, error) {
	args := make([]*glittertypes.Argument, 0, len(columnValues))
	for _, v := range columnValues {
		a, err := toGlitterArgument(v)
		if err != nil {
			return nil, fmt.Errorf("failed to convert argument: rowIndex=%d %w", rowIndex, err)
		}
		args = append(args, a)
	}
	return args, nil
}

func BuildBatchInsertStatement(table string, columns []string, rowValues [][]interface{}) (string, []*glittertypes.Argument, error) {
	args := make([]*glittertypes.Argument, 0, len(columns)*len(rowValues))
	sqlBuilder := strings.Builder{}
	write := func(s string, args ...interface{}) {
		sqlBuilder.WriteString(fmt.Sprintf(s, args...))
	}
	if len(columns) == 0 {
		return "", nil, fmt.Errorf("empty columns")
	}
	write("INSERT INTO %s (%s) VALUES ", table, strings.Join(columns, ","))
	rowPlaceholders := "(" + "?" + strings.Repeat(",?", len(columns)-1) + ")"
	for i, v := range rowValues {
		if len(v) != len(columns) {
			return "", nil, fmt.Errorf("column values length not match with columns: row_index=%d columns=%+v", i, columns)
		}
		if i > 0 {
			write(",")
		}
		write(rowPlaceholders)
		rowArgs, err := toGlitterArguments(i, v)
		if err != nil {
			return "", nil, err
		}
		args = append(args, rowArgs...)
	}
	return sqlBuilder.String(), args, nil
}

func BuildInsertStatement(table string, columnToValues map[string]interface{}) (string, []*glittertypes.Argument, error) {
	columns := make([]string, 0, len(columnToValues))
	values := make([]interface{}, 0, len(columnToValues))
	for col, val := range columnToValues {
		columns = append(columns, col)
		values = append(values, val)
	}
	return BuildBatchInsertStatement(table, columns, [][]interface{}{values})
}

// BuildUpdateStatement where connected by and
func BuildUpdateStatement(table string, columns map[string]interface{}, whereEqual map[string]interface{}) (string, []*glittertypes.Argument, error) {
	var setKey []string
	var whereKey []string
	args := make([]*glittertypes.Argument, 0, len(columns)+len(whereEqual))
	for k, v := range columns {
		setKey = append(setKey, fmt.Sprintf("%s=?", k))
		a, err := toGlitterArgument(v)
		if err != nil {
			return "", nil, fmt.Errorf("failed to convert argument: column=%s %w", k, err)
		}
		args = append(args, a)
	}
	for k, v := range whereEqual {
		whereKey = append(whereKey, fmt.Sprintf("%s=?", k))
		a, err := toGlitterArgument(v)
		if err != nil {
			return "", nil, fmt.Errorf("failed to convert argument: whereKey=%s %w", k, err)
		}
		args = append(args, a)
	}
	setKeyGather := strings.Join(setKey, ",")
	whereKeyGather := strings.Join(whereKey, " and ")
	updateSql := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, setKeyGather, whereKeyGather)
	return updateSql, args, nil
}

const AccountAddressPrefix = "glitter"

// BuildDeleteStatement where connected by and
func BuildDeleteStatement(table string, whereEqual map[string]interface{}, orderBy string, asc bool, limit int) (string, []*glittertypes.Argument, error) {
	var whereKey []string
	args := make([]*glittertypes.Argument, 0, len(whereEqual))
	for k, v := range whereEqual {
		whereKey = append(whereKey, fmt.Sprintf("%s=?", k))
		a, err := toGlitterArgument(v)
		if err != nil {
			return "", nil, fmt.Errorf("failed to convert argument: whereKey=%s %w", k, err)
		}
		args = append(args, a)
	}
	whereKeyGather := strings.Join(whereKey, " and ")
	updateSql := fmt.Sprintf("DELETE FROM %s WHERE %s", table, whereKeyGather)
	if len(orderBy) > 0 {
		orderBySC := "ASC"
		if !asc {
			orderBySC = "DESC"
		}
		updateSql = fmt.Sprintf("%s ORDER BY %s %s", updateSql, orderBy, orderBySC)
	}
	if limit > 0 {
		updateSql = fmt.Sprintf("%s LIMIT %d", updateSql, limit)
	}
	return updateSql, args, nil
}

func GetEvmAddrFromGlitterAddr(glitterAddr string) (string, error) {
	accountPubKeyPrefix := AccountAddressPrefix + "pub"
	validatorAddressPrefix := AccountAddressPrefix + "valoper"
	validatorPubKeyPrefix := AccountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := AccountAddressPrefix + "valcons"
	consNodePubKeyPrefix := AccountAddressPrefix + "valconspub"
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	SetBip44CoinType(config)
	accAddr, err := sdk.AccAddressFromBech32(glitterAddr)
	if err != nil {
		return "", err
	}
	evmAddr := common.BytesToAddress(accAddr).String()
	return evmAddr, nil
}

func GetGlitterAddrFromEvmAddr(evmAddr string) (string, error) {
	accountPubKeyPrefix := AccountAddressPrefix + "pub"
	validatorAddressPrefix := AccountAddressPrefix + "valoper"
	validatorPubKeyPrefix := AccountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := AccountAddressPrefix + "valcons"
	consNodePubKeyPrefix := AccountAddressPrefix + "valconspub"
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	SetBip44CoinType(config)

	glitterAddrStrFromEth := sdk.AccAddress(common.HexToAddress(evmAddr).Bytes()).String()
	accAddr, err := sdk.AccAddressFromBech32(glitterAddrStrFromEth)
	if err != nil {
		return "", err
	}
	return accAddr.String(), nil
}

func SetBip44CoinType(config *sdk.Config) {
	config.SetCoinType(ethermint.Bip44CoinType)
	config.SetPurpose(sdk.Purpose)                      // Shared
	config.SetFullFundraiserPath(ethermint.BIP44HDPath) // nolint: staticcheck
}

func FullTableName(db, table string) string {
	return fmt.Sprintf("%s.%s", db, table)
}
