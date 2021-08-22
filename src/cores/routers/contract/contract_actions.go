package contract

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/validate"
	"cos-backend-com/src/cores/routers"
	"cos-backend-com/src/libs/apierror"
	"cos-backend-com/src/libs/models/contractmodels"
	"cos-backend-com/src/libs/sdk/cores"
	"net/http"

	"github.com/wujiu2020/strip/utils/apires"
)

type ContractActionsHandler struct {
	routers.Base
}

func (h *ContractActionsHandler) Create(actionType cores.ContractActionType) (res interface{}) {
	var input cores.CreateContractActionInput
	if err := h.Params.BindJsonBody(&input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	if err := validate.Default.Struct(input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var uid flake.ID
	h.Ctx.Find(&uid, "uid")

	if err := contractmodels.ContractActions.Create(h.Ctx, uid, actionType, &input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(http.StatusOK)
	return
}

func (h *ContractActionsHandler) Get(actionType cores.ContractActionType, id string) (res interface{}) {
	var output cores.ContractActionResult

	var uid flake.ID
	h.Ctx.Find(&uid, "uid")

	if err := contractmodels.ContractActions.Get(h.Ctx, id, uid, actionType, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *ContractActionsHandler) List(actionType cores.ContractActionType) (res interface{}) {
	var output []cores.ContractActionResult

	var uid flake.ID
	h.Ctx.Find(&uid, "uid")

	if err := contractmodels.ContractActions.List(h.Ctx, uid, actionType, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}
