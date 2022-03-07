package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	actionApiV1 "trainings-go/api/action/v1"
	categoryApiV1 "trainings-go/api/category/v1"
	userApiV1 "trainings-go/api/user/v1"
	"trainings-go/internal/conf"
	"trainings-go/internal/service"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, user *service.UserService, logger log.Logger, category *service.CategorySrv, action *service.ActionService) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			tracing.Server(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	userApiV1.RegisterUserServer(srv, user)
	categoryApiV1.RegisterCategoryServer(srv, category)
	actionApiV1.RegisterActionServer(srv, action)
	return srv
}
