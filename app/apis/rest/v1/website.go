package v1

import (
	"github.com/labstack/echo/v4"
	pb "github.com/mises-id/mises-websitesvc/proto"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	"github.com/mises-id/sns-apigateway/lib/codes"
)

type (
	WebsiteParams struct {
		WebSiteCategoryID string `json:"website_category_id" query:"website_category_id"`
		Keywords          string `json:"keywords" query:"keywords"`
		rest.PageQuickParams
	}

	WebsiteResp struct {
		ID                string `json:"id"`
		WebsiteCategoryID string `json:"website_category_id"`
		Title             string `json:"title"`
		Url               string `json:"url"`
		Logo              string `json:"logo"`
		Desc              string `json:"desc"`
	}
)

func PageWebsite(c echo.Context) error {

	params := &WebsiteParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}
	grpcsvc, ctx, err := rest.GrpcWebsiteService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.WebsitePage(ctx, &pb.WebsitePageRequest{
		Type:              "web3",
		Keywords:          params.Keywords,
		WebsiteCategoryId: params.WebSiteCategoryID,
		Paginator: &pb.PageQuick{
			NextId: params.PageQuickParams.NextID,
			Limit:  uint64(params.PageQuickParams.Limit),
		},
	})
	if err != nil {
		return err
	}

	return rest.BuildSuccessRespWebsiteWithPagination(c, BuildWebsiteSliceResp(svcresp.Data), svcresp.Paginator)
}

func BuildWebsiteSliceResp(data []*pb.Website) []*WebsiteResp {
	resp := make([]*WebsiteResp, len(data))
	for i, v := range data {
		resp[i] = BuildWebsiteResp(v)
	}
	return resp
}

func BuildWebsiteResp(data *pb.Website) *WebsiteResp {
	if data == nil {
		return nil
	}
	resp := &WebsiteResp{
		ID:                data.Id,
		WebsiteCategoryID: data.WebsiteCategoryId,
		Title:             data.Title,
		Url:               data.Url,
		Logo:              data.Logo,
		Desc:              data.Desc,
	}
	return resp
}
