package v1

import (
    "github.com/labstack/echo/v4"
    swapsvc "github.com/mises-id/mises-swapsvc/proto"
    "github.com/mises-id/sns-apigateway/app/apis/rest"
    "github.com/mises-id/sns-apigateway/lib/codes"
)

func BridgeGetCurrencies(c echo.Context) (err error) {
    grpcsvc, ctx, err := rest.GrpcSwapService()
    if err != nil {
        return err
    }
    ret, err := grpcsvc.BridgeGetCurrencies(ctx, &swapsvc.BridgeGetCurrenciesRequest{
    })
    if err != nil {
        return err
    }
    return rest.BuildSuccessResp(c, ret.Data)
}

func BridgeGetPairsParams(c echo.Context) (err error) {
    params := &swapsvc.BridgeGetPairsParamsRequest{}
    if err := c.Bind(params); err != nil {
        return codes.ErrInvalidArgument.New("invalid query params")
    }
    grpcsvc, ctx, err := rest.GrpcSwapService()
    if err != nil {
        return err
    }
    ret, err := grpcsvc.BridgeGetPairsParams(ctx, params)
    if err != nil {
        return err
    }
    return rest.BuildSuccessResp(c, ret.Data)
}

func BridgeGetExchangeAmount(c echo.Context) (err error) {
    params := &swapsvc.BridgeGetExchangeAmountRequest{}
    if err := c.Bind(params); err != nil {
        return codes.ErrInvalidArgument.New("invalid query params")
    }
    grpcsvc, ctx, err := rest.GrpcSwapService()
    if err != nil {
        return err
    }
    ret, err := grpcsvc.BridgeGetExchangeAmount(ctx, params)
    if err != nil {
        return err
    }
    return rest.BuildSuccessResp(c, ret.Data)
}

func BridgeCreateTransaction(c echo.Context) (err error) {
    params := &swapsvc.BridgeCreateTransactionRequest{}
    if err := c.Bind(params); err != nil {
        return codes.ErrInvalidArgument.New("invalid query params")
    }
    grpcsvc, ctx, err := rest.GrpcSwapService()
    if err != nil {
        return err
    }
    ret, err := grpcsvc.BridgeCreateTransaction(ctx, params)
    if err != nil {
        return err
    }
    return rest.BuildSuccessResp(c, ret.Data)
}

func BridgeGetTransactionInfo(c echo.Context) (err error) {
    params := &swapsvc.BridgeGetTransactionInfoRequest{}
    if err := c.Bind(params); err != nil {
        return codes.ErrInvalidArgument.New("invalid query params")
    }
    grpcsvc, ctx, err := rest.GrpcSwapService()
    if err != nil {
        return err
    }
    ret, err := grpcsvc.BridgeGetTransactionInfo(ctx, params)
    if err != nil {
        return err
    }
    return rest.BuildSuccessResp(c, ret.Data)
}

func BridgeGetTransactionStatus(c echo.Context) (err error) {
    params := &swapsvc.BridgeGetTransactionStatusRequest{}
    if err := c.Bind(params); err != nil {
        return codes.ErrInvalidArgument.New("invalid query params")
    }
    grpcsvc, ctx, err := rest.GrpcSwapService()
    if err != nil {
        return err
    }
    ret, err := grpcsvc.BridgeGetTransactionStatus(ctx, params)
    if err != nil {
        return err
    }
    return rest.BuildSuccessResp(c, ret.Data)
}

func BridgeValidateAddress(c echo.Context) (err error) {
    params := &swapsvc.BridgeValidateAddressRequest{}
    if err := c.Bind(params); err != nil {
        return codes.ErrInvalidArgument.New("invalid query params")
    }
    grpcsvc, ctx, err := rest.GrpcSwapService()
    if err != nil {
        return err
    }
    ret, err := grpcsvc.BridgeValidateAddress(ctx, params)
    if err != nil {
        return err
    }
    return rest.BuildSuccessResp(c, ret.Data)
}

func BridgeGetFixRateForAmount(c echo.Context) (err error) {
    params := &swapsvc.BridgeGetFixRateForAmountRequest{}
    if err := c.Bind(params); err != nil {
        return codes.ErrInvalidArgument.New("invalid query params")
    }
    grpcsvc, ctx, err := rest.GrpcSwapService()
    if err != nil {
        return err
    }
    ret, err := grpcsvc.BridgeGetFixRateForAmount(ctx, params)
    if err != nil {
        return err
    }
    return rest.BuildSuccessResp(c, ret.Data)
}

func BridgeCreateFixTransaction(c echo.Context) (err error) {
    params := &swapsvc.BridgeCreateFixTransactionRequest{}
    if err := c.Bind(params); err != nil {
        return codes.ErrInvalidArgument.New("invalid query params")
    }
    grpcsvc, ctx, err := rest.GrpcSwapService()
    if err != nil {
        return err
    }
    ret, err := grpcsvc.BridgeCreateFixTransaction(ctx, params)
    if err != nil {
        return err
    }
    return rest.BuildSuccessResp(c, ret.Data)
}

func BridgeHistoryList(c echo.Context) (err error) {
    params := &swapsvc.BridgeHistoryListRequest{}
    if err := c.Bind(params); err != nil {
        return codes.ErrInvalidArgument.New("invalid query params")
    }
    grpcsvc, ctx, err := rest.GrpcSwapService()
    if err != nil {
        return err
    }
    ret, err := grpcsvc.BridgeHistoryList(ctx, params)
    if err != nil {
        return err
    }
    return rest.BuildSuccessResp(c, ret.Data)
}
