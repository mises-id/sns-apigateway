package middleware

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/mises-id/sns-apigateway/config/env"
	"github.com/mises-id/sns-apigateway/lib/codes"
)

var (
	secret           = env.Envs.JWTSecret
	validAuthMethods = []string{
		"Bearer",
	}
)

type UserSession struct {
	UID      uint64 `bson:"_id"`
	Username string `bson:"username,omitempty"`
	Misesid  string `bson:"misesid,omitempty"`
}

func Auth(ctx context.Context, authToken string) (*UserSession, error) {
	claim, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		if err.Error() == "Token is expired" {
			return nil, codes.ErrTokenExpired
		}
		return nil, err
	}
	mapClaims := claim.Claims.(jwt.MapClaims)
	return &UserSession{
		UID:      uint64(mapClaims["uid"].(float64)),
		Misesid:  mapClaims["misesid"].(string),
		Username: mapClaims["username"].(string),
	}, nil
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
			userSession, err := Auth(c.Request().Context(), strs[1])
			if err != nil {
				return err
			}
			c.Set("CurrentUser", userSession)
			c.Set("CurrentUID", userSession.UID)
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
