package v1

import (
	"github.com/labstack/echo/v4"
	pb "github.com/mises-id/mises-vpnsvc/proto"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	"github.com/mises-id/sns-apigateway/lib/codes"
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
	return rest.BuildSuccessResp(c, resp)
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
	return rest.BuildSuccessResp(c, resp)
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
	return rest.BuildSuccessResp(c, resp)
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
	return rest.BuildSuccessResp(c, resp)
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
	return rest.BuildSuccessResp(c, resp)
}
