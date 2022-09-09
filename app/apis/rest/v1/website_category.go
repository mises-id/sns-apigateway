package v1

import (
	"github.com/labstack/echo/v4"
	pb "github.com/mises-id/mises-websitesvc/proto"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	"github.com/mises-id/sns-apigateway/lib/codes"
)

type (
	ListWebsiteCategoryParams struct{}
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

	svcresp, err := grpcsvc.WebsiteCategoryList(ctx, &pb.WebsiteCategoryListRequest{})
	if err != nil {
		return err
	}

	return rest.BuildSuccessResp(c, svcresp)
}
