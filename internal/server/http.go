package server

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	actionApiV1 "trainings-go/api/action/v1"
	categoryApiV1 "trainings-go/api/category/v1"
	userApiV1 "trainings-go/api/user/v1"
	"trainings-go/internal/conf"
	"trainings-go/internal/service"
)

func NewWhiteListMatcher() selector.MatchFunc {

	whiteList := make(map[string]struct{})
	whiteList["/user.v1.User/Register"] = struct{}{}
	whiteList["/user.v1.User/Login"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Server, user *service.UserService, logger log.Logger, ca *conf.Auth, category *service.CategorySrv, action *service.ActionService) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			logging.Server(logger), // 记录接口请求日志
			// 白名单
			selector.Server(jwt.Server(func(token *jwtv4.Token) (interface{}, error) {
				return []byte(ca.GetKey()), nil
			}, jwt.WithSigningMethod(jwtv4.SigningMethodHS256))).Match(NewWhiteListMatcher()).Build(),
			// 限流
			ratelimit.Server(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	openAPIhandler := openapiv2.NewHandler()
	srv.HandlePrefix("/q/", openAPIhandler)
	userApiV1.RegisterUserHTTPServer(srv, user)
	categoryApiV1.RegisterCategoryHTTPServer(srv, category)
	actionApiV1.RegisterActionHTTPServer(srv, action)
	return srv
}
