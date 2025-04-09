package service

import (
	"context"
	commonv1 "kratosdemo/api/common/v1"
	v1 "kratosdemo/api/connect/v1"
	"kratosdemo/internal/biz"
	"kratosdemo/internal/pkg/response"

	"github.com/go-kratos/kratos/v2/log"
)

// ConnectService 是Connect服务的实现
type ConnectService struct {
	v1.UnimplementedConnectServer

	uc  *biz.ConnectUsecase
	log *log.Helper
}

// NewConnectService 创建一个新的ConnectService
func NewConnectService(uc *biz.ConnectUsecase, logger log.Logger) *ConnectService {
	return &ConnectService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// TestConnect 实现TestConnect接口
func (s *ConnectService) TestConnect(ctx context.Context, req *v1.TestConnectRequest) (*commonv1.Response, error) {
	logger := s.log.WithContext(ctx)
	logger.Info("Received TestConnect request")
	success, err := s.uc.TestConnect(ctx)
	if err != nil {
		// 创建响应数据
		data := &v1.TestConnectData{
			Message: "Database connection test failed: " + err.Error(),
			Success: false,
		}
		// 返回错误响应
		return response.ErrorWithData(response.CodeError, "Database connection test failed", data), nil
	}

	message := "Database connection test successful"
	if !success {
		message = "Database connection test failed"
	}

	// 创建响应数据
	data := &v1.TestConnectData{
		Message: message,
		Success: success,
	}

	// 返回成功响应
	return response.Success(data), nil
}
