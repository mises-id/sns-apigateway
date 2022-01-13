package v1

import (
	"encoding/json"
	"time"

	"github.com/labstack/echo"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	"github.com/mises-id/sns-apigateway/lib/codes"
	pb "github.com/mises-id/sns-socialsvc/proto"
)

type ListUserStatusParams struct {
	rest.PageQuickParams
}

type LinkMeta struct {
	Title          string `json:"title"`
	Host           string `json:"host"`
	Link           string `json:"link"`
	AttachmentPath string `json:"attachment_path"`
}

type CreateStatusParams struct {
	FromType     string        `json:"from_type"`
	StatusType   string        `json:"status_type"`
	ParentID     string        `json:"parent_id"`
	Content      string        `json:"content"`
	LinkMeta     *LinkMeta     `json:"link_meta"`
	Images       []string      `json:"images"`
	IsPrivate    bool          `json:"is_private"`
	ShowDuration time.Duration `json:"show_duration"`
}

type LinkMetaResp struct {
	Title          string `json:"title"`
	Host           string `json:"host"`
	Link           string `json:"link"`
	AttachmentPath string `json:"attachment_path"`
	AttachmentURL  string `json:"attachment_url"`
}

type StatusResp struct {
	ID            string           `json:"id"`
	User          *UserSummaryResp `json:"user"`
	Content       string           `json:"content"`
	FromType      string           `json:"from_type"`
	StatusType    string           `json:"status_type"`
	ParentStatus  *StatusResp      `json:"parent_status"`
	OriginStatus  *StatusResp      `json:"origin_status"`
	CommentsCount uint64           `json:"comments_count"`
	LikesCount    uint64           `json:"likes_count"`
	ForwardsCount uint64           `json:"forwards_count"`
	IsLiked       bool             `json:"is_liked"`
	LinkMeta      *LinkMetaResp    `json:"link_meta"`
	CreatedAt     time.Time        `json:"created_at"`
	ThumbImages   []string         `json:"thumb_images"`
	Images        []string         `json:"images"`
	IsPublic      bool             `json:"is_public"`
	HideTime      time.Time        `json:"hide_time"`
}

func BuildStatusResp(info *pb.StatusInfo) *StatusResp {
	if info == nil {
		return nil
	}
	resp := &StatusResp{
		ID:            info.Id,
		User:          BuildUserSummaryResp(info.User),
		Content:       info.Content,
		FromType:      info.FromType,
		StatusType:    info.StatusType,
		CommentsCount: info.CommentCount,
		LikesCount:    info.LikeCount,
		ForwardsCount: info.ForwardCount,
		IsLiked:       info.IsLiked,
		CreatedAt:     time.Unix(int64(info.CreatedAt), 0),
		IsPublic:      info.IsPublic,
		HideTime:      time.Unix(int64(info.HideTime), 0),
	}

	if info.LinkMeta != nil {
		resp.LinkMeta = &LinkMetaResp{
			Title:          info.LinkMeta.Title,
			Host:           info.LinkMeta.Host,
			Link:           info.LinkMeta.Link,
			AttachmentPath: info.LinkMeta.ImagePath,
			AttachmentURL:  info.LinkMeta.ImageUrl,
		}
	}
	if info.Parent != nil {
		resp.ParentStatus = BuildStatusResp(info.Parent)
	}
	if info.Origin != nil {
		resp.OriginStatus = BuildStatusResp(info.Origin)
	}
	if info.ImageMeta != nil {
		resp.Images = info.ImageMeta.Images
		resp.ThumbImages = info.ImageMeta.Images
	} else {
		resp.Images = []string{}
		resp.ThumbImages = []string{}
	}

	return resp
}
func BuildStatusRespSlice(infos []*pb.StatusInfo) []*StatusResp {
	resp := []*StatusResp{}
	for _, info := range infos {
		resp = append(resp, BuildStatusResp(info))
	}

	return resp
}
func GetStatus(c echo.Context) error {

	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.GetStatus(ctx, &pb.GetStatusRequest{
		CurrentUid: GetCurrentUID(c),
		Statusid:   c.Param("id"),
	})
	if err != nil {
		return err
	}

	return rest.BuildSuccessResp(c, BuildStatusResp(svcresp.Status))
}

// list user status
func ListUserStatus(c echo.Context) error {
	uid, err := GetUIDParam(c)
	if err != nil {
		return err
	}

	params := &ListUserStatusParams{}
	if err = c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}

	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.ListStatus(ctx, &pb.ListStatusRequest{
		CurrentUid: GetCurrentUID(c),
		TargetUid:  uid,
		FromTypes:  []string{"post", "forward"},
		Paginator: &pb.PageQuick{
			NextId: params.PageQuickParams.NextID,
			Limit:  uint64(params.PageQuickParams.Limit),
		},
	})
	if err != nil {
		return err
	}

	return rest.BuildSuccessRespWithPagination(c, BuildStatusRespSlice(svcresp.Statuses), svcresp.Paginator)
}

func Timeline(c echo.Context) error {
	params := &ListUserStatusParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}

	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.ListUserTimeline(ctx, &pb.ListStatusRequest{
		CurrentUid: GetCurrentUID(c),
		Paginator: &pb.PageQuick{
			NextId: params.PageQuickParams.NextID,
			Limit:  uint64(params.PageQuickParams.Limit),
		},
	})
	if err != nil {
		return err
	}

	return rest.BuildSuccessRespWithPagination(c, BuildStatusRespSlice(svcresp.Statuses), svcresp.Paginator)
}

func RecommendStatus(c echo.Context) error {

	params := &ListUserStatusParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}

	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.ListRecommended(ctx, &pb.ListStatusRequest{
		CurrentUid: GetCurrentUID(c),
		Paginator: &pb.PageQuick{
			NextId: params.PageQuickParams.NextID,
			Limit:  uint64(params.PageQuickParams.Limit),
		},
	})
	if err != nil {
		return err
	}

	return rest.BuildSuccessRespWithPagination(c, BuildStatusRespSlice(svcresp.Statuses), svcresp.Paginator)

}

func UpdateStatus(c echo.Context) error {
	params := &CreateStatusParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid status params")
	}

	return rest.BuildSuccessResp(c, nil)
}

func CreateStatus(c echo.Context) error {
	params := &CreateStatusParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid status params")
	}
	fromType := "post"
	if len(params.ParentID) > 0 {
		fromType = "forward"
	}
	var meta json.RawMessage
	var err error
	if params.LinkMeta != nil {
		if meta, err = json.Marshal(params.LinkMeta); err != nil {
			return err
		}
	}
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}

	svcresp, err := grpcsvc.CreateStatus(ctx, &pb.CreateStatusRequest{
		CurrentUid: GetCurrentUID(c),
		StatusType: params.StatusType,
		Content:    params.Content,
		ParentId:   params.ParentID,
		Meta:       string(meta),
		FromType:   fromType,
		Images:     params.Images,
	})
	if err != nil {
		return err
	}

	return rest.BuildSuccessResp(c, BuildStatusResp(svcresp.Status))
}

func DeleteStatus(c echo.Context) error {
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}

	_, err = grpcsvc.DeleteStatus(ctx, &pb.DeleteStatusRequest{
		CurrentUid: GetCurrentUID(c),
		Statusid:   c.Param("id"),
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, nil)
}

func LikeStatus(c echo.Context) error {
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}

	_, err = grpcsvc.LikeStatus(ctx, &pb.LikeStatusRequest{
		CurrentUid: GetCurrentUID(c),
		Statusid:   c.Param("id"),
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, nil)
}

func UnlikeStatus(c echo.Context) error {
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}

	_, err = grpcsvc.UnLikeStatus(ctx, &pb.UnLikeStatusRequest{
		CurrentUid: GetCurrentUID(c),
		Statusid:   c.Param("id"),
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, nil)
}
