package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"time"
	v1 "trainings-go/api/common/v1"
	userApiV1 "trainings-go/api/user/v1"
	"trainings-go/internal/data/model"
	"trainings-go/internal/pkg/util"
)

type UserBiz struct {
	repo UserRepo
	log  *log.Helper
}

type UserRepo interface {
	GetUserInfo(context.Context, string) uint64
	GetUserByName(ctx context.Context, name string) (model.User, error)
	SaveUser(ctx context.Context, user model.User) error
}

func NewUserBiz(repo UserRepo, logger log.Logger) *UserBiz {
	return &UserBiz{repo: repo, log: log.NewHelper(logger)}
}

func (bz *UserBiz) GetUserInfo(ctx context.Context, hashId string) (id uint64) {
	bz.log.WithContext(ctx).Infof("UserBiz GetUserInfo Received: %s", hashId)
	id = bz.repo.GetUserInfo(ctx, hashId)
	return
}

func (bz *UserBiz) Register(ctx context.Context, in *userApiV1.RegisterRequest) (err error) {
	bz.log.WithContext(ctx).Infof("UserBiz Register Received: %s", in)
	if err = bz.checkRegister(ctx, in); err != nil {
		return
	}
	err = bz.repo.SaveUser(ctx, model.User{
		Name:      in.GetName(),
		Password:  util.Md5(in.GetPassword()),
		NickName:  in.GetNickName(),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	})
	return
}

func (bz *UserBiz) checkRegister(ctx context.Context, in *userApiV1.RegisterRequest) (err error) {
	data, err := bz.repo.GetUserByName(ctx, in.GetName())
	if err != nil && err != gorm.ErrRecordNotFound {
		err = v1.ErrorUserExist("系统错误")
		return
	}
	if data.Id != 0 {
		err = v1.ErrorUserExist("用户已存在")
		return
	}
	if len(in.GetPassword()) < 6 {
		err = v1.ErrorUserExist("密码必须大于6位数")
		return
	}
	if in.GetPassword() != in.GetPassword2() {
		err = v1.ErrorUserExist("两次密码不一致")
		return
	}
	return
}

func (bz *UserBiz) Login(ctx context.Context, in *userApiV1.LoginRequest) (err error) {
	if len(in.GetName()) == 0 {
		err = v1.ErrorUserNotFound("请输入用户名")
		return
	}
	data, err := bz.repo.GetUserByName(ctx, in.GetName())
	if err != nil && err != gorm.ErrRecordNotFound {
		err = v1.ErrorUserExist("系统错误")
		return
	}
	if data.Id == 0 {
		err = v1.ErrorUserNotFound("用户不存在:%s", in.GetName())
		return
	}
	if data.Status != 0 {
		err = v1.ErrorUserNotFound("用户不存在:%s", in.GetName())
		return
	}
	if len(in.GetPassword()) == 0 {
		err = v1.ErrorUserExist("密码不能为空")
		return
	}
	if len(in.GetPassword()) < 6 {
		err = v1.ErrorUserExist("密码过短")
		return
	}
	if util.Md5(in.GetPassword()) != data.Password {
		err = v1.ErrorUserExist("密码")
		return
	}
	return
}
