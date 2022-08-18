package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mises-id/sns-apigateway/app/apis/rest"

	pb "github.com/mises-id/sns-socialsvc/proto"
)

type (
	OpenseaSingleAssetInput struct {
		AssetContractAddress string `json:"asset_contract_address" query:"asset_contract_address"`
		TokenId              string `json:"token_id" query:"token_id"`
		AccountAddress       string `json:"account_address" query:"account_address"`
		IncludeOrders        string `json:"include_orders" query:"include_orders"`
		Network              string `json:"network" query:"network"`
	}
	OpenseaContractAssetInput struct {
		AssetContractAddress string `json:"asset_contract_address" query:"asset_contract_address"`
		Network              string `json:"network" query:"network"`
	}
	OpenseaSingleAssetOuput struct {
		Id       uint64 `json:"id"`
		ImageUrl string `json:"image_url"`
		Name     string `json:"name"`
	}
	ListOpenseaAssetInput struct {
		Owner   string `json:"owner" query:"owner"`
		Limit   uint64 `json:"limit" query:"limit"`
		Cursor  string `json:"cursor" query:"cursor"`
		Network string `json:"network" query:"network"`
	}
	ListOpenseaAssetOutput struct {
		Next     string                     `json:"next"`
		Previous string                     `json:"previous"`
		Data     []*OpenseaSingleAssetOuput `json:"data"`
	}
)

func GetOpenseaAsset(c echo.Context) error {
	params := &OpenseaSingleAssetInput{}
	if err := c.Bind(params); err != nil {
		return err
	}
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.GetOpenseaAsset(ctx, &pb.GetOpenseaAssetRequest{
		AssetContractAddress: params.AssetContractAddress,
		TokenId:              params.TokenId,
		AccountAddress:       params.AccountAddress,
		IncludeOrders:        params.IncludeOrders,
		Network:              params.Network,
	})
	if err != nil {
		return err
	}

	return c.String(200, svcresp.OpenseaAsset)
}
func GetOpenseaAssetContract(c echo.Context) error {
	params := &OpenseaSingleAssetInput{}
	if err := c.Bind(params); err != nil {
		return err
	}
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.GetOpenseaAssetContract(ctx, &pb.GetOpenseaAssetContractRequest{
		AssetContractAddress: params.AssetContractAddress,
		Network:              params.Network,
	})
	if err != nil {
		return err
	}

	return c.String(200, svcresp.OpenseaAsset)
}
func ListOpenseaAsset(c echo.Context) error {
	params := &ListOpenseaAssetInput{}
	if err := c.Bind(params); err != nil {
		return err
	}
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.ListOpenseaAsset(ctx, &pb.ListOpenseaAssetRequest{
		CurrentUid: GetCurrentUID(c),
		Owner:      params.Owner,
		Limit:      params.Limit,
		Cursor:     params.Cursor,
		Network:    params.Network,
	})
	if err != nil {
		return err
	}

	return c.String(200, svcresp.Assets)
}

func BuildOpenseaAssetResp(in *pb.OpenseaAsset) *OpenseaSingleAssetOuput {
	if in == nil {
		return nil
	}
	resp := &OpenseaSingleAssetOuput{
		Id:       in.Id,
		ImageUrl: in.ImageUrl,
		Name:     in.Name,
	}
	return resp
}
func BuildOpenseaAssetSliceResp(assets []*pb.OpenseaAsset) []*OpenseaSingleAssetOuput {
	res := make([]*OpenseaSingleAssetOuput, len(assets))
	for i, v := range assets {
		res[i] = BuildOpenseaAssetResp(v)
	}
	return res
}
