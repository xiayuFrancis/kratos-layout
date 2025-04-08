package data

import (
	"context"
	"kratosdemo/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// ConnectRepo 是数据库连接测试的仓库实现
type ConnectRepo struct {
	data *Data
	log  *log.Helper
}

// NewConnectRepo 创建一个新的ConnectRepo实例
func NewConnectRepo(data *Data, logger log.Logger) biz.ConnectRepo {
	return &ConnectRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// TestConnect 测试数据库连接
func (r *ConnectRepo) TestConnect(ctx context.Context) (bool, error) {
	// 使用Ent客户端测试连接
	if r.data.db == nil {
		return false, nil
	}
	
	// 尝试执行一个简单的查询来测试连接
	_, err := r.data.db.User.Query().Count(ctx)
	if err != nil {
		r.log.Errorf("database connection test failed: %v", err)
		return false, err
	}
	
	return true, nil
}
