// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"trainings-go/internal/biz"
	"trainings-go/internal/conf"
	"trainings-go/internal/data"
	"trainings-go/internal/server"
	"trainings-go/internal/service"
)

// Injectors from wire.go:

// initApp init kratos application.
func initApp(confServer *conf.Server, confData *conf.Data, logger log.Logger, auth *conf.Auth) (*kratos.App, func(), error) {
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	userRepo := data.NewUserRepo(dataData, logger)
	userBiz := biz.NewUserBiz(userRepo, logger, auth)
	userService := service.NewUserService(userBiz, logger)
	categoryRepo := data.NewCategoryRepo(logger, dataData)
	categoryBiz := biz.NewCategoryBiz(logger, categoryRepo, userBiz)
	categorySrv := service.NewCategorySrv(logger, categoryBiz)
	actionRepo := data.NewActionRepo(logger, dataData)
	actionBiz := biz.NewActionBiz(logger, actionRepo, userBiz)
	actionService := service.NewActionService(actionBiz, logger)
	httpServer := server.NewHTTPServer(confServer, userService, logger, auth, categorySrv, actionService)
	grpcServer := server.NewGRPCServer(confServer, userService, logger, categorySrv, actionService)
	app := newApp(logger, httpServer, grpcServer)
	return app, func() {
		cleanup()
	}, nil
}
