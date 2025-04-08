package data

import (
	"context"
	"kratosdemo/ent"
	"kratosdemo/internal/conf"
	"time"

	"entgo.io/ent/dialect"
	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewConnectRepo, NewRedisClient)

// Data .
type Data struct {
	db    *ent.Client
	redis *redis.Client
	log   *log.Helper
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	log := log.NewHelper(logger)
	// 初始化 MySQL 连接
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

	// 初始化 Redis 连接
	rdb := redis.NewClient(&redis.Options{
		Network:  c.Redis.Network,
		Addr:     c.Redis.Addr,
		Password: "", // 如果有密码，可以从配置中获取
		DB:       0,  // 默认使用 0 号数据库
	})

	// 测试 Redis 连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Errorf("failed connecting to redis: %v", err)
		return nil, nil, err
	}

	d := &Data{
		db:    client,
		redis: rdb,
		log:   log,
	}

	cleanup := func() {
		log.Info("closing the data resources")
		if err := d.db.Close(); err != nil {
			log.Error(err)
		}
		if err := d.redis.Close(); err != nil {
			log.Error(err)
		}
	}
	return d, cleanup, nil
}
