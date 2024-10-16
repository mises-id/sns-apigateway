package rest

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	grpcpool "github.com/go-kit/kit/util/grpcpool"
	"github.com/labstack/echo/v4"
	airdropsvcpb "github.com/mises-id/mises-airdropsvc/proto"
	airdropsvcgrpcclient "github.com/mises-id/mises-airdropsvc/svc/client/grpc"
	miningsvcpb "github.com/mises-id/mises-miningsvc/proto"
	miningsvcgrpcclient "github.com/mises-id/mises-miningsvc/svc/client/grpc"
	newsflowpb "github.com/mises-id/mises-news-flow/pkg/proto/apiserver/v1"
	swapvcpb "github.com/mises-id/mises-swapsvc/proto"
	swapsvcgrpcclient "github.com/mises-id/mises-swapsvc/svc/client/grpc"
	websitesvcpb "github.com/mises-id/mises-websitesvc/proto"
	websitesvcgrpcclient "github.com/mises-id/mises-websitesvc/svc/client/grpc"
	pb "github.com/mises-id/sns-socialsvc/proto"
	grpcclient "github.com/mises-id/sns-socialsvc/svc/client/grpc"
	storagepb "github.com/mises-id/sns-storagesvc/proto"
	storagesvcgrpcclient "github.com/mises-id/sns-storagesvc/svc/client/grpc"
	"google.golang.org/grpc"
)

var (
	socialSvcPool   *grpcpool.Pool
	storageSvcPool  *grpcpool.Pool
	websiteSvcPool  *grpcpool.Pool
	airdropSvcPool  *grpcpool.Pool
	swapSvcPool     *grpcpool.Pool
	mingsvcSvcPool  *grpcpool.Pool
	newsFlowSvcPool *grpcpool.Pool
	store           sync.Map
)

type PoolCfg struct {
	SocialSvcURI   string
	StorageSvcURI  string
	WebsiteSvcURI  string
	AirdropSvcURI  string
	SwapSvcURI     string
	MiningSvcURI   string
	NewsFlowSvcURI string
	Capacity       int
	IdleTimeout    time.Duration
}
type PageQuickParams struct {
	Limit  int64  `json:"limit" query:"limit"`
	Total  int64  `json:"total" query:"total"`
	NextID string `json:"last_id" query:"last_id"`
}
type PageParams struct {
	PageNum      int64 `json:"page_num" query:"page_num"`
	PageSize     int64 `json:"page_size" query:"page_size"`
	TotalPage    int64 `json:"total_page"`
	TotalRecords int64 `json:"total_records"`
}

// BuildSuccessResp return a success response with payload
func Build403Resp(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusForbidden, echo.Map{
		"code": http.StatusForbidden,
		"data": data,
	})
}

// BuildSuccessResp return a success response with payload
func BuildSuccessResp(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, echo.Map{
		"code": 0,
		"data": data,
	})
}

func BuildSuccessRespWithRequestID(c echo.Context, requestID string, data interface{}) error {
	return c.JSON(http.StatusOK, echo.Map{
		"code":       0,
		"data":       data,
		"request_id": requestID,
	})
}
func BuildSuccessRespWithWebsitePageAndRequestID(c echo.Context, requestID string, data interface{}, pagination *websitesvcpb.Page) error {
	return c.JSON(http.StatusOK, echo.Map{
		"code":       0,
		"data":       data,
		"request_id": requestID,
		"pagination": PageParams{
			PageNum:      int64(pagination.PageNum),
			PageSize:     int64(pagination.PageSize),
			TotalRecords: int64(pagination.TotalRecords),
			TotalPage:    int64(pagination.TotalPage),
		},
	})
}
func BuildSuccessRespWithSwapPage(c echo.Context, requestID string, data interface{}, pagination *swapvcpb.Page) error {
	return c.JSON(http.StatusOK, echo.Map{
		"code":       0,
		"data":       data,
		"request_id": requestID,
		"pagination": PageParams{
			PageNum:      int64(pagination.PageNum),
			PageSize:     int64(pagination.PageSize),
			TotalRecords: int64(pagination.TotalRecords),
			TotalPage:    int64(pagination.TotalPage),
		},
	})
}

// BuildSuccessResp return a success response with payload
func BuildSuccessRespWithPagination(c echo.Context, data interface{}, pagination *pb.PageQuick) error {
	return c.JSON(http.StatusOK, echo.Map{
		"code": 0,
		"data": data,
		"pagination": PageQuickParams{
			Limit:  int64(pagination.Limit),
			Total:  int64(pagination.Total),
			NextID: pagination.NextId,
		},
	})
}
func BuildSuccessRespWebsiteWithPagination(c echo.Context, data interface{}, pagination *websitesvcpb.PageQuick) error {
	return c.JSON(http.StatusOK, echo.Map{
		"code": 0,
		"data": data,
		"pagination": PageQuickParams{
			Limit:  int64(pagination.Limit),
			Total:  int64(pagination.Total),
			NextID: pagination.NextId,
		},
	})
}
func BuildSuccessRespWithWebsitePage(c echo.Context, data interface{}, pagination *websitesvcpb.Page) error {
	return c.JSON(http.StatusOK, echo.Map{
		"code": 0,
		"data": data,
		"pagination": PageParams{
			PageNum:      int64(pagination.PageNum),
			PageSize:     int64(pagination.PageSize),
			TotalRecords: int64(pagination.TotalRecords),
			TotalPage:    int64(pagination.TotalPage),
		},
	})
}
func BuildSuccessRespAirdropWithPagination(c echo.Context, data interface{}, pagination *airdropsvcpb.PageQuick) error {
	return c.JSON(http.StatusOK, echo.Map{
		"code": 0,
		"data": data,
		"pagination": PageQuickParams{
			Limit:  int64(pagination.Limit),
			Total:  int64(pagination.Total),
			NextID: pagination.NextId,
		},
	})
}
func BuildSuccessRespWithAirdropPage(c echo.Context, data interface{}, pagination *airdropsvcpb.Page) error {
	return c.JSON(http.StatusOK, echo.Map{
		"code": 0,
		"data": data,
		"pagination": PageParams{
			PageNum:      int64(pagination.PageNum),
			PageSize:     int64(pagination.PageSize),
			TotalRecords: int64(pagination.TotalRecords),
			TotalPage:    int64(pagination.TotalPage),
		},
	})
}
func BuildSuccessRespWithPage(c echo.Context, data interface{}, pagination *pb.Page) error {
	return c.JSON(http.StatusOK, echo.Map{
		"code": 0,
		"data": data,
		"pagination": PageParams{
			PageNum:      int64(pagination.PageNum),
			PageSize:     int64(pagination.PageSize),
			TotalRecords: int64(pagination.TotalRecords),
			TotalPage:    int64(pagination.TotalPage),
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

func GrpcStorageService() (storagepb.StoragesvcServer, context.Context, error) {
	conn, err := storageSvcPool.Get(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create grpcclient: %q", err)
	}
	defer conn.Close()

	// Create a context with the header key and value
	ctx := context.WithValue(context.Background(), "key", "value")

	svcclient, err := storagesvcgrpcclient.New(conn.ClientConn)
	return svcclient, ctx, err
}
func GrpcWebsiteService() (websitesvcpb.WebsitesvcServer, context.Context, error) {
	conn, err := websiteSvcPool.Get(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create grpcclient: %q", err)
	}
	defer conn.Close()

	// Create a context with the header key and value
	ctx := context.WithValue(context.Background(), "key", "value")

	svcclient, err := websitesvcgrpcclient.New(conn.ClientConn)
	return svcclient, ctx, err
}
func GrpcSwapService() (swapvcpb.SwapsvcServer, context.Context, error) {
	conn, err := swapSvcPool.Get(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create grpcclient: %q", err)
	}
	defer conn.Close()

	// Create a context with the header key and value
	ctx := context.WithValue(context.Background(), "key", "value")

	svcclient, err := swapsvcgrpcclient.New(conn.ClientConn)
	return svcclient, ctx, err
}

func GrpcMiningService() (miningsvcpb.MiningsvcServer, context.Context, error) {
	conn, err := mingsvcSvcPool.Get(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create grpcclient: %q", err)
	}
	defer conn.Close()

	// Create a context with the header key and value
	ctx := context.WithValue(context.Background(), "key", "value")

	svcclient, err := miningsvcgrpcclient.New(conn.ClientConn)
	return svcclient, ctx, err
}

func GrpcAirdropService() (airdropsvcpb.AirdropsvcServer, context.Context, error) {
	conn, err := airdropSvcPool.Get(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create grpcclient: %q", err)
	}
	defer conn.Close()

	// Create a context with the header key and value
	ctx := context.WithValue(context.Background(), "key", "value")

	svcclient, err := airdropsvcgrpcclient.New(conn.ClientConn)
	return svcclient, ctx, err
}

func GrpcNewsFlowService() (newsflowpb.ApiserverClient, context.Context, error) {
	conn, err := newsFlowSvcPool.Get(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create news-flow grpcclient: %q", err)
	}
	defer conn.Close()

	// Create a context with the header key and value
	ctx := context.WithValue(context.Background(), "key", "value")

	client := newsflowpb.NewApiserverClient(conn.ClientConn)
	return client, ctx, nil
}

func InMemoryStore() *sync.Map {
	return &store
}

func ResetSvrPool(cfg PoolCfg) {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	socialSvcPool, err = grpcpool.NewWithContext(ctx, func(ctx context.Context) (*grpc.ClientConn, error) {
		return grpc.DialContext(ctx, cfg.SocialSvcURI, grpc.WithInsecure())
	}, 0, cfg.Capacity, cfg.IdleTimeout*time.Second)

	storageSvcPool, err = grpcpool.NewWithContext(ctx, func(ctx context.Context) (*grpc.ClientConn, error) {
		return grpc.DialContext(ctx, cfg.StorageSvcURI, grpc.WithInsecure())
	}, 0, cfg.Capacity, cfg.IdleTimeout*time.Second)

	websiteSvcPool, err = grpcpool.NewWithContext(ctx, func(ctx context.Context) (*grpc.ClientConn, error) {
		return grpc.DialContext(ctx, cfg.WebsiteSvcURI, grpc.WithInsecure())
	}, 0, cfg.Capacity, cfg.IdleTimeout*time.Second)

	airdropSvcPool, err = grpcpool.NewWithContext(ctx, func(ctx context.Context) (*grpc.ClientConn, error) {
		return grpc.DialContext(ctx, cfg.AirdropSvcURI, grpc.WithInsecure())
	}, 0, cfg.Capacity, cfg.IdleTimeout*time.Second)

	swapSvcPool, err = grpcpool.NewWithContext(ctx, func(ctx context.Context) (*grpc.ClientConn, error) {
		return grpc.DialContext(ctx, cfg.SwapSvcURI, grpc.WithInsecure())
	}, 0, cfg.Capacity, cfg.IdleTimeout*time.Second)

	mingsvcSvcPool, err = grpcpool.NewWithContext(ctx, func(ctx context.Context) (*grpc.ClientConn, error) {
		return grpc.DialContext(ctx, cfg.MiningSvcURI, grpc.WithInsecure())
	}, 0, cfg.Capacity, cfg.IdleTimeout*time.Second)

	newsFlowSvcPool, err = grpcpool.NewWithContext(ctx, func(ctx context.Context) (*grpc.ClientConn, error) {
		return grpc.DialContext(ctx, cfg.NewsFlowSvcURI, grpc.WithInsecure())
	}, 0, cfg.Capacity, cfg.IdleTimeout*time.Second)

	if err != nil {
		panic(err)
	}
}

func init() {
	ResetSvrPool(
		PoolCfg{
			SocialSvcURI:   ":5040",
			StorageSvcURI:  ":6040",
			WebsiteSvcURI:  ":4040",
			AirdropSvcURI:  ":3040",
			SwapSvcURI:     ":4540",
			MiningSvcURI:   ":3540",
			NewsFlowSvcURI: ":7000",
			Capacity:       1,
			IdleTimeout:    60,
		})
}
