package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/server"
	"github.com/wendisx/puzzle/pkg/util"
)

/*  auth middleware for echo */
func (m EchoMiddleware) SimpleJwtAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			res := server.NewEchoResponder(c)
			tokenStr := c.Request().Header.Get("Authorization")
			clog.Debug(tokenStr)
			if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
				return res.Error(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			}
			tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
			jwtClaim, err := util.ParseToken(tokenStr)
			if err != nil || !jwtClaim.Check() {
				return res.Error(http.StatusUnauthorized, "Invalid Token")
			}
			// TODO: to be optimized...
			clog.Info(fmt.Sprintf("context.user{id=%s,name=%s}", string(jwtClaim.ExternId), jwtClaim.Name))
			c.Set("userId", jwtClaim.ExternId)
			c.Set("name", jwtClaim.Name)
			return next(c)
		}
	}
}

/* check middleware for echo */
func (m EchoMiddleware) ParseAndCheckBody(enable bool, s interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			res := server.NewEchoResponder(c)
			if err := util.Bind(c, s); err != nil {
				return res.Error(http.StatusBadRequest, "Invalid Request Payload.")
			}
			if enable {
				err := util.GetGlobalValidator().Check(s)
				if err != nil {
					return res.Error(http.StatusBadRequest, "Parameter Verification Failed.")
				}
			}
			clog.Debug(fmt.Sprintf("%#v", s))
			c.Set("body", s)
			return next(c)
		}
	}
}
