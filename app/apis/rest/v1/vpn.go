package v1

import (
	"github.com/labstack/echo/v4"
	pb "github.com/mises-id/mises-vpnsvc/proto"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	"github.com/mises-id/sns-apigateway/lib/codes"
	"net/http"
)

func FetchOrderInfo(c echo.Context) error {
	grpcsvc, ctx, err := rest.GrpcVpnService()
	if err != nil {
		return err
	}
	params := &pb.FetchOrderInfoRequest{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}
	params.EthAddress = GetCurrentEthAddress(c)
	resp, err := grpcsvc.FetchOrderInfo(ctx, params)
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, resp.Data)
}

func FetchOrders(c echo.Context) error {
	grpcsvc, ctx, err := rest.GrpcVpnService()
	if err != nil {
		return err
	}
	params := &pb.FetchOrdersRequest{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}
	params.EthAddress = GetCurrentEthAddress(c)
	resp, err := grpcsvc.FetchOrders(ctx, params)
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, resp.Data)
}

func VpnInfo(c echo.Context) error {
	grpcsvc, ctx, err := rest.GrpcVpnService()
	if err != nil {
		return err
	}
	params := &pb.VpnInfoRequest{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}
	params.EthAddress = GetCurrentEthAddress(c)
	resp, err := grpcsvc.VpnInfo(ctx, params)
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, resp.Data)
}

func CreateOrder(c echo.Context) error {
	grpcsvc, ctx, err := rest.GrpcVpnService()
	if err != nil {
		return err
	}
	params := &pb.CreateOrderRequest{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}
	params.EthAddress = GetCurrentEthAddress(c)
	resp, err := grpcsvc.CreateOrder(ctx, params)
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, resp.Data)
}

func UpdateOrder(c echo.Context) error {
	grpcsvc, ctx, err := rest.GrpcVpnService()
	if err != nil {
		return err
	}
	params := &pb.UpdateOrderRequest{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}
	params.EthAddress = GetCurrentEthAddress(c)
	resp, err := grpcsvc.UpdateOrder(ctx, params)
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, resp.Data)
}

func GetServerList(c echo.Context) error {
	grpcsvc, ctx, err := rest.GrpcVpnService()
	if err != nil {
		return buildErrorRespForVpnClient(c, "service error")
	}
	params := &pb.GetServerListRequest{}
	params.EthAddress = GetCurrentEthAddress(c)
	resp, err := grpcsvc.GetServerList(ctx, params)
	if err != nil {
		return buildErrorRespForVpnClient(c, "get server list error")
	}
	return c.JSON(http.StatusOK, echo.Map{
		"code": 0,
		"msg": "success",
		"servers": resp.Data.Servers,
	})
}

func GetServerLink(c echo.Context) error {
	grpcsvc, ctx, err := rest.GrpcVpnService()
	if err != nil {
		return buildErrorRespForVpnClient(c, "service error")
	}
	params := &pb.GetServerLinkRequest{}
	if err := c.Bind(params); err != nil {
		return buildErrorRespForVpnClient(c, "params error")
	}
	params.EthAddress = GetCurrentEthAddress(c)
	resp, err := grpcsvc.GetServerLink(ctx, params)
	if err != nil {
		return buildErrorRespForVpnClient(c, "get vpn link error")
	}
	return c.JSON(http.StatusOK, echo.Map{
		"code": 0,
		"msg": "success",
		"config": resp.Data,
	})
}

func buildErrorRespForVpnClient(c echo.Context, errMsg string) error {
	return c.JSON(http.StatusOK, echo.Map{
		"code": -1,
		"msg": errMsg,
	})
}
