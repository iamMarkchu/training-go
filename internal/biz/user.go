package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"time"
	v1 "trainings-go/api/common/v1"
	userCoreApi "trainings-go/api/user/core"
	userApiV1 "trainings-go/api/user/v1"
	"trainings-go/internal/conf"
	"trainings-go/internal/data/model"
	"trainings-go/internal/pkg/define"
	"trainings-go/internal/pkg/util"
)

type UserBiz struct {
	repo UserRepo
	log  *log.Helper
	conf *conf.Auth
}

type UserRepo interface {
	GetUserInfo(context.Context, string) uint64
	GetUserByName(ctx context.Context, name string) (model.User, error)
	SaveUser(context.Context, model.User) error
	MGetUserInfo(context.Context, []uint64) ([]model.User, error)
}

func NewUserBiz(repo UserRepo, logger log.Logger, conf *conf.Auth) *UserBiz {
	return &UserBiz{repo: repo, log: log.NewHelper(logger), conf: conf}
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

func (bz *UserBiz) Login(ctx context.Context, in *userApiV1.LoginRequest) (t string, ex int64, err error) {
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
	claims := &define.ApiClaims{
		data.Id,
		jwtv4.RegisteredClaims{
			NotBefore: jwtv4.NewNumericDate(time.Now()),
			ExpiresAt: jwtv4.NewNumericDate(time.Now().Add(time.Second * 86400)),
			Issuer:    "147258",
		},
	}
	token := jwtv4.NewWithClaims(jwtv4.SigningMethodHS256, claims)
	t, err = token.SignedString([]byte(bz.conf.Key))
	if err != nil {
		return
	}
	ex = 86400
	return
}

func (bz *UserBiz) MGetUserInfo(ctx context.Context, in *userCoreApi.MGetInfoRequest) (res map[uint64]*userCoreApi.User, err error) {
	res = make(map[uint64]*userCoreApi.User)
	list, err := bz.repo.MGetUserInfo(ctx, in.GetUids())
	if err != nil {
		return
	}
	for _, user := range list {
		res[user.Id] = &userCoreApi.User{
			Id:       user.Id,
			Name:     user.Name,
			NickName: user.NickName,
		}
	}
	return
}
