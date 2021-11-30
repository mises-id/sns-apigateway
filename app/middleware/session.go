package middleware

import (
	"strings"

	"github.com/labstack/echo"
	"github.com/mises-id/apigateway/lib/codes"
	"github.com/mises-id/socialsvc/app/services/session"
)

var (
	validAuthMethods = []string{
		"Bearer",
	}
)

type UserSession struct {
	UID      uint64 `bson:"_id"`
	Username string `bson:"username,omitempty"`
	Misesid  string `bson:"misesid,omitempty"`
}

var SetCurrentUserMiddleware = func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorization := c.Request().Header.Get("Authorization")
		if authorization != "" {
			strs := strings.Split(authorization, " ")
			if err := validateAuthToken(strs); err != nil {
				return err
			}

			//TODO we should move Auth to this repo
			user, err := session.Auth(c.Request().Context(), strs[1])
			if err != nil {
				return err
			}
			userSession := UserSession{
				UID:      user.UID,
				Username: user.Username,
				Misesid:  user.Misesid,
			}
			c.Set("CurrentUser", userSession)
			c.Set("CurrentUID", user.UID)
		}

		return next(c)
	}
}

var RequireCurrentUserMiddleware = func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("CurrentUser").(*UserSession)
		if !ok || user == nil {
			return codes.ErrUnauthorized
		}
		return next(c)
	}
}

func validateAuthToken(strs []string) error {
	if len(strs) != 2 {
		return codes.ErrInvalidAuth
	}
	authMethod, authToken := strs[0], strs[1]
	if len(authToken) > 1000 {
		return codes.ErrInvalidAuthToken
	}
	for _, m := range validAuthMethods {
		if m == authMethod {
			return nil
		}
	}
	return codes.ErrInvalidAuthMethod
}
