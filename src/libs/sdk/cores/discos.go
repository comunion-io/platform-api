package cores

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/pagination"
	"cos-backend-com/src/libs/sdk/eth"
	"time"
)

type DiscoState int

const (
	DiscoStateDefault DiscoState = iota
	DiscoStateWaitingForEnable
	DiscoStateCreateFailed
	DiscoStateEnabled
	DiscoStateWaitingForStart
	DiscoStateFundraisingSuccess
	DiscoStateFundraisingFailed
	DiscoStateFundraising
	DiscoStateFundraisingComplete
)

// DiscosModel represents a row from 'discos'.
type DiscosModel struct {
	Id                   flake.ID   `json:"id" db:"id"`                                        // id (PK)
	StartupId            flake.ID   `json:"startupId" db:"startup_id"`                         // startup_id
	WalletAddr           string     `json:"walletAddr" db:"wallet_addr"`                       // wallet_addr
	TokenAddr            string     `json:"tokenAddr" db:"token_addr"`                         // token_addr
	Description          string     `json:"description" db:"description"`                      // description
	FundRaisingStartedAt time.Time  `json:"fundRaisingStartedAt" db:"fund_raising_started_at"` // fund_raising_started_at
	FundRaisingEndedAt   time.Time  `json:"fundRaisingEndedAt" db:"fund_raising_ended_at"`     // fund_raising_ended_at
	InvestmentReward     int64      `json:"investmentReward" db:"investment_reward"`           // investment_reward
	RewardDeclineRate    int        `json:"rewardDeclineRate" db:"reward_decline_rate"`        // reward_decline_rate
	ShareToken           int64      `json:"shareToken" db:"share_token"`                       // share_token
	MinFundRaising       int64      `json:"minFundRaising" db:"min_fund_raising"`              // min_fund_raising
	AddLiquidityPool     int64      `json:"addLiquidityPool" db:"add_liquidity_pool"`          // add_liquidity_pool
	TotalDepositToken    float64    `json:"totalDepositToken" db:"total_deposit_token"`        // total_deposit_token
	State                DiscoState `json:"state" db:"state"`                                  // state
	FundRaisingAddr      string     `json:"fundRaisingAddr" db:"fund_raising_addr"`            // fund_raising_addr
	CreatedAt            time.Time  `json:"createdAt" db:"created_at"`                         // created_at
	UpdatedAt            time.Time  `json:"updatedAt" db:"updated_at"`                         // updated_at
}

// DiscosInput represents an input for 'discos'.
type CreateDiscosInput struct {
	Id                   flake.ID  `json:"id" validate:"required"`
	WalletAddr           string    `json:"walletAddr" validate:"required"`
	TokenAddr            string    `json:"tokenAddr" validate:"required"`
	Description          string    `json:"description" validate:"required"`
	FundRaisingStartedAt time.Time `json:"fundRaisingStartedAt" validate:"required"`
	FundRaisingEndedAt   time.Time `json:"fundRaisingEndedAt" validate:"required"`
	InvestmentReward     int64     `json:"investmentReward" validate:"required"`
	RewardDeclineRate    int       `json:"rewardDeclineRate" validate:"required"`
	ShareToken           int64     `json:"shareToken" validate:"required"`
	MinFundRaising       int64     `json:"minFundRaising" validate:"required"`
	AddLiquidityPool     int64     `json:"addLiquidityPool" validate:"required"`
	TotalDepositToken    float64   `json:"totalDepositToken" validate:"required"`
	TxId                 string    `json:"txId" validate:"required"`
}

// StartupDiscosResult represents an output for 'discos'.
type StartupDiscosResult struct {
	Id                   flake.ID             `json:"id" db:"id"`                                        // id (PK)
	WalletAddr           string               `json:"walletAddr" db:"wallet_addr"`                       // wallet_addr
	TokenAddr            string               `json:"tokenAddr" db:"token_addr"`                         // token_addr
	Description          string               `json:"description" db:"description"`                      // description
	FundRaisingStartedAt time.Time            `json:"fundRaisingStartedAt" db:"fund_raising_started_at"` // fund_raising_started_at
	FundRaisingEndedAt   time.Time            `json:"fundRaisingEndedAt" db:"fund_raising_ended_at"`     // fund_raising_ended_at
	InvestmentReward     int64                `json:"investmentReward" db:"investment_reward"`           // investment_reward
	RewardDeclineRate    int                  `json:"rewardDeclineRate" db:"reward_decline_rate"`        // reward_decline_rate
	ShareToken           int64                `json:"shareToken" db:"share_token"`                       // share_token
	MinFundRaising       int64                `json:"minFundRaising" db:"min_fund_raising"`              // min_fund_raising
	AddLiquidityPool     int64                `json:"addLiquidityPool" db:"add_liquidity_pool"`          // add_liquidity_pool
	TotalDepositToken    float64              `json:"totalDepositToken" db:"total_deposit_token"`        // total_deposit_token
	State                DiscoState           `json:"state" db:"state"`                                  // state
	TxId                 string               `json:"txId" db:"tx_id"`                                   // tx_id
	FundRaisingAddr      string               `json:"fundRaisingAddr" db:"fund_raising_addr"`            // fund_raising_addr
	TransactionState     eth.TransactionState `json:"transactionState" db:"transaction_state"`           // transaction_state
}

// ListDiscoResult represents an output for 'discos'.
type DiscoOutput struct {
	Id                flake.ID            `json:"id" db:"id"`                                 // id (PK)
	Startup           DiscosStartupResult `json:"startup" db:"startup"`                       // startup
	InvestmentReward  int64               `json:"investmentReward" db:"investment_reward"`    // investment_reward
	RewardDeclineRate int                 `json:"rewardDeclineRate" db:"reward_decline_rate"` // reward_decline_rate
	ShareToken        int64               `json:"shareToken" db:"share_token"`                // share_token
	MinFundRaising    int64               `json:"minFundRaising" db:"min_fund_raising"`       // min_fund_raising
	AddLiquidityPool  int64               `json:"addLiquidityPool" db:"add_liquidity_pool"`   // add_liquidity_pool
	State             DiscoState          `json:"state" db:"state"`                           // state
	FundRaisingAddr   string              `json:"fundRaisingAddr" db:"fund_raising_addr"`     // fund_raising_addr
}

type DiscosStartupResult struct {
	Id          flake.ID `json:"id" db:"id"`                    // id (PK)
	Name        string   `json:"name" db:"name"`                // name
	Logo        string   `json:"logo" db:"logo"`                // logo
	TokenSymbol string   `json:"tokenSymbol" db:"token_symbol"` // name
}

type ListDiscosOrderBy string

const (
	ListDiscosOrderByTime             ListDiscosOrderBy = "time"
	ListDiscosOrderByName             ListDiscosOrderBy = "name"
	ListDiscosOrderByInvestmentReward ListDiscosOrderBy = "investmentReward"
	ListDiscosOrderByLiquidityPool    ListDiscosOrderBy = "liquidityPool"
)

type ListDiscosInput struct {
	pagination.ListRequest
	Keyword string             `param:"keyword"`
	OrderBy *ListDiscosOrderBy `param:"orderBy"`
	IsAsc   bool               `param:"isAsc"`
}

type ListDiscosResult struct {
	pagination.ListResult
	Result []DiscoOutput `json:"result"`
}

type StatDiscoEthIncreaseInput struct {
	TimeFrom time.Time `json:"timeFrom" validate:"required"`
	TimeTo   time.Time `json:"timeTo" validate:"required"`
}

type StatDiscoEthIncreaseOutput struct {
	Date  time.Time `json:"date" db:"date"`
	Count int64     `json:"count" db:"count"`
}

type StatDiscoEthTotalInput struct {
	StartupId *flake.ID `json:"startupId" validate:"startup_id"`
}

type StatDiscoEthTotalResult struct {
	Count int64 `json:"count" db:"count"`
}

type StatDiscoTotalResult struct {
	Count        int64   `json:"count" db:"count"`
	IcreaseCount int64   `json:"-" db:"increase_count"`
	Rate         float64 `json:"rate" db:"rate"`
}

type DiscoSwapStateOutput struct {
	DiscoState int `json:"discoState" db:"disco_state"`
	SwapState  int `json:"swapState" db:"swap_state"`
}
