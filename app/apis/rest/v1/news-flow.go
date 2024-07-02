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

type ListNewsBeforeParams struct {
	NewsId *string `json:"news_id" query:"news_id"`
}

func ListNewsBefore(c echo.Context) error {
	params := &ListNewsBeforeParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.Newf("invalid query params")
	}

	grpcsvc, ctx, err := rest.GrpcNewsFlowService()
	if err != nil {
		return err
	}

	resp, err := grpcsvc.FindNewsInPageBefore(
		ctx,
		&pb.FindNewsInPageBeforeRequest{
			NewsId: params.NewsId,
		},
	)
	if err != nil {
		return err
	}

	return rest.BuildSuccessResp(c, NewListNewsBeforeResponseFromPB(resp))
}

type GetNewsParams struct {
	Id string `json:"id" query:"id"`
}

func GetNews(c echo.Context) error {
	params := &GetNewsParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.Newf("invalid query params")
	}

	grpcsvc, ctx, err := rest.GrpcNewsFlowService()
	if err != nil {
		return err
	}

	resp, err := grpcsvc.GetNewsById(
		ctx,
		&pb.GetNewsByIdRequest{
			Id: params.Id,
		},
	)
	if err != nil {
		return err
	}

	return rest.BuildSuccessResp(c, NewGetNewsByIdRespFromPB(resp))
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

type ListNewsBeforeResponse struct {
	NewsArray []*News `json:"news_array"`
	HaveMore  bool    `json:"have_more"`
}

func NewListNewsBeforeResponseFromPB(pbResp *pb.FindNewsInPageBeforeResponse) *ListNewsBeforeResponse {
	newsArray := make([]*News, 0)
	for _, pbNews := range pbResp.NewsArray {
		if news := NewNewsFromPB(pbNews); news != nil {
			newsArray = append(newsArray, news)
		}
	}
	return &ListNewsBeforeResponse{
		NewsArray: newsArray,
		HaveMore:  pbResp.HaveMore,
	}
}

func NewGetNewsByIdRespFromPB(pbNews *pb.News) *News {
	return NewNewsFromPB(pbNews)
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

type Media struct {
	Medium    string `json:"medium"`
	URL       string `json:"url"`
	Thumbnail string `json:"thumbnail"`
}

type News struct {
	Id            string     `json:"id"`
	CrawledSource string     `json:"crawled_source"`
	Source        NewsSource `json:"source"`
	PublishedAt   time.Time  `json:"published_at"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	Content       string     `json:"content"`
	Thumbnail     string     `json:"thumbnail"`
	Link          string     `json:"link"`
	Medias        []Media    `json:"medias"`
	Categories    []string   `json:"categories"`
	Currencies    []Currency `json:"currencies"`
}

func NewNewsFromPB(pbNews *pb.News) *News {
	if pbNews == nil {
		return nil
	}

	medias := make([]Media, 0)
	for _, pbM := range pbNews.Medias {
		medias = append(
			medias,
			Media{
				Medium:    pbM.Medium,
				URL:       pbM.Url,
				Thumbnail: pbM.Thumbnail,
			})
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

	thumbnail := pbNews.Thumbnail
	if len(thumbnail) == 0 {
		if len(pbNews.Medias) > 0 {
			media := pbNews.Medias[0]
			if len(media.Thumbnail) > 0 {
				thumbnail = media.Thumbnail
			} else if len(media.Url) > 0 {
				thumbnail = media.Url
			}
		}
	}

	return &News{
		Id:            pbNews.Id,
		CrawledSource: pbNews.CrawledSource,
		Source: NewsSource{
			Title:  pbNews.Source.Title,
			Domain: pbNews.Source.Domain,
			Region: pbNews.Source.Region,
		},
		PublishedAt: pbNews.PublishedAt.AsTime(),
		Title:       pbNews.Title,
		Description: pbNews.Description,
		Content:     pbNews.Content,
		Thumbnail:   thumbnail,
		Link:        pbNews.Link,
		Medias:      medias,
		Categories:  pbNews.Categories,
		Currencies:  currencies,
	}
}
