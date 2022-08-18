package v1

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	"github.com/mises-id/sns-apigateway/lib/codes"

	pb "github.com/mises-id/sns-socialsvc/proto"
)

type CreateCommentParams struct {
	CommentableID string `json:"status_id"`
	NftAssetID    string `json:"nft_asset_id"`
	ParentID      string `json:"parent_id"`
	Content       string `json:"content"`
}

type ListCommentParams struct {
	rest.PageQuickParams
	CommentableID string `query:"status_id"`
	NftAssetID    string `query:"nft_asset_id"`
	TopicID       string `query:"topic_id"`
}

type CommentResp struct {
	ID            string         `json:"id"`
	ParentID      string         `json:"parent_id"`
	TopicID       string         `json:"topic_id"`
	Content       string         `json:"content"`
	Comments      []*CommentResp `json:"comments"`
	CommentsCount uint64         `json:"comments_count"`
	LikesCount    uint64         `json:"likes_count"`

	User     *UserSummaryResp `json:"user,omitempty"`
	Opponent *UserSummaryResp `json:"opponent,omitempty"`

	IsLiked   bool      `json:"is_liked"`
	CreatedAt time.Time `json:"created_at"`
}

func BuildCommentResp(in *pb.Comment) *CommentResp {
	if in == nil {
		return &CommentResp{}
	} else {
		return &CommentResp{
			ID:            in.Id,
			ParentID:      in.ParentId,
			TopicID:       in.GroupId,
			Content:       in.Content,
			Comments:      BuildCommentRespSlice(in.Comments),
			CommentsCount: in.CommentCount,
			LikesCount:    in.LikeCount,
			IsLiked:       in.IsLiked,
			CreatedAt:     time.Unix(int64(in.CreatedAt), 0),
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
func GetComment(c echo.Context) error {

	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.GetComment(ctx, &pb.GetCommentRequest{
		CurrentUid: GetCurrentUID(c),
		CommentId:  c.Param("id"),
	})
	if err != nil {
		return err
	}

	return rest.BuildSuccessResp(c, BuildCommentResp(svcresp.Comment))
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
		NftAssetId: params.NftAssetID,
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
		NftAssetId: params.NftAssetID,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, BuildCommentResp(svcresp.Comment))
}

func DeleteComment(c echo.Context) error {
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	_, err = grpcsvc.DeleteComment(ctx, &pb.DeleteCommentRequest{
		CurrentUid: GetCurrentUID(c),
		Id:         c.Param("id"),
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, nil)
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
