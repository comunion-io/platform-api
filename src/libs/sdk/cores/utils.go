package cores

import (
	"cos-backend-com/src/common/flake"
)

type PrepareIdOutput struct {
	Id flake.ID `json:"id"`
}

type Token struct {
	Name   string `json:"name" db:"token_name"`
	Symbol string `json:"symbol" db:"token_symbol"`
}

type PayTokenListOutput struct {
	PayTokens []Token `json:"payTokens"`
}

func AvailableTokens(token Token) PayTokenListOutput {
	var arrTokens []Token = []Token{
		Token{"ETH", "ETH"},
		Token{"BTC", "BTC"},
		Token{"USDT", "USDT"},
		Token{"USD", "USD"},
		Token{"CNY", "CNY"},
	}
	arrTokens = append(arrTokens, token)
	var payTokens PayTokenListOutput = PayTokenListOutput{arrTokens}
	return payTokens
}
