package v1

import (
    "github.com/labstack/echo/v4"
    swapsvc "github.com/mises-id/mises-swapsvc/proto"
    "github.com/mises-id/sns-apigateway/app/apis/rest"
    "github.com/mises-id/sns-apigateway/lib/codes"
    "time"
)

type BridgeGetCurrenciesItem struct {
    Ticker             string   `json:"symbol"`
    Name               string   `json:"name"`
    ContractAddress    string   `json:"address"`
    Image              string   `json:"logo_uri"`
    Enabled            bool     `json:"bridgeEnabled"`
    EnabledFrom        bool     `json:"bridgeEnabledFrom"`
    EnabledTo          bool     `json:"bridgeEnabledTo"`
    FixRateEnabled     bool     `json:"bridgeFixRateEnabled"`
    PayinConfirmations uint64   `json:"bridgePayinConfirmations"`
    ExtraIdName        string   `json:"bridgeExtraIdName"`
    FixedTime          uint64   `json:"bridgeFixedTime"`
    Protocol           string   `json:"bridgeProtocol"`
    Blockchain         string   `json:"bridgeBlockchain"`
    Notifications      *swapsvc.BridgeNotifications `json:"bridgeNotifications"`
}

type NewBridgeHistoryItem struct {
    Id           string `json:"id"`
    Status       string `json:"status"`
    CreatedAt    string `json:"createdAt"`
    CurrencyFrom string `json:"currencyFrom"`
    AmountFrom   string `json:"amountFrom"`
    CurrencyTo   string `json:"currencyTo"`
    AmountTo     string `json:"amountTo"`
    TrackUrl     string `json:"trackUrl"`
}

func buildBridgeGetCurrenciesSuccessResp(c echo.Context, data []*swapsvc.BridgeCurrencyInfo) error {
    ret := make([]*BridgeGetCurrenciesItem, 0, len(data))
    for _, v := range data {
        ret = append(ret, &BridgeGetCurrenciesItem{
            Ticker: v.Ticker,
            Name: v.Fullname,
            ContractAddress: v.ContractAddress,
            Image: v.Image,
            Enabled: v.Enabled,
            EnabledFrom: v.EnabledFrom,
            EnabledTo: v.EnabledTo,
            FixRateEnabled: v.FixRateEnabled,
            PayinConfirmations: v.PayinConfirmations,
            ExtraIdName: v.ExtraIdName,
            FixedTime: v.FixedTime,
            Protocol: v.Protocol,
            Blockchain: v.Blockchain,
            Notifications: v.Notifications,
        })
    }
    return rest.BuildSuccessResp(c, ret)
}

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
    return buildBridgeGetCurrenciesSuccessResp(c, ret.Data)
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
    params.EthAddress = GetCurrentEthAddress(c)
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
    params.EthAddress = GetCurrentEthAddress(c)
    ret, err := grpcsvc.BridgeCreateFixTransaction(ctx, params)
    if err != nil {
        return err
    }
    return rest.BuildSuccessResp(c, ret.Data)
}

func buildBridgeHistoryListSuccessResp(c echo.Context, data []*swapsvc.BridgeHistoryItem) error {
    ret := make([]*NewBridgeHistoryItem, 0, len(data))
    for _, v := range data {
        ret = append(ret, &NewBridgeHistoryItem{
            Id: v.Id,
            Status: v.Status,
            CreatedAt: bridgeHistoryFormatTime(v.CreatedAt),
            CurrencyFrom: v.CurrencyFrom,
            CurrencyTo: v.CurrencyTo,
            AmountFrom: v.AmountFrom,
            AmountTo: v.AmountTo,
            TrackUrl: v.TrackUrl,
        })
    }
    return rest.BuildSuccessResp(c, ret)
}

func bridgeHistoryFormatTime(unixTimestamp int64) string {
    t := time.Unix(unixTimestamp, 0)
    t = t.UTC()
    return t.Format("02 Jan 2006, 15:04:05")
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
    params.EthAddress = GetCurrentEthAddress(c)
    ret, err := grpcsvc.BridgeHistoryList(ctx, params)
    if err != nil {
        return err
    }
    return buildBridgeHistoryListSuccessResp(c, ret.Data)
}
