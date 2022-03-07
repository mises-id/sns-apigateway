package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
)

type GasPricesResp struct {
	SafeGasPrice    float64 `json:"safe_gasprice"`
	ProposeGasPrice float64 `json:"propose_gasprice"`
	FastGasPrice    float64 `json:"fast_gasprice"`
}

func GasPrices(c echo.Context) error {

	return rest.BuildSuccessResp(c, &GasPricesResp{
		SafeGasPrice:    0.001,
		ProposeGasPrice: 0.001,
		FastGasPrice:    0.001,
	})
}
