package data

import (
	"blog/internal/biz"
	"blog/internal/data/ent/user"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ar *userRepo) GetUser(ctx context.Context, name string) (*biz.User, error) {
	p, err := ar.data.db.User.Query().Where(user.Name(name)).First(ctx)
	if err != nil {
		return nil, err
	}
	return &biz.User{
		Name:     p.Name,
		PassWord: p.Password,
		Email:    p.Email,
		Phone:    p.Phone,
		Address:  p.Address,
		Role:     p.Role,
		Desc:     p.Description,
	}, nil
}
