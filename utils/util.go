package utils

import (
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ethereum/go-ethereum/common"
	ethermint "github.com/evmos/ethermint/types"
	curl "github.com/idoubi/goz"
	jsoniter "github.com/json-iterator/go"
	abci "github.com/tendermint/tendermint/abci/types"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/types"
)

const AccountAddressPrefix = "glitter"

func GetEvmAddrFromGlitterAddr(glitterAddr string) (string, error) {
	//accountPubKeyPrefix := AccountAddressPrefix + "pub"
	//validatorAddressPrefix := AccountAddressPrefix + "valoper"
	//validatorPubKeyPrefix := AccountAddressPrefix + "valoperpub"
	//consNodeAddressPrefix := AccountAddressPrefix + "valcons"
	//consNodePubKeyPrefix := AccountAddressPrefix + "valconspub"
	//config := sdk.GetConfig()
	//config.SetBech32PrefixForAccount(AccountAddressPrefix, accountPubKeyPrefix)
	//config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	//config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	//SetBip44CoinType(config)
	accAddr, err := sdk.AccAddressFromBech32(glitterAddr)
	if err != nil {
		return "", err
	}
	evmAddr := common.BytesToAddress(accAddr).String()
	return evmAddr, nil
}

func GetGlitterAddrFromEvmAddr(evmAddr string) (string, error) {
	//accountPubKeyPrefix := AccountAddressPrefix + "pub"
	//validatorAddressPrefix := AccountAddressPrefix + "valoper"
	//validatorPubKeyPrefix := AccountAddressPrefix + "valoperpub"
	//consNodeAddressPrefix := AccountAddressPrefix + "valcons"
	//consNodePubKeyPrefix := AccountAddressPrefix + "valconspub"
	//config := sdk.GetConfig()
	//config.SetBech32PrefixForAccount(AccountAddressPrefix, accountPubKeyPrefix)
	//config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	//config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	//SetBip44CoinType(config)

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

// curlGet ...
func CurlGet(reqUrl string) (data string, err error) {
	headers := map[string]interface{}{
		"accept":       "application/json",
		"Content-Type": "application/json",
	}
	option := curl.Options{
		Headers: headers,
	}
	response, err := curl.NewClient(option).Get(reqUrl)
	if err == nil {
		if body, e := response.GetBody(); e == nil {
			if response.GetStatusCode() == 200 {
				return body.String(), nil
			} else { // nolint
				return "", errors.New(fmt.Sprintf("Http status %d body %s", response.GetStatusCode(), body.String()))
			}
		} else {
			return "", e
		}
	}

	return "", err
}

// curlPost ...
func CurlPost(reqUrl string, postData map[string]interface{}) (data string, err error) {
	headers := map[string]interface{}{
		"accept":       "application/json",
		"Content-Type": "application/json",
	}
	option := curl.Options{
		Headers: headers,
		JSON:    postData,
	}
	response, err := curl.NewClient().Post(reqUrl, option)

	if err == nil {
		if body, e := response.GetBody(); e == nil {
			if response.GetStatusCode() == 200 {
				return body.String(), nil
			} else { // nolint
				return "", errors.New(fmt.Sprintf("Http status %d body %s", response.GetStatusCode(), body.String()))
			}
		} else {
			return "", e
		}
	}

	return "", err
}

// CurlGetV2 ...
func CurlGetV2(reqUrl string, options curl.Options) (data string, err error) {
	response, err := curl.NewClient(options).Get(reqUrl)
	if err == nil {
		if body, e := response.GetBody(); e == nil {
			if response.GetStatusCode() == 200 {
				return body.String(), nil
			} else { // nolint
				return "", errors.New(fmt.Sprintf("Http status %d body %s", response.GetStatusCode(), body.String()))
			}
		} else {
			return "", e
		}
	}

	return "", err
}

// ConvertValidatorAddressToBech32String converts the given validator address to its Bech32 string representation
func ConvertValidatorAddressToBech32String(address types.Address) string {
	return sdk.ConsAddress(address).String()
}

// ConvertValidatorPubKeyToBech32String converts the given pubKey to a Bech32 string
func ConvertValidatorPubKeyToBech32String(pubKey tmcrypto.PubKey) (string, error) {
	bech32Prefix := sdk.GetConfig().GetBech32ConsensusPubPrefix()
	return bech32.ConvertAndEncode(bech32Prefix, pubKey.Bytes())
}

func FindEventByType(events []abci.Event, eventType string) (abci.Event, error) {
	for _, event := range events {
		if event.Type == eventType {
			return event, nil
		}
	}

	return abci.Event{}, fmt.Errorf("no event with type %s found", eventType)
}

func FindEventsByType(events []abci.Event, eventType string) []abci.Event {
	var found []abci.Event
	for _, event := range events {
		if event.Type == eventType {
			found = append(found, event)
		}
	}

	return found
}

func FindAttributeByKey(event abci.Event, attrKey string) (abci.EventAttribute, error) {
	for _, attr := range event.Attributes {
		if string(attr.Key) == attrKey {
			return attr, nil
		}
	}

	return abci.EventAttribute{}, fmt.Errorf("no attribute with key %s found inside event with type %s", attrKey, event.Type)
}

var Json = jsoniter.ConfigCompatibleWithStandardLibrary

func ConvToJSON(v interface{}) (str string) {
	bytes, _ := Json.Marshal(v)
	return string(bytes)
}
