package v1

import (
	"github.com/labstack/echo/v4"
	pb "github.com/mises-id/mises-websitesvc/proto"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	"github.com/mises-id/sns-apigateway/lib/codes"
)

type (
	PhishingSiteParams struct {
		DomainName string `json:"domain" query:"domain"`
		Logo string `json:"logo" query:"logo"`
		Content string `json:"content"`
	}
	PhishingSiteResp struct {
		DomainName string `json:"domain"`
		Level string `json:"level"`
		SuggestedUrl     string `json:"suggested_url"`
		HtmlBodyFuzzyHash     string `json:"html_body_fuzzy_hash"`
		LogoPhash string `json:"logo_phash"`
		TitleKeyword string `json:"title_keyword"`
		Tag string `json:"tag"`
		HtmlTextSimhash string `json:"html_text_simhash"`
	}
	VerifyContractParams struct {
		Address string `json:"address" query:"address"`
		DomainName string `json:"domain" query:"domain"`
	}
	VerifyContractResp struct {
		Address         string `json:"address" bson:"address"`
		TrustPercentage int    `json:"trust_percentage" bson:"trust_percentage"`
		Level string    `json:"level" bson:"level"`
	}
)

func PhishingCheck(c echo.Context) error {
	params := &PhishingSiteParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}
	grpcsvc, ctx, err := rest.GrpcWebsiteService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.PhishingCheck(ctx, &pb.PhishingCheckRequest{
		DomainName: params.DomainName,
		Logo:params.Logo,
		Content:params.Content,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, BuildPhishingSiteResp(svcresp))
}

func BuildPhishingSiteResp(data *pb.PhishingCheckResponse) *PhishingSiteResp {
	if data == nil {
		return nil
	}
	resp := &PhishingSiteResp{
		DomainName: data.DomainName,
		Level: data.Level,
		SuggestedUrl: data.SuggestedUrl,
		HtmlBodyFuzzyHash: data.HtmlBodyFuzzyHash,
		LogoPhash: data.LogoPhash,
		TitleKeyword: data.TitleKeyword,
		Tag: data.Tag,
		HtmlTextSimhash: data.HtmlTextSimhash,
	}
	return resp
}

func VerifyContract(c echo.Context) error {
	params := &VerifyContractParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}
	grpcsvc, ctx, err := rest.GrpcWebsiteService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.VerifyContract(ctx, &pb.VerifyContractRequest{
		Address: params.Address,
		DomainName: params.DomainName,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, BuildVerifyContractResp(svcresp))
}

func BuildVerifyContractResp(data *pb.VerifyContractResponse) *VerifyContractResp {
	if data == nil {
		return nil
	}
	resp := &VerifyContractResp{
		Address: data.Address,
		TrustPercentage: int(data.TrustPercentage),
		Level: data.Level,
	}
	return resp
}
