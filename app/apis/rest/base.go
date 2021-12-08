package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	grpcpool "github.com/go-kit/kit/util/grpcpool"
	"github.com/labstack/echo"
	pb "github.com/mises-id/sns-socialsvc/proto"
	grpcclient "github.com/mises-id/sns-socialsvc/svc/client/grpc"
	"google.golang.org/grpc"
)

var (
	socialSvcPool *grpcpool.Pool
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
	conn, err := socialSvcPool.Get(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create grpcclient: %q", err)
	}
	defer conn.Close()

	// Create a context with the header key and value
	ctx := context.WithValue(context.Background(), "key", "value")

	svcclient, err := grpcclient.New(conn.ClientConn)
	return svcclient, ctx, err
}

func init() {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	socialSvcPool, err = grpcpool.NewWithContext(ctx, func(ctx context.Context) (*grpc.ClientConn, error) {
		println("grpcpool", "new connection created")
		return grpc.DialContext(ctx, ":5040", grpc.WithInsecure())
	}, 0, 1, 60*time.Second)
	if err != nil {
		panic(err)
	}
}
