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

type CommentResp struct {
	ID            string         `json:"id"`
	TopicID       string         `json:"topic_id"`
	Content       string         `json:"content"`
	Comments      []*CommentResp `json:"comments"`
	CommentsCount uint64         `json:"comments_count"`
	LikesCount    uint64         `json:"likes_count"`

	User     *UserSummaryResp `json:"user,omitempty"`
	Opponent *UserSummaryResp `json:"opponent,omitempty"`
}

func BuildCommentResp(in *pb.Comment) *CommentResp {
	if in == nil {
		return &CommentResp{}
	} else {
		return &CommentResp{
			ID:            in.Id,
			TopicID:       in.GroupId,
			Content:       in.Content,
			Comments:      BuildCommentRespSlice(in.Comments),
			CommentsCount: in.CommentCount,
			LikesCount:    in.LikeCount,
			User:          BuildUserSummaryResp(in.User),
			Opponent:      BuildUserSummaryResp(in.Opponent),
		}
	}
}

func BuildCommentRespSlice(in []*pb.Comment) []*CommentResp {
	if in == nil {
		return []*CommentResp{}
	}

	resp := []*CommentResp{}
	for _, i := range in {
		resp = append(resp, BuildCommentResp(i))
	}

	return resp

}

func ListComment(c echo.Context) error {
	params := &ListCommentParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.Newf("invalid query params")
	}

	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.ListComment(ctx, &pb.ListCommentRequest{
		CurrentUid: GetCurrentUID(c),
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
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.CreateComment(ctx, &pb.CreateCommentRequest{
		CurrentUid: GetCurrentUID(c),
		StatusId:   params.CommentableID,
		Content:    params.Content,
		ParentId:   params.ParentID,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, BuildCommentResp(svcresp.Comment))
}

func LikeComment(c echo.Context) error {
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}

	_, err = grpcsvc.LikeComment(ctx, &pb.LikeCommentRequest{
		CurrentUid: GetCurrentUID(c),
		CommentId:  c.Param("id"),
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, nil)
}

func UnlikeComment(c echo.Context) error {
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}

	_, err = grpcsvc.UnlikeComment(ctx, &pb.UnlikeCommentRequest{
		CurrentUid: GetCurrentUID(c),
		CommentId:  c.Param("id"),
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, nil)
}
