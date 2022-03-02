package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"trainings-go/internal/biz"
	"trainings-go/internal/data/model"
)

type userRepo struct {
	data  *Data
	log   *log.Helper
	table string
}

type userModel struct {
	Id   uint64
	Name string
}

// NewUserRepo
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data:  data,
		log:   log.NewHelper(logger),
		table: "iron_users",
	}
}

func (r *userRepo) GetUserInfo(ctx context.Context, hashId string) (id uint64) {
	r.log.WithContext(ctx).Infof("userRepo GetUserInfo Received:%s", hashId)
	data := userModel{}
	err := r.data.db.Table(r.table).Where("name = ?", hashId).Find(&data).Error
	if err != nil {
		r.log.WithContext(ctx).Error("userRepo GetUserInfo error", "error", err)
		return
	}
	id = data.Id
	r.log.WithContext(ctx).Info(data)
	return
}

func (r *userRepo) GetUserByName(ctx context.Context, name string) (data model.User, err error) {
	r.log.WithContext(ctx).Infof("userRepo GetUserInfo Received:%s", name)
	data = model.User{}
	err = r.data.db.Table("iron_users").Where("name = ?", name).Find(&data).Error
	if err != nil {
		r.log.WithContext(ctx).Error("userRepo GetUserByName error", "error", err)
		return
	}
	r.log.WithContext(ctx).Info(data)
	return
}

func (r *userRepo) SaveUser(ctx context.Context, user model.User) (err error) {
	r.log.WithContext(ctx).Infof("userRepo SaveUser Received:%s", user)
	err = r.data.db.Transaction(func(tx *gorm.DB) error {
		data := model.User{}
		err := r.data.db.Table(r.table).Where("name = ?", user.Name).Find(&data).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if data.Id > 0 {
			return errors.New(500, "USER_EXIST", "用户已存在")
		}
		err = r.data.db.Table(r.table).Save(&user).Error
		if err != nil {
			r.log.WithContext(ctx).Error("SaveUser error:%s", err)
			return err
		}
		return nil
	})
	return
}
