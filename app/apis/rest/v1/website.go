package v1

import (
	"encoding/json"
	"os"
	"path"

	"github.com/labstack/echo/v4"
	pb "github.com/mises-id/mises-websitesvc/proto"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	"github.com/mises-id/sns-apigateway/lib/codes"
)

type (
	WebsiteParams struct {
		WebSiteCategoryID string `json:"website_category_id" query:"website_category_id"`
		SubcategoryID     string `json:"subcategory_id" query:"subcategory_id"`
		ListNum           uint64 `json:"list_num" query:"list_num"`
		Keywords          string `json:"keywords" query:"keywords"`
		rest.PageParams
	}

	WebsiteResp struct {
		ID                string               `json:"id"`
		WebsiteCategoryID string               `json:"website_category_id"`
		SubcategoryID     string               `json:"subcategory_id"`
		Title             string               `json:"title"`
		Url               string               `json:"url"`
		Logo              string               `json:"logo"`
		Desc              string               `json:"desc"`
		WebSiteCategory   *WebsiteCategoryResp `json:"website_category"`
		Subcategory       *WebsiteCategoryResp `json:"subcategory"`
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
		SubcategoryId:     params.SubcategoryID,
		Paginator: &pb.Page{
			PageNum:  uint64(params.PageParams.PageNum),
			PageSize: uint64(params.PageParams.PageSize),
		},
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessRespWithWebsitePage(c, BuildWebsiteSliceResp(svcresp.Data), svcresp.Paginator)
}
func PageExtensions(c echo.Context) error {
	params := &WebsiteParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}
	grpcsvc, ctx, err := rest.GrpcWebsiteService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.WebsitePage(ctx, &pb.WebsitePageRequest{
		Type:              "extensions",
		Keywords:          params.Keywords,
		WebsiteCategoryId: params.WebSiteCategoryID,
		SubcategoryId:     params.SubcategoryID,
		Paginator: &pb.Page{
			PageNum:  uint64(params.PageParams.PageNum),
			PageSize: uint64(params.PageParams.PageSize),
		},
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessRespWithWebsitePage(c, BuildWebsiteSliceResp(svcresp.Data), svcresp.Paginator)
}

func CreateRecommendJson(c echo.Context) error {
	params := &WebsiteParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}
	grpcsvc, ctx, err := rest.GrpcWebsiteService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.WebsiteRecommend(ctx, &pb.WebsiteRecommendRequest{ListNum: params.ListNum})
	if err != nil {
		return err
	}
	data := BuildWebsiteSliceResp(svcresp.Data)
	filePath := "./assets/website"
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	filename := "recommend.json"
	localfile := path.Join(filePath, filename)
	filePtr, err := os.Create(localfile)
	if err != nil {
		return err
	}
	defer filePtr.Close()
	encoder := json.NewEncoder(filePtr)
	return encoder.Encode(data)
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
		SubcategoryID:     data.SubcategoryId,
		Title:             data.Title,
		Url:               data.Url,
		Logo:              data.Logo,
		Desc:              data.Desc,
	}
	if data.WebsiteCategory != nil {
		resp.WebSiteCategory = BuildWebsiteCategoryResp(data.WebsiteCategory)
	}
	if data.Subcategory != nil {
		resp.Subcategory = BuildWebsiteCategoryResp(data.Subcategory)
	}
	return resp
}
