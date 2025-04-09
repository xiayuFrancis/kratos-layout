// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"kratosdemo/internal/biz"
	"kratosdemo/internal/conf"
	"kratosdemo/internal/data"
	"kratosdemo/internal/server"
	"kratosdemo/internal/service"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	connectRepo := data.NewConnectRepo(dataData, logger)
	connectUsecase := biz.NewConnectUsecase(connectRepo, logger)
	connectService := service.NewConnectService(connectUsecase, logger)
	grpcServer := server.NewGRPCServer(confServer, connectService, logger)
	httpServer := server.NewHTTPServer(confServer, connectService, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
