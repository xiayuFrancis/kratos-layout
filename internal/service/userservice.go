package service

import (
	"blog/api/common"
	"context"
	"github.com/go-kratos/kratos/v2/errors"

	pb "blog/api/user/v1"
)

type UserServiceService struct {
	pb.UnimplementedUserServiceServer
	response *ResponseWrapper
}

func NewUserServiceService(response *ResponseWrapper) *UserServiceService {
	return &UserServiceService{response: response}
}

func (s *UserServiceService) Login(ctx context.Context, req *pb.LoginRequest) (*common.Response, error) {
	if req.Name == "admin" {
		return s.response.Success(&pb.User{Name: "admin", Password: "123456"})
	}
	return s.response.Error(errors.NotFound(common.ErrorReason_USER_NOT_FOUND.String(), "")), nil

}
