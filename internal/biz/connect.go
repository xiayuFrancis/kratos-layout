package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// ConnectRepo 是数据库连接测试的仓库接口
type ConnectRepo interface {
	TestConnect(ctx context.Context) (bool, error)
}

// ConnectUsecase 是数据库连接测试的用例
type ConnectUsecase struct {
	repo ConnectRepo
	log  *log.Helper
}

// NewConnectUsecase 创建一个新的ConnectUsecase
func NewConnectUsecase(repo ConnectRepo, logger log.Logger) *ConnectUsecase {
	return &ConnectUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// TestConnect 测试数据库连接
func (uc *ConnectUsecase) TestConnect(ctx context.Context) (bool, error) {
	return uc.repo.TestConnect(ctx)
}
