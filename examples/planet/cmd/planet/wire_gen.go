// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/databonfire/bonfire/examples/planet/internal/conf"
	"github.com/databonfire/bonfire/examples/planet/internal/server"
	"github.com/databonfire/bonfire/examples/planet/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(biz *conf.Biz, confServer *conf.Server, data *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	influencerService := service.NewInfluencerService()
	grpcServer := server.NewGRPCServer(confServer, influencerService, logger)
	httpServer := server.NewHTTPServer(confServer, biz, data, influencerService, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
	}, nil
}