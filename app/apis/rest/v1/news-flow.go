package v1

import (
	"time"

	"github.com/labstack/echo/v4"
	pb "github.com/mises-id/mises-news-flow/pkg/proto/apiserver/v1"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	"github.com/mises-id/sns-apigateway/lib/codes"
)

type ListNewsParams struct {
	PageIndex int32 `json:"page_index" query:"page_index"`
}

func ListNews(c echo.Context) error {
	params := &ListNewsParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.Newf("invalid query params")
	}

	grpcsvc, ctx, err := rest.GrpcNewsFlowService()
	if err != nil {
		return err
	}

	resp, err := grpcsvc.FindNewsInPage(
		ctx,
		&pb.FindNewsInPageParams{
			PageIndex: params.PageIndex,
		},
	)
	if err != nil {
		return err
	}

	return rest.BuildSuccessResp(c, NewListNewsRespFromPB(resp))
}

type ListNewsResp struct {
	NewsArray     []*News `json:"news_array"`
	NextPageIndex int32   `json:"next_page_index"`
}

func NewListNewsRespFromPB(pbResp *pb.FindNewsInPageResult) *ListNewsResp {
	newsArray := make([]*News, 0)
	for _, pbNews := range pbResp.NewsArray {
		if news := NewNewsFromPB(pbNews); news != nil {
			newsArray = append(newsArray, news)
		}
	}
	return &ListNewsResp{
		NewsArray:     newsArray,
		NextPageIndex: pbResp.NextPageIndex,
	}
}

type NewsSource struct {
	Title  string `json:"title"`
	Domain string `json:"domain"`
	Region string `json:"region"`
}

type Currency struct {
	Code  string `json:"code"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

type News struct {
	Id            string     `json:"id"`
	CrawledSource string     `json:"crawled_source"`
	CreatedAt     time.Time  `json:"created_at"`
	Source        NewsSource `json:"source"`
	PublishedAt   time.Time  `json:"published_at"`
	Title         string     `json:"title"`
	ImageURL      string     `json:"image_url"`
	Description   string     `json:"description"`
	URL           string     `json:"url"`
	Currencies    []Currency `json:"currencies"`
}

func NewNewsFromPB(pbNews *pb.News) *News {
	if pbNews == nil {
		return nil
	}

	currencies := make([]Currency, 0)
	for _, pbC := range pbNews.Currencies {
		currencies = append(
			currencies,
			Currency{
				Code:  pbC.Code,
				Title: pbC.Title,
				URL:   pbC.Url,
			})
	}

	return &News{
		Id:            pbNews.Id,
		CrawledSource: pbNews.CrawledSource,
		CreatedAt:     pbNews.CreatedAt.AsTime(),
		Source: NewsSource{
			Title:  pbNews.Source.Title,
			Domain: pbNews.Source.Domain,
			Region: pbNews.Source.Region,
		},
		PublishedAt: pbNews.PublishedAt.AsTime(),
		Title:       pbNews.Title,
		ImageURL:    pbNews.ImageURL,
		Description: pbNews.Description,
		URL:         pbNews.Url,
		Currencies:  currencies,
	}
}
