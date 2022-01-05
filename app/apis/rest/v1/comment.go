package v1

import (
	"github.com/labstack/echo"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	"github.com/mises-id/sns-apigateway/lib/codes"

	pb "github.com/mises-id/sns-socialsvc/proto"
)

type CreateCommentParams struct {
	CommentableID string `json:"status_id"`
	ParentID      string `json:"parent_id"`
	Content       string `json:"content"`
}

type ListCommentParams struct {
	rest.PageQuickParams
	CommentableID string `query:"status_id"`
	TopicID       string `query:"topic_id"`
}

func BuildCommentRespSlice(in []*pb.Comment) []*pb.Comment {
	if in == nil {
		return []*pb.Comment{}
	} else {
		return in
	}
}

func ListComment(c echo.Context) error {
	params := &ListCommentParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.Newf("invalid query params")
	}
	var currentUID uint64
	if c.Get("CurrentUID") != nil {
		currentUID = c.Get("CurrentUID").(uint64)
	}

	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.ListComment(ctx, &pb.ListCommentRequest{
		CurrentUid: currentUID,
		StatusId:   params.CommentableID,
		TopicId:    params.TopicID,
		Paginator: &pb.PageQuick{
			NextId: params.PageQuickParams.NextID,
			Limit:  uint64(params.PageQuickParams.Limit),
		},
	})
	if err != nil {
		return err
	}

	return rest.BuildSuccessRespWithPagination(c, BuildCommentRespSlice(svcresp.Comments), svcresp.Paginator)
}

func CreateComment(c echo.Context) error {
	params := &CreateCommentParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid comment params")
	}
	var currentUID uint64
	if c.Get("CurrentUID") != nil {
		currentUID = c.Get("CurrentUID").(uint64)
	} else {
		return codes.ErrInvalidArgument
	}

	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.CreateComment(ctx, &pb.CreateCommentRequest{
		CurrentUid: currentUID,
		StatusId:   params.CommentableID,
		Content:    params.Content,
		ParentId:   params.ParentID,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, svcresp.Comment)
}
