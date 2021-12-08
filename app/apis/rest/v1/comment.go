package v1

import (
	"github.com/labstack/echo"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	"github.com/mises-id/sns-apigateway/lib/codes"

	pb "github.com/mises-id/sns-socialsvc/proto"
)

type CreateCommentParams struct {
	CommentableID string `json:"status_id"`
	Content       string `json:"content"`
}

type ListCommentParams struct {
	rest.PageQuickParams
	CommentableID string `query:"status_id"`
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
	svcresp, err := grpcsvc.ListStatus(ctx, &pb.ListStatusRequest{
		CurrentUid: currentUID,
		TargetUid:  0,
		ParrentId:  params.CommentableID,
		FromTypes:  []string{"comment"},
		Paginator: &pb.PageQuick{
			NextId: params.PageQuickParams.NextID,
			Limit:  uint64(params.PageQuickParams.Limit),
		},
	})
	if err != nil {
		return err
	}

	return rest.BuildSuccessRespWithPagination(c, svcresp.Statuses, svcresp.Paginator)
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
	svcresp, err := grpcsvc.CreateStatus(ctx, &pb.CreateStatusRequest{
		CurrentUid: currentUID,
		StatusType: "text",
		Content:    params.Content,
		ParentId:   params.CommentableID,
		FromType:   "comment",
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, svcresp.Status)
}
