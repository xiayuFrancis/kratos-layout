package service

import (
	"blog/api/common"
	"blog/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/errors"

	pb "blog/api/user/v1"
)

type UserServiceService struct {
	pb.UnimplementedUserServiceServer
	response *ResponseWrapper
	user     *biz.UserUsecase
}

func NewUserServiceService(response *ResponseWrapper, user *biz.UserUsecase) *UserServiceService {
	return &UserServiceService{response: response, user: user}
}

func (s *UserServiceService) Login(ctx context.Context, req *pb.LoginRequest) (*common.Response, error) {

	user, err := s.user.Get(ctx, req.Name)
	if err != nil {
		return s.response.Error(errors.NotFound(common.ErrorReason_USER_NOT_FOUND.String(), "")), nil
	}

	return s.response.Success(&user)
}
