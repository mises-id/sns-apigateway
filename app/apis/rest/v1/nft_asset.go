package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	"github.com/mises-id/sns-apigateway/lib/codes"
	pb "github.com/mises-id/sns-socialsvc/proto"
)

type (
	PageNftAssetParams struct {
		UID    uint64
		SortBy string `query:"sort_by" json:"sort_by"`
		Scene  string `json:"scene" query:"scene"`
		rest.PageQuickParams
	}
	PaymentToken struct {
		ID       int64  `json:"id" bson:"id"`
		Symbol   string `json:"symbol,omitempty" bson:"symbol"`
		Address  string `json:"address,omitempty" bson:"address"`
		ImageURL string `json:"image_url,omitempty" bson:"image_url"`
		Name     string `json:"name,omitempty" bson:"name"`
		Decimals int64  `json:"decimals" bson:"decimals"`
		ETHPrice string `json:"eth_price,omitempty" bson:"eth_price"`
		USDPrice string `json:"usd_price,omitempty" bson:"usd_price"`
	}
	Stats struct {
		OneDayVolume          float64 `json:"one_day_volume" bson:"one_day_volume"`
		OneDayChange          float64 `json:"one_day_change" bson:"one_day_change"`
		OneDaySales           float64 `json:"one_day_sales" bson:"one_day_sales"`
		OneDayAveragePrice    float64 `json:"one_day_average_price" bson:"one_day_average_price"`
		SevenDayVolume        float64 `json:"seven_day_volume" bson:"seven_day_volume"`
		SevenDayChange        float64 `json:"seven_day_change" bson:"seven_day_change"`
		SevenDaySales         float64 `json:"seven_day_sales" bson:"seven_day_sales"`
		SevenDayAveragePrice  float64 `json:"seven_day_average_price" bson:"seven_day_average_price"`
		ThirtyDayVolume       float64 `json:"thirty_day_volume" bson:"thirty_day_volume"`
		ThirtyDayChange       float64 `json:"thirty_day_change" bson:"thirty_day_change"`
		ThirtyDaySales        float64 `json:"thirty_day_sales" bson:"thirty_day_sales"`
		ThirtyDayAveragePrice float64 `json:"thirty_day_average_price" bson:"thirty_day_average_price"`
		TotalVolume           float64 `json:"total_volume" bson:"total_volume"`
		TotalSales            float64 `json:"total_sales" bson:"total_sales"`
		TotalSupply           float64 `json:"total_supply" bson:"total_supply"`
		Count                 float64 `json:"count" bson:"count"`
		NumOwners             int64   `json:"num_owners" bson:"num_owners"`
		AveragePrice          float64 `json:"average_price" bson:"average_price"`
		NumReports            int64   `json:"num_reports" bson:"num_reports"`
		MarketCap             float64 `json:"market_cap" bson:"market_cap"`
		FloorPrice            float64 `json:"floor_price" bson:"floor_price"`
	}

	NftCollection struct {
		Id            string          `json:"id"`
		Name          string          `json:"name"`
		Slug          string          `json:"slug"`
		PaymentTokens []*PaymentToken `json:"payment_tokens"`
		Stats         *Stats          `json:"stats"`
	}
	AssetContract struct {
		Address string `json:"address"`
	}

	NftAssetResp struct {
		Id                string           `json:"id"`
		ImageUrl          string           `json:"image_url"`
		ImagePreviewUrl   string           `json:"image_preview_url"`
		ImageThumbnailUrl string           `json:"image_thumbnail_url"`
		PermaLink         string           `json:"perma_link"`
		LikesCount        uint64           `json:"likes_count"`
		CommentsCount     uint64           `json:"comments_count"`
		TokenId           string           `json:"token_id"`
		Name              string           `json:"name"`
		Collection        *NftCollection   `json:"collection"`
		AssetContract     *AssetContract   `json:"asset_contract"`
		IsLiked           bool             `json:"is_liked"`
		User              *UserSummaryResp `json:"user"    `
	}
	LikeRequest struct {
		rest.PageQuickParams
	}
	LikeResp struct {
		ID   string           `json:"id"`
		User *UserSummaryResp `json:"user"`
	}
)

func ListNftAssetLike(c echo.Context) error {
	params := &LikeRequest{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.Newf("invalid query params")
	}

	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	paginator := &pb.PageQuick{
		NextId: params.PageQuickParams.NextID,
		Limit:  uint64(params.PageQuickParams.Limit),
	}
	svcresp, err := grpcsvc.ListLike(ctx, &pb.ListLikeUserRequest{
		TargerId:   c.Param("id"),
		CurrentUid: GetCurrentUID(c),
		Paginator:  paginator,
	})
	if err != nil {
		return err
	}

	return rest.BuildSuccessRespWithPagination(c, BuildLikeSliceResp(svcresp.Likes), svcresp.Paginator)
}

func BuildLikeSliceResp(likes []*pb.Like) []*LikeResp {
	resp := make([]*LikeResp, len(likes))
	for k, like := range likes {
		resp[k] = BuildLikeResp(like)
	}
	return resp
}

func BuildLikeResp(like *pb.Like) *LikeResp {
	if like == nil {
		return nil
	}
	resp := &LikeResp{
		ID: like.Id,
	}
	if like.User != nil {
		resp.User = BuildUserSummaryResp(like.User)
	}
	return resp
}

func GetNftAsset(c echo.Context) error {

	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.GetNftAsset(ctx, &pb.GetNftAssetRequest{
		CurrentUid: GetCurrentUID(c),
		NftAssetId: c.Param("id"),
	})
	if err != nil {
		return err
	}

	return rest.BuildSuccessResp(c, BuildNftAssetResp(svcresp.Asset))
}

func BuildNftAssetRespSlice(assets []*pb.NftAsset) []*NftAssetResp {
	resp := make([]*NftAssetResp, len(assets))

	for k, asset := range assets {
		resp[k] = BuildNftAssetResp(asset)
	}
	return resp
}

func BuildNftAssetResp(asset *pb.NftAsset) *NftAssetResp {
	if asset == nil {
		return nil
	}
	resp := &NftAssetResp{
		Id:                asset.Id,
		ImageUrl:          asset.ImageUrl,
		ImagePreviewUrl:   asset.ImagePreviewUrl,
		ImageThumbnailUrl: asset.ImageThumbnailUrl,
		TokenId:           asset.TokenId,
		PermaLink:         asset.PermaLink,
		LikesCount:        asset.LikesCount,
		CommentsCount:     asset.CommentsCount,
		Name:              asset.Name,
		IsLiked:           asset.IsLiked,
		User:              BuildUserSummaryResp(asset.User),
	}
	if asset.Collection != nil {
		resp.Collection = BuildNftCollectionResp(asset.Collection)
	}
	if asset.AssetContract != nil {
		resp.AssetContract = BuildAssetContractResp(asset.AssetContract)
	}

	return resp
}

func BuildAssetContractResp(in *pb.AssetContract) *AssetContract {
	resp := &AssetContract{
		Address: in.Address,
	}
	return resp
}

func BuildNftCollectionResp(in *pb.NftCollection) *NftCollection {
	resp := &NftCollection{
		Name: in.Name,
		Slug: in.Slug,
	}
	return resp
}

func LikeNftAsset(c echo.Context) error {
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}

	_, err = grpcsvc.LikeNftAsset(ctx, &pb.LikeNftAssetRequest{
		CurrentUid: GetCurrentUID(c),
		NftAssetId: c.Param("id"),
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, nil)
}

func UnlikeNftAsset(c echo.Context) error {
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}

	_, err = grpcsvc.UnlikeNftAsset(ctx, &pb.UnLikeNftAssetRequest{
		CurrentUid: GetCurrentUID(c),
		NftAssetId: c.Param("id"),
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, nil)
}
