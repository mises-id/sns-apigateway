package v1

import (
	"github.com/labstack/echo/v4"
	miningsvc "github.com/mises-id/mises-miningsvc/proto"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	"github.com/mises-id/sns-apigateway/lib/codes"
)

type AdMiningCallbackResponse struct {
	Point float32 `json:"point"`
	Msg   string  `json:"msg"`
}

type EstimateAdBonusRequest struct {
	AdType  string `json:"ad_type" query:"ad_type"`
	TxID    string `json:"tx_id" query:"tx_id"`
	Address string `json:"address" query:"address"`
}

type EstimateAdBonusResponse struct {
	Point float32 `json:"point"`
	Msg   string  `json:"msg"`
}

func EstimateAdBonus(c echo.Context) error {

	params := &EstimateAdBonusRequest{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}
	grpcsvc, ctx, err := rest.GrpcMiningService()
	if err != nil {
		return err
	}
	resp, err := grpcsvc.EstimateAdBonus(ctx, &miningsvc.EstimateAdBonusRequest{
		AdType:  params.AdType,
		Address: params.Address,
		TxId:    params.TxID,
	})
	if err != nil {
		return err
	}

	return rest.BuildSuccessResp(c, buildEstimateAdBonusResponse(resp))
}

func buildEstimateAdBonusResponse(in *miningsvc.EstimateAdBonusResponse) *EstimateAdBonusResponse {

	if in == nil {
		return nil
	}
	resp := &EstimateAdBonusResponse{
		Point: in.Point,
		Msg:   in.Msg,
	}

	return resp
}

// ADmob ssv
func ADMobSSV(c echo.Context) error {

	urlStr := c.Request().URL.String()
	grpcsvc, ctx, err := rest.GrpcMiningService()
	if err != nil {
		return rest.BuildSuccessResp(c, buildADMobSSVResponseOnError(err))
	}
	resp, err := grpcsvc.AdMiningCallback(ctx, &miningsvc.AdMiningCallbackRequest{
		AdType:      "admob",
		CallbackUrl: urlStr,
	})
	if err != nil {
		return rest.BuildSuccessResp(c, buildADMobSSVResponseOnError(err))
	}

	return rest.BuildSuccessResp(c, buildADMobSSVResponse(resp))
}

func buildADMobSSVResponseOnError(err error) *AdMiningCallbackResponse {
	resp := &AdMiningCallbackResponse{
		Point: 0,
		Msg:   err.Error(),
	}

	return resp
}
func buildADMobSSVResponse(in *miningsvc.AdMiningCallbackResponse) *AdMiningCallbackResponse {

	if in == nil {
		return nil
	}
	resp := &AdMiningCallbackResponse{
		Point: in.Point,
	}

	return resp
}
