package data

import (
	"context"
	"kratosdemo/ent"
	"kratosdemo/internal/conf"

	"entgo.io/ent/dialect"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	_ "github.com/go-sql-driver/mysql"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewConnectRepo)

// Data .
type Data struct {
	db  *ent.Client
	log *log.Helper
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	log := log.NewHelper(logger)
	client, err := ent.Open(
		dialect.MySQL,
		c.Database.Source,
	)
	if err != nil {
		log.Errorf("failed opening connection to mysql: %v", err)
		return nil, nil, err
	}
	// 运行数据库自动迁移工具
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Errorf("failed creating schema resources: %v", err)
		return nil, nil, err
	}

	d := &Data{
		db:  client,
		log: log,
	}

	cleanup := func() {
		log.Info("closing the data resources")
		if err := d.db.Close(); err != nil {
			log.Error(err)
		}
	}
	return d, cleanup, nil
}
