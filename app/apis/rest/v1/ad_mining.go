package v1

import (
	"github.com/labstack/echo/v4"
	miningsvc "github.com/mises-id/mises-miningsvc/proto"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	"github.com/mises-id/sns-apigateway/lib/codes"
	"github.com/sirupsen/logrus"
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
type AdMiningLogRequest struct {
	AdType string `json:"ad_type" query:"ad_type"`
}

type EstimateAdBonusResponse struct {
	Point float32 `json:"point"`
	Msg   string  `json:"msg"`
}
type AdMiningUserResponse struct {
	EthAddress      string `json:"eth_address"`
	LimitPerDay     uint32 `json:"limit_per_day"`
	TodayBonusCount uint32 `json:"today_bonus_count"`
}

func AdMiningLog(c echo.Context) (err error) {

	params := &AdMiningLogRequest{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}
	ethAddress := GetCurrentEthAddress(c)
	grpcsvc, ctx, err := rest.GrpcMiningService()
	if err != nil {
		return err
	}
	user_agent := userAgent(c)
	_, err = grpcsvc.AdMiningLog(ctx, &miningsvc.AdMiningLogRequest{
		AdType:  params.AdType,
		Address: ethAddress,
		UserAgent: &miningsvc.UserAgent{
			Ua:       user_agent.ua,
			Ipaddr:   user_agent.ipaddr,
			Os:       user_agent.os,
			Browser:  user_agent.browser,
			Platform: user_agent.platform,
			DeviceId: user_agent.device_id,
		},
	})
	if err != nil {
		return err
	}

	return rest.BuildSuccessResp(c, nil)
}

func MyAdMining(c echo.Context) (err error) {

	ethAddress := GetCurrentEthAddress(c)
	grpcsvc, ctx, err := rest.GrpcMiningService()
	if err != nil {
		return err
	}
	resp, err := grpcsvc.FindAdMiningUser(ctx, &miningsvc.FindAdMiningUserRequest{
		EthAddress: ethAddress,
	})
	if err != nil {
		return err
	}

	return rest.BuildSuccessResp(c, builAdMiningUserResponse(resp))
}

func builAdMiningUserResponse(in *miningsvc.FindAdMiningUserResponse) *AdMiningUserResponse {

	if in == nil {
		return nil
	}
	resp := &AdMiningUserResponse{
		EthAddress:      in.EthAddress,
		LimitPerDay:     in.LimitPerDay,
		TodayBonusCount: in.TodayBonusCount,
	}

	return resp
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

// mintegral callback
func MintegralCallback(c echo.Context) error {

	urlStr := c.Request().URL.String()
	grpcsvc, ctx, err := rest.GrpcMiningService()
	if err != nil {
		logrus.Error("MintegralCallback:grpc:", err)
		return rest.Build403Resp(c, "internal error")
	}
	_, err = grpcsvc.AdMiningCallback(ctx, &miningsvc.AdMiningCallbackRequest{
		AdType:      "mintegral",
		CallbackUrl: urlStr,
	})
	if err != nil {
		logrus.Error("MintegralCallback:response:", err)
		return rest.Build403Resp(c, "process error")
	}

	return rest.BuildSuccessResp(c, nil)
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
