package biz

import "context"

type User struct {
	Name     string `json:"name"`
	PassWord string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Role     string `json:"role"`
	Desc     string `json:"description"`
}

type UserRepo interface {
	GetUser(ctx context.Context, name string) (*User, error)
}

type UserUsecase struct {
	repo UserRepo
}

func NewUserUsecase(repo UserRepo) *UserUsecase {
	return &UserUsecase{repo: repo}

}

func (u *UserUsecase) Get(ctx context.Context, name string) (*User, error) {
	return u.repo.GetUser(ctx, name)
}
