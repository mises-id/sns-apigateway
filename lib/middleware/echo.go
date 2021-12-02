package middleware

import (
	"github.com/labstack/echo"
	"github.com/mises-id/apigateway/lib/codes"
	log "github.com/sirupsen/logrus"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrorResponseMiddleware = func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			if _, ok := err.(*echo.HTTPError); ok {
				return err
			}
			statusErr, ok := status.FromError(err)
			if ok {
				switch statusErr.Code() {
				case grpccodes.NotFound:
					err = codes.ErrNotFound
				case grpccodes.PermissionDenied:
					err = codes.ErrForbidden
				case grpccodes.Unauthenticated:
					err = codes.ErrUnauthorized
				case grpccodes.AlreadyExists:
					err = codes.ErrUsernameExisted
				}
			}

			code, ok := err.(codes.Code)
			if !ok {
				log.WithFields(map[string]interface{}{
					"RequestID": c.Response().Header().Get(echo.HeaderXRequestID),
				}).Error("Unkown Error:", err)
				code = codes.ErrInternal
			}

			return c.JSON(code.HTTPStatus, code)
		}
		return nil
	}
}
