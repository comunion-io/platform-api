package cores

import (
	"cos-backend-com/src/common/flake"
	"time"
)

type ContractActionType int

const (
	ContractActionTypeDefault ContractActionType = iota
	ContractActionTypeDefaultDisco
)

type ConctractActionsModel struct {
	Id        string             `json:"id" db:"id"`
	Uid       flake.ID           `json:"uid" db:"uid"`
	Type      ContractActionType `json:"type" db:"type"`
	CreatedAt time.Time          `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time          `json:"updatedAt" db:"updated_at"`
}

type CreateContractActionInput struct {
	Id string `json:"id" validate:"required"`
}

type ContractActionResult struct {
	Id        string             `json:"id" db:"id"`
	Uid       flake.ID           `json:"uid" db:"uid"`
	Type      ContractActionType `json:"type" db:"type"`
	CreatedAt time.Time          `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time          `json:"updatedAt" db:"updated_at"`
}
