package v1

import (
	"github.com/labstack/echo/v4"
	pb "github.com/mises-id/mises-websitesvc/proto"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	"github.com/mises-id/sns-apigateway/lib/codes"
)

type (
	ListWebsiteCategoryParams struct{}
	WebsiteCategoryResp       struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		ShorterName string `json:"shorter_name"`
		Desc        string `json:"desc"`
		TypeString  string `json:"type_string"`
	}
)

func ListWebsiteCategory(c echo.Context) error {

	params := &ListWebsiteCategoryParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}

	grpcsvc, ctx, err := rest.GrpcWebsiteService()
	if err != nil {
		return err
	}

	svcresp, err := grpcsvc.WebsiteCategoryList(ctx, &pb.WebsiteCategoryListRequest{
		Type: "web3",
	})
	if err != nil {
		return err
	}

	return rest.BuildSuccessResp(c, BuildWebsiteCategorySliceResp(svcresp.Data))
}
func ListExtensionsCategory(c echo.Context) error {

	params := &ListWebsiteCategoryParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}

	grpcsvc, ctx, err := rest.GrpcWebsiteService()
	if err != nil {
		return err
	}

	svcresp, err := grpcsvc.WebsiteCategoryList(ctx, &pb.WebsiteCategoryListRequest{
		Type: "extensions",
	})
	if err != nil {
		return err
	}

	return rest.BuildSuccessResp(c, BuildWebsiteCategorySliceResp(svcresp.Data))
}

func BuildWebsiteCategorySliceResp(data []*pb.WebsiteCategory) []*WebsiteCategoryResp {
	resp := make([]*WebsiteCategoryResp, len(data))
	for i, v := range data {
		resp[i] = BuildWebsiteCategoryResp(v)
	}
	return resp
}

func BuildWebsiteCategoryResp(data *pb.WebsiteCategory) *WebsiteCategoryResp {
	if data == nil {
		return nil
	}
	resp := &WebsiteCategoryResp{
		ID:          data.Id,
		Name:        data.Name,
		ShorterName: data.ShorterName,
		TypeString:  data.TypeString,
		Desc:        data.Desc,
	}
	return resp
}
