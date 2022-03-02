//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"trainings-go/internal/biz"
	"trainings-go/internal/conf"
	"trainings-go/internal/data"
	"trainings-go/internal/server"
	"trainings-go/internal/service"
)

// initApp init kratos application.
func initApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, biz.ProviderSet, data.ProviderSet, service.ProviderSet, newApp))
}
