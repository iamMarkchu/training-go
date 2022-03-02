package service

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	v1 "trainings-go/api/common/v1"
	userApiV1 "trainings-go/api/user/v1"
	"trainings-go/internal/biz"
)

type UserService struct {
	userApiV1.UnimplementedUserServer
	userBiz *biz.UserBiz
	log     *log.Helper
}

func NewUserService(userBiz *biz.UserBiz, logger log.Logger) *UserService {
	return &UserService{userBiz: userBiz, log: log.NewHelper(logger)}
}

// SayHello implements user.GreeterServer
func (s *UserService) SayHello(ctx context.Context, in *userApiV1.HelloRequest) (*userApiV1.HelloReply, error) {
	s.log.WithContext(ctx).Infof("User SayHello Received: %v", in.GetName())

	if in.GetName() == "error" {
		return nil, v1.ErrorUserNotFound("user not found: %s", in.GetName())
	}
	id := s.userBiz.GetUserInfo(ctx, in.GetName())
	return &userApiV1.HelloReply{Message: "Hello " + in.GetName() + " Id:" + fmt.Sprintf("%d", id)}, nil
}

func (s *UserService) Register(ctx context.Context, in *userApiV1.RegisterRequest) (*userApiV1.RegisterResponse, error) {
	s.log.WithContext(ctx).Info("User Register Received ", "params", in)
	err := s.userBiz.Register(ctx, in)
	if err != nil {
		return &userApiV1.RegisterResponse{}, err
	}
	return &userApiV1.RegisterResponse{}, nil
}

func (s *UserService) Login(ctx context.Context, in *userApiV1.LoginRequest) (*userApiV1.LoginResponse, error) {
	s.log.WithContext(ctx).Info("User Login Received ", "params", in)
	t, ex, err := s.userBiz.Login(ctx, in)
	if err != nil {
		return &userApiV1.LoginResponse{}, err
	}
	return &userApiV1.LoginResponse{
		Data: &userApiV1.LoginResponse_LoginData{
			Expire: uint64(ex),
			Token:  t,
		},
	}, nil
}
