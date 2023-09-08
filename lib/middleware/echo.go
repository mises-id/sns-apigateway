package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/mises-id/sns-apigateway/lib/codes"
	log "github.com/sirupsen/logrus"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrTooManyRequestFunc = func(c echo.Context, identifier string, err error) error {
	code := codes.ErrTooManyRequest
	if err != nil {
		code = code.New(err.Error())
	}
	return c.JSON(code.HTTPStatus, echo.Map{
		"code":    code.Code,
		"message": code.Msg,
	})
}
var ErrorResponseMiddleware = func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			if _, ok := err.(*echo.HTTPError); ok {
				return err
			}
			var msg string
			statusErr, ok := status.FromError(err)
			if ok {
				msg = statusErr.Message()
				switch statusErr.Code() {
				case grpccodes.NotFound:
					err = codes.ErrNotFound
				case grpccodes.InvalidArgument:
					err = codes.ErrInvalidArgument
				case grpccodes.PermissionDenied:
					err = codes.ErrForbidden
				case grpccodes.Unauthenticated:
					err = codes.ErrUnauthorized
				case grpccodes.AlreadyExists:
					err = codes.ErrUsernameExisted
				case grpccodes.DeadlineExceeded:
					err = codes.ErrRequestTimeoutCode
				case grpccodes.Unavailable:
					err = codes.ErrInternal
				case grpccodes.DeadlineExceeded:
					err = codes.ErrRequestTimeout
				}
			}

			code, ok := err.(codes.Code)
			if !ok {
				log.WithFields(map[string]interface{}{
					"RequestID": c.Response().Header().Get(echo.HeaderXRequestID),
					"Uri":       c.Request().RequestURI,
				}).Error("Unkown Error:", err)
				code = codes.ErrInternal
			} else {
				if msg != "" {
					code = code.New(msg)
				}
			}

			return c.JSON(code.HTTPStatus, code)
		}
		return nil
	}
}
