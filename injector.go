//go:build wireinject
// +build wireinject

package main

import (
	"khaira-admin/controller"
	"khaira-admin/helper"
	"khaira-admin/repository"
	"khaira-admin/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

var ServerSet = wire.NewSet(
	helper.NewElasticClient,
	repository.NewRepositoryImpl,
	service.NewServiceImpl,
	controller.NewControllerImpl,
	helper.NewDb,
	NewServer,
)

func InitServer() (*fiber.App, func(), error) {
	wire.Build(ServerSet)
	return nil, nil, nil
}
