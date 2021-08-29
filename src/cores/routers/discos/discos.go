package discos

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/validate"
	"cos-backend-com/src/cores/routers"
	"cos-backend-com/src/libs/apierror"
	"cos-backend-com/src/libs/models/discomodels"
	"cos-backend-com/src/libs/sdk/cores"
	"net/http"
	"time"

	"github.com/wujiu2020/strip/utils/apires"
)

type DiscosHandler struct {
	routers.Base
}

func (h *DiscosHandler) CreateStartupDisco(startupId flake.ID) (res interface{}) {
	var input cores.CreateDiscosInput
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

	if err := discomodels.Discos.CreateDisco(h.Ctx, startupId, &input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(http.StatusOK)
	return
}

func (h *DiscosHandler) GetStartupDisco(startupId flake.ID) (res interface{}) {
	var output cores.StartupDiscosResult
	if err := discomodels.Discos.GetDisco(h.Ctx, startupId, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *DiscosHandler) ListDisco() (res interface{}) {
	var input cores.ListDiscosInput
	h.Params.BindValuesToStruct(&input)

	if err := validate.Default.Struct(input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var output cores.ListDiscosResult
	total, err := discomodels.Discos.ListDisco(h.Ctx, &input, &output.Result)
	if err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	if total == 0 {
		output.Result = make([]cores.DiscoOutput, 0, 0)
	}

	output.Total = total

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *DiscosHandler) StatDiscoEthIncrease() (res interface{}) {
	var input cores.StatDiscoEthIncreaseInput
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

	var outputs []cores.StatDiscoEthIncreaseOutput
	if err := discomodels.Discos.StatDiscoEthIncrease(h.Ctx, &input, &outputs); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var totalEthCount int64
	if err := discomodels.Discos.StatDiscoEthTotal(h.Ctx, input.TimeFrom, &totalEthCount); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	for i := 0; i < len(outputs); i++ {
		if i == 0 {
			outputs[i].Count += totalEthCount
			continue
		}
		outputs[i].Count += outputs[i-1].Count
	}

	res = apires.With(&outputs)
	return
}

func (h *DiscosHandler) StatDiscoEthTotal() (res interface{}) {
	var totalEthCount int64
	if err := discomodels.Discos.StatDiscoEthTotal(h.Ctx, time.Now(), &totalEthCount); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&cores.StatDiscoEthTotalResult{
		Count: totalEthCount,
	})
	return
}

func (h *DiscosHandler) StatDiscoTotal() (res interface{}) {
	var output cores.StatDiscoTotalResult
	if err := discomodels.Discos.StatDiscoTotal(h.Ctx, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	if output.IcreaseCount > 0 {
		output.Rate = float64(output.IcreaseCount) / float64(output.Count)
	} else {
		output.Rate = 0
	}

	res = apires.With(&output)
	return
}

func (h *DiscosHandler) GetDiscoSwapState(startupId flake.ID) (res interface{}) {
	var output cores.DiscoSwapStateOutput
	if err := discomodels.Discos.GetDiscoSwapState(h.Ctx, &startupId, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output)
	return
}

func (h *DiscosHandler) DeleteStartupDisco(startupId flake.ID) (res interface{}) {
	var uid flake.ID
	h.Ctx.Find(&uid, "uid")

	if err := discomodels.Discos.DeleteDisco(h.Ctx, startupId, uid); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(http.StatusOK)
	return
}
