package cores

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/pagination"
)

type ProposalStatus int

const (
	ProposalStatusPending    ProposalStatus = 0
	ProposalStatusNotStarted ProposalStatus = 1
	ProposalStatusVoting     ProposalStatus = 2
	ProposalStatusOver       ProposalStatus = 3
)

type CreateTermInput struct {
	Amount  float64 `json:"amount"`
	Content string  `json:"content"`
}

type CreateProposalInput struct {
	Id                        flake.ID          `json:"id" validate:"required"`
	TxId                      string            `json:"txId" validate:"required"`
	StartupId                 string            `json:"startupId" validate:"required"`
	WalletAddr                string            `json:"walletAddr" validate:"required"`
	ContractAddr              string            `json:"contractAddr" validate:"required"`
	Status                    *ProposalStatus   `json:"status" validate:"required"`
	Title                     string            `json:"title" validate:"required"`
	Type                      int               `json:"type" validate:"required"`
	UserId                    flake.ID          `json:"userId"`
	Contact                   string            `json:"contact" validate:"required"`
	Description               string            `json:"description" validate:"required"`
	VoterType                 *int              `json:"voterType" validate:"required"`
	Supporters                *int              `json:"supporters" validate:"required"`
	MinimumApprovalPercentage *int              `json:"minApprovalPercent" validate:"required"`
	Duration                  *int              `json:"duration" validate:"required"`
	HasPayment                *bool             `json:"hasPayment" validate:"required"`
	PaymentAddr               string            `json:"paymentAddr"`
	PaymentType               int               `json:"paymentType"`
	PaymentMonths             int               `json:"paymentMonths"`
	PaymentDate               string            `json:"paymentDate"`
	PaymentAmount             float64           `json:"paymentAmount"`
	TotalPaymentAmount        float64           `json:"totalPaymentAmount"`
	Terms                     []CreateTermInput `json:"terms"`
}

type CreateProposalResult struct {
	Id     flake.ID       `json:"id" db:"id"`
	Status ProposalStatus `json:"status" db:"status"`
}

type UpdateProposalStatusInput struct {
	Id     flake.ID       `json:"id"`
	Status ProposalStatus `json:"status" validate:"func=parent.Validate"`
}

func (u UpdateProposalStatusInput) Validate() bool {
	v := u.Status
	return v == 4 || v == 5 || v == 6
}

type UpdateProposalStatusResult struct {
	Id flake.ID `json:"id"`
}

type VoteProposalInput struct {
	Id         flake.ID `json:"id"`
	TxId       string   `json:"txId" validate:"required"`
	Amount     float32  `json:"amount" validate:"required"`
	IsApproved bool     `json:"isApproved" validate:"required"`
	WalletAddr string   `json:"walletAddr" validate:"required"`
	CreatedAt  string   `json:"createdAt"`
}

type VoteProposalResult struct {
}

type CreateTermResult struct {
	Id flake.ID `json:"id" db:"id"`
}

type GetProposalInput struct {
	Id   flake.ID `json:"id"`
	TxId string   `json:"txId"`
}

type ProposalResult struct {
	Id      flake.ID `json:"id" db:"id"`
	TxId    string   `json:"txId" db:"tx_id"`
	Startup struct {
		Id          flake.ID `json:"id" db:"id"`
		Name        string   `json:"name" db:"name"`
		Logo        string   `json:"logo" db:"logo"`
		TokenSymbol string   `json:"tokenSymbol" db:"token_symbol"`
	} `json:"startup" db:"startup"`
	Comer struct {
		Id   flake.ID `json:"id" db:"id"`
		Name string   `json:"name" db:"name"`
	} `json:"comer" db:"comer"`
	WalletAddr         string         `json:"walletAddr" db:"wallet_addr"`
	ContractAddr       string         `json:"contractAddr" db:"contract_addr"`
	CreatedAt          string         `json:"createdAt" db:"created_at"`
	UpdatedAt          string         `json:"updatedAt" db:"updated_at"`
	Status             ProposalStatus `json:"status" db:"status"`
	Title              string         `json:"title" db:"title"`
	Type               int            `json:"type" db:"type"`
	Contact            string         `json:"contact" db:"contact"`
	Description        string         `json:"description" db:"description"`
	VoterType          int            `json:"voterType" db:"voter_type"`
	Supporters         int            `json:"supporters" db:"supporters"`
	MinApprovalPercent int            `json:"minApprovalPercent" db:"minimum_approval_percentage"`
	Duration           int            `json:"duration" db:"duration"`
	HasPayment         bool           `json:"hasPayment" db:"has_payment"`
	PaymentAddr        string         `json:"paymentAddr" db:"payment_addr"`
	PaymentType        int            `json:"paymentType" db:"payment_type"`
	PaymentMonths      int            `json:"paymentMonths" db:"payment_months"`
	PaymentDate        string         `json:"paymentDate" db:"payment_date"`
	PaymentAmount      float64        `json:"paymentAmount" db:"payment_amount"`
	TotalPaymentAmount float64        `json:"totalPaymentAmount" db:"total_payment_amount"`
	Votes              []struct {
		Amount     float64 `json:"amount" db:"amount"`
		IsApproved bool    `json:"isApproved" db:"is_approved"`
		WalletAddr string  `json:"walletAddr" db:"wallet_addr"`
		CreatedAt  string  `json:"createdAt" db:"created_at"`
	} `json:"votes" db:"votes"`
	Terms []struct {
		Amount  float64 `json:"amount" db:"amount"`
		Content string  `json:"content" db:"content"`
	} `json:"terms" db:"terms"`
}

const (
	ListProposalsTypeAll     string = "all"
	ListProposalsTypeCreated string = "created"
	ListProposalsTypeVoted   string = "voted"
)

type ListProposalsInput struct {
	pagination.ListRequest
	Keyword   string   `param:"keyword"`
	Type      string   `param:"type" validate:"required"`
	StartupId flake.ID `param:"startupId"`
	Statuses  string   `param:"statuses[]"`
	OrderBy   string   `param:"orderBy"`
	IsDesc    bool     `param:"isDesc"`
}

type ListProposalsResult struct {
	pagination.ListResult
	Result []struct {
		Id      flake.ID `json:"id" db:"id"`
		Startup struct {
			Id          flake.ID `json:"id" db:"id"`
			Name        string   `json:"name" db:"name"`
			Logo        string   `json:"logo" db:"logo"`
			TokenSymbol string   `json:"tokenSymbol" db:"token_symbol"`
		} `json:"startup" db:"startup"`
		Comer struct {
			Id   flake.ID `json:"id" db:"id"`
			Name string   `json:"name" db:"name"`
		} `json:"comer" db:"comer"`
		Title              string         `json:"title" db:"title"`
		Status             ProposalStatus `json:"status" db:"status"`
		HasPayment         bool           `json:"hasPayment" db:"has_payment"`
		TotalPaymentAmount float64        `json:"totalPaymentAmount" db:"total_payment_amount"`
		CreatedAt          string         `json:"createdAt" db:"created_at"`
		UpdatedAt          string         `json:"updatedAt" db:"updated_at"`
		Duration           int            `json:"duration" db:"duration"`
	} `json:"result"`
}

type ProposalOverResult struct {
	Done bool `json:"done" db:"done"`
}
