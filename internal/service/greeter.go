package service

import (
	"context"

	commonv1 "kratosdemo/api/common/v1"
	v1 "kratosdemo/api/helloworld/v1"
	"kratosdemo/internal/biz"
	"kratosdemo/internal/pkg/response"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc *biz.GreeterUsecase
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
	return &GreeterService{uc: uc}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*commonv1.Response, error) {
	g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
	if err != nil {
		// 返回错误响应
		return response.Error(response.CodeError, "创建 Greeter 失败: " + err.Error()), nil
	}
	
	// 创建响应数据
	data := &v1.HelloData{
		Message: "Hello " + g.Hello,
	}
	
	// 返回成功响应
	return response.Success(data), nil
}
