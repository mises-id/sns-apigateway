package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	pb "github.com/mises-id/socialsvc/proto"
	grpcclient "github.com/mises-id/socialsvc/svc/client/grpc"
	"google.golang.org/grpc"
)

type PageQuickParams struct {
	Limit  int64  `json:"limit" query:"limit"`
	NextID string `json:"last_id" query:"last_id"`
}

// BuildSuccessResp return a success response with payload
func BuildSuccessResp(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, echo.Map{
		"code": 0,
		"data": data,
	})
}

// BuildSuccessResp return a success response with payload
func BuildSuccessRespWithPagination(c echo.Context, data interface{}, pagination *pb.PageQuick) error {
	return c.JSON(http.StatusOK, echo.Map{
		"code": 0,
		"data": data,
		"pagination": PageQuickParams{
			Limit:  int64(pagination.Limit),
			NextID: pagination.NextId,
		},
	})
}

// Probe for k8s liveness
func Probe(c echo.Context) error {
	return BuildSuccessResp(c, nil)
}

// build a service client, we are currently not using service discover
func GrpcSocialService() (pb.SocialServer, context.Context, error) {
	conn, err := grpc.Dial(":5040", grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create grpcclient: %q", err)
	}

	// Create a context with the header key and value
	ctx := context.WithValue(context.Background(), "key", "value")

	svcclient, err := grpcclient.New(conn)
	return svcclient, ctx, err
}
