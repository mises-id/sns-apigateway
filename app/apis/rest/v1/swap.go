package v1

import (
	"strconv"

	"github.com/labstack/echo/v4"
	pb "github.com/mises-id/mises-swapsvc/proto"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	"github.com/mises-id/sns-apigateway/lib/codes"
)

type (
	//requests
	//----------------------------------------------------------------
	SwapPublicRequest struct {
		RequestID string `json:"request_id" query:"request_id"`
		ChainID   uint64 `json:"chain_id" query:"chain_id"`
	}
	SwapOrderRequest struct {
		FromAddress string `json:"from_address" query:"from_address"`
		TxHash      string `json:"txHash" query:"txHash"`
		rest.PageParams
		SwapPublicRequest
	}
	TokenRequest struct {
		SwapPublicRequest
	}
	SwapApproveRequest struct {
		SwapPublicRequest
		TokenAddress      string `json:"token_address" query:"token_address"`
		WalletAddress     string `json:"wallet_address" query:"wallet_address"`
		AggregatorAddress string `json:"aggregator_address" query:"aggregator_address"`
		Amount            string `json:"amount" query:"amount"`
	}
	SwapTradeRequest struct {
		SwapPublicRequest
		Amount            string  `json:"amount" query:"amount"`
		FromTokenAddress  string  `json:"from_token_address" query:"from_token_address"`
		ToTokenAddress    string  `json:"to_token_address" query:"to_token_address"`
		FromAddress       string  `json:"from_address" query:"from_address"`
		DestReceiver      string  `json:"dest_receiver" query:"dest_receiver"`
		Slippage          float32 `json:"slippage" query:"slippage"`
		AggregatorAddress string  `json:"aggregator_address" query:"aggregator_address"`
	}
	SwapQuoteRequest struct {
		SwapPublicRequest
		Amount           string `json:"amount" query:"amount"`
		FromTokenAddress string `json:"from_token_address" query:"from_token_address"`
		ToTokenAddress   string `json:"to_token_address" query:"to_token_address"`
	}
	//response
	//----------------------------------------------------------------
	SwapProvider struct {
		Key  string `json:"key"`
		Name string `json:"name"`
		Logo string `json:"logo"`
	}
	Token struct {
		Address  string `json:"address"`
		Decimals int32  `json:"decimals"`
		LogoUri  string `json:"logo_uri"`
		Name     string `json:"name"`
		Symbol   string `json:"symbol"`
		Price    string `json:"-"`
		Value    string `json:"-"`
		ChainID  uint64 `json:"chain_id"`
	}
	Transaction struct {
		Hash        string `json:"hash"`
		Gas         string `json:"gas"`
		GasUsed     string `json:"gas_used"`
		GasPrice    string `json:"gas_price"`
		Nonce       string `json:"nonce"`
		BlockNumber int64  `json:"block_number"`
	}

	SwapOrderResponse struct {
		Id              string        `json:"id"`
		ChainID         uint64        `json:"chain_id"`
		FromAddress     string        `json:"from_address"`
		DestReceiver    string        `json:"dest_receiver"`
		ReceiptStatus   int32         `json:"receipt_status"`
		FromToken       *Token        `json:"from_token"`
		ToToken         *Token        `json:"to_token"`
		FromTokenAmount string        `json:"from_token_amount"`
		ToTokenAmount   string        `json:"to_token_amount"`
		Provider        *SwapProvider `json:"provider"`
		ContractAddress string        `json:"contract_address"`
		ReferrerAddress string        `json:"referrer_address"`
		Fee             float32       `json:"fee"`
		BlockAt         int64         `json:"block_at"`
		Tx              *Transaction  `json:"tx"`
	}
	Trade struct {
		Data     string `json:"data"`
		From     string `json:"from"`
		To       string `json:"to"`
		GasPrice string `json:"gas_price"`
		GasLimit string `json:"gas_limit"`
		Value    string `json:"value"`
	}
	Aggregator struct {
		Type            string `json:"type"`
		Name            string `json:"name"`
		Logo            string `json:"logo"`
		ContractAddress string `json:"contract_address"`
	}
	SwapTradeInfo struct {
		Aggregator       *Aggregator `json:"aggregator"`
		FromTokenAddress string      `json:"from_token_address"`
		ToTokenAddress   string      `json:"to_token_address"`
		FromTokenAmount  string      `json:"from_token_amount"`
		ToTokenAmount    string      `json:"to_token_amount"`
		Trade            *Trade      `json:"trade"`
		Error            string      `json:"error"`
		FetchTime        int64       `json:"fetch_time"`
		Fee              float32     `json:"fee"`
	}
	ApproveAllowanceResp struct {
		Allowance string `json:"allowance"`
	}
	ApproveSwapTransactionResponse struct {
		Data     string `json:"data"`
		To       string `json:"to"`
		GasPrice string `json:"gas_price"`
		GasLimit string `json:"gas_limit"`
		Value    string `json:"value"`
	}
	SwapQuoteInfo struct {
		Aggregator       *Aggregator `json:"aggregator"`
		FromTokenAddress string      `json:"from_token_address"`
		ToTokenAddress   string      `json:"to_token_address"`
		FromTokenAmount  string      `json:"from_token_amount"`
		ToTokenAmount    string      `json:"to_token_amount"`
		EstimateGasFee   string      `json:"estimate_gas_fee"`
		Error            string      `json:"error"`
		FetchTime        int64       `json:"fetch_time"`
		Fee              float32     `json:"fee"`
		ComparePercent   float32     `json:"compare_percent"`
	}
	SwapQuoteResponse struct {
		BestQuote *SwapQuoteInfo   `json:"best_quote"`
		Error     string           `json:"error"`
		AllQuote  []*SwapQuoteInfo `json:"all_quote"`
	}
	WalletsAndTokensRequest struct {
		ChainID uint64
		Wallets []string
		Tokens  []string
	}

	WalletsAndTokensResponse *map[string]map[string]string
)

func GetChainIDParam(c echo.Context) (uint64, error) {
	chain_idParam := c.Param("chain_id")
	chain_id, err := strconv.ParseUint(chain_idParam, 10, 64)
	if err != nil {
		return 0, codes.ErrInvalidArgument.Newf("invalid chain_id %s", chain_idParam)
	}
	return chain_id, nil
}

func PageSwapOrder(c echo.Context) error {
	params := &SwapOrderRequest{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}
	grpcsvc, ctx, err := rest.GrpcSwapService()
	if err != nil {
		return err
	}
	from_address := c.Param("from_address")
	svcresp, err := grpcsvc.SwapOrderPage(ctx, &pb.SwapOrderPageRequest{
		ChainID:     params.ChainID,
		FromAddress: from_address,
		Paginator: &pb.Page{
			PageNum:  uint64(params.PageParams.PageNum),
			PageSize: uint64(params.PageParams.PageSize),
		},
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessRespWithSwapPage(c, params.RequestID, BuildSwapOrderSliceResp(svcresp.Data), svcresp.Paginator)
}

func FindSwapOrder(c echo.Context) error {
	params := &SwapOrderRequest{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}
	grpcsvc, ctx, err := rest.GrpcSwapService()
	if err != nil {
		return err
	}
	from_address := c.Param("from_address")
	txHash := c.Param("tx_hash")
	svcresp, err := grpcsvc.FindSwapOrder(ctx, &pb.FindSwapOrderRequest{
		ChainID:     params.ChainID,
		FromAddress: from_address,
		TxHash:      txHash,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessRespWithRequestID(c, params.RequestID, BuildSwapOrderResponse(svcresp.Data))
}

func ListTokens(c echo.Context) error {
	params := &TokenRequest{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}
	grpcsvc, ctx, err := rest.GrpcSwapService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.ListSwapToken(ctx, &pb.ListSwapTokenRequest{
		ChainID: params.ChainID,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessRespWithRequestID(c, params.RequestID, BuildSwapTokenSlice(svcresp.Data))
}

// approve allowance
func GetSwapApproveAllowance(c echo.Context) error {
	params := &SwapApproveRequest{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}
	grpcsvc, ctx, err := rest.GrpcSwapService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.GetSwapApproveAllowance(ctx, &pb.GetSwapApproveAllowanceRequest{
		ChainID:           params.ChainID,
		TokenAddress:      params.TokenAddress,
		WalletAddress:     params.WalletAddress,
		AggregatorAddress: params.AggregatorAddress,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessRespWithRequestID(c, params.RequestID, buildApproveAllowance(svcresp))
}

// approve transaction
func ApproveSwapTransaction(c echo.Context) error {
	params := &SwapApproveRequest{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}
	grpcsvc, ctx, err := rest.GrpcSwapService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.ApproveSwapTransaction(ctx, &pb.ApproveSwapTransactionRequest{
		ChainID:           params.ChainID,
		TokenAddress:      params.TokenAddress,
		Amount:            params.Amount,
		AggregatorAddress: params.AggregatorAddress,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessRespWithRequestID(c, params.RequestID, buildApproveTransaction(svcresp))
}

func SwapHealth(c echo.Context) error {
	grpcsvc, ctx, err := rest.GrpcSwapService()
	if err != nil {
		return err
	}
	_, err = grpcsvc.Health(ctx, &pb.HealthRequest{})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, nil)
}

func SwapTrade(c echo.Context) error {
	params := &SwapTradeRequest{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}
	grpcsvc, ctx, err := rest.GrpcSwapService()
	if err != nil {
		return err
	}
	ethAddress := GetCurrentEthAddress(c)
	svcresp, err := grpcsvc.SwapTrade(ctx, &pb.SwapTradeRequest{
		ChainID:           params.ChainID,
		FromTokenAddress:  params.FromTokenAddress,
		ToTokenAddress:    params.ToTokenAddress,
		FromAddress:       params.FromAddress,
		DestReceiver:      params.DestReceiver,
		Amount:            params.Amount,
		Slippage:          params.Slippage,
		AggregatorAddress: params.AggregatorAddress,
		Misesid:           ethAddress,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessRespWithRequestID(c, params.RequestID, buildSwapTradeInfo(svcresp.Data))
}
func SwapQuote(c echo.Context) error {
	params := &SwapQuoteRequest{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}
	grpcsvc, ctx, err := rest.GrpcSwapService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.SwapQuote(ctx, &pb.SwapQuoteRequest{
		ChainID:          params.ChainID,
		FromTokenAddress: params.FromTokenAddress,
		ToTokenAddress:   params.ToTokenAddress,
		Amount:           params.Amount,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessRespWithRequestID(c, params.RequestID, buildSwapQuoteResponse(svcresp))
}

func buildSwapQuoteResponse(data *pb.SwapQuoteResponse) *SwapQuoteResponse {
	if data == nil {
		return nil
	}
	resp := &SwapQuoteResponse{
		BestQuote: buildSwapQuoteInfo(data.BestQuote),
		AllQuote:  buildSwapQuoteSlice(data.AllQuote),
		Error:     data.Error,
	}
	return resp
}

func buildSwapQuoteSlice(data []*pb.SwapQuoteInfo) []*SwapQuoteInfo {
	resp := make([]*SwapQuoteInfo, len(data))
	for i, v := range data {
		resp[i] = buildSwapQuoteInfo(v)
	}
	return resp
}

func buildSwapQuoteInfo(data *pb.SwapQuoteInfo) *SwapQuoteInfo {
	if data == nil {
		return nil
	}
	resp := &SwapQuoteInfo{
		FromTokenAddress: data.FromTokenAddress,
		ToTokenAddress:   data.ToTokenAddress,
		FromTokenAmount:  data.FromTokenAmount,
		ToTokenAmount:    data.ToTokenAmount,
		FetchTime:        data.FetchTime,
		Fee:              data.Fee,
		EstimateGasFee:   data.EstimateGasFee,
		Error:            data.Error,
		ComparePercent:   data.ComparePercent,
	}
	if data.Aggregator != nil {
		resp.Aggregator = &Aggregator{
			Name:            data.Aggregator.Name,
			Logo:            data.Aggregator.Logo,
			Type:            data.Aggregator.Type,
			ContractAddress: data.Aggregator.ContractAddress,
		}
	}
	return resp
}
func buildSwapTradeSlice(data []*pb.SwapTradeInfo) []*SwapTradeInfo {
	resp := make([]*SwapTradeInfo, len(data))
	for i, v := range data {
		resp[i] = buildSwapTradeInfo(v)
	}
	return resp
}

func buildSwapTradeInfo(data *pb.SwapTradeInfo) *SwapTradeInfo {
	if data == nil {
		return nil
	}
	resp := &SwapTradeInfo{
		FromTokenAddress: data.FromTokenAddress,
		ToTokenAddress:   data.ToTokenAddress,
		FromTokenAmount:  data.FromTokenAmount,
		ToTokenAmount:    data.ToTokenAmount,
		FetchTime:        data.FetchTime,
		Fee:              data.Fee,
		Error:            data.Error,
	}
	if data.Aggregator != nil {
		resp.Aggregator = &Aggregator{
			Name:            data.Aggregator.Name,
			Logo:            data.Aggregator.Logo,
			Type:            data.Aggregator.Type,
			ContractAddress: data.Aggregator.ContractAddress,
		}
	}
	if data.Trade != nil {
		resp.Trade = &Trade{
			Data:     data.Trade.Data,
			From:     data.Trade.From,
			To:       data.Trade.To,
			GasPrice: data.Trade.GasPrice,
			GasLimit: data.Trade.GasLimit,
			Value:    data.Trade.Value,
		}
	}
	return resp
}

func buildApproveTransaction(data *pb.ApproveSwapTransactionResponse) *ApproveSwapTransactionResponse {
	if data == nil {
		return nil
	}
	resp := &ApproveSwapTransactionResponse{
		Data:     data.Data,
		To:       data.To,
		GasPrice: data.GasPrice,
		Value:    data.Value,
		GasLimit: data.GasLimit,
	}
	return resp
}

func buildApproveAllowance(data *pb.GetSwapApproveAllowanceResponse) *ApproveAllowanceResp {
	if data == nil {
		return nil
	}
	resp := &ApproveAllowanceResp{
		Allowance: data.Allowance,
	}
	return resp
}

func BuildSwapTokenSlice(data []*pb.Token) []*Token {
	resp := make([]*Token, len(data))
	for i, v := range data {
		resp[i] = NewToken(v)
	}
	return resp
}

func BuildSwapOrderSliceResp(data []*pb.SwapOrder) []*SwapOrderResponse {
	resp := make([]*SwapOrderResponse, len(data))
	for i, v := range data {
		resp[i] = BuildSwapOrderResponse(v)
	}
	return resp
}

func BuildSwapOrderResponse(data *pb.SwapOrder) *SwapOrderResponse {
	if data == nil {
		return nil
	}
	resp := &SwapOrderResponse{
		Id:              data.Id,
		ChainID:         data.ChainID,
		FromAddress:     data.FromAddress,
		DestReceiver:    data.DestReceiver,
		ReceiptStatus:   int32(data.ReceiptStatus),
		ContractAddress: data.ContractAddress,
		//TransactionFee: data.TransactionFee,
		BlockAt: data.BlockAt,
	}
	if data.Tx != nil {
		resp.Tx = NewTransaction(data.Tx)
	}
	if data.Provider != nil {
		resp.Provider = NewSwapProvider(data.Provider)
	}
	if data.FromToken != nil {
		resp.FromToken = NewToken(data.FromToken)
		resp.FromTokenAmount = data.FromToken.Value
	}
	if data.ToToken != nil {
		resp.ToToken = NewToken(data.ToToken)
		resp.ToTokenAmount = data.ToToken.Value
	}
	//resp.NativeToken = NewToken(data.NativeToken)
	return resp
}
func NewSwapProvider(data *pb.SwapProvider) *SwapProvider {
	if data == nil {
		return nil
	}
	resp := SwapProvider{
		Key:  data.Key,
		Name: data.Name,
		Logo: data.Logo,
	}
	return &resp
}
func NewTransaction(data *pb.Transaction) *Transaction {
	if data == nil {
		return nil
	}
	resp := Transaction{
		Hash:        data.Hash,
		Gas:         data.Gas,
		GasUsed:     data.GasUsed,
		GasPrice:    data.GasPrice,
		Nonce:       data.Nonce,
		BlockNumber: data.BlockNumber,
	}
	return &resp
}

func NewToken(data *pb.Token) *Token {
	if data == nil {
		return nil
	}
	resp := Token{
		Address:  data.Address,
		Value:    data.Value,
		LogoUri:  data.LogoUri,
		Decimals: data.Decimals,
		Symbol:   data.Symbol,
		ChainID:  data.ChainID,
	}
	resp.Name = data.Name
	return &resp
}

func WalletsAndTokens(c echo.Context) error {
	params := &WalletsAndTokensRequest{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}
	grpcsvc, ctx, err := rest.GrpcSwapService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.WalletsAndTokens(ctx, &pb.WalletsAndTokensRequest{
		ChainID: params.ChainID,
		Wallets: params.Wallets,
		Tokens:  params.Tokens,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, buildWalletsAndTokens(svcresp))
}

func buildWalletsAndTokens(in *pb.WalletsAndTokensResponse) WalletsAndTokensResponse {
	var defaultResponse WalletsAndTokensResponse = &map[string]map[string]string{}
	if in == nil {
		return defaultResponse
	}
	resp := map[string]map[string]string{}
	for k, v := range in.Data {
		resp[k] = v.Token
	}
	return &resp
}
