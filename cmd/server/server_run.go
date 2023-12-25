package main

import (
	"github.com/caarlos0/env/v8"
	"github.com/mamoru777/chatservice2/internal/app"
	"github.com/mamoru777/chatservice2/internal/config"
	"github.com/mamoru777/chatservice2/internal/mylogger"
)

func main() {
	file := mylogger.Init()
	defer file.Close()
	dataBaseconfig := config.DataBaseConfig{}
	grpcServerConfig := config.GrpcServerConfig{}
	if err := env.Parse(&dataBaseconfig); err != nil {
		mylogger.Logger.Fatalf("ошибка при получении переменных окружения, %v", err)
	}
	if err := env.Parse(&grpcServerConfig); err != nil {
		mylogger.Logger.Fatalf("ошибка при получении переменных окружения, %v", err)
	}
	if err := app.Run(dataBaseconfig, grpcServerConfig); err != nil {
		mylogger.Logger.Fatal("ошибка при запуске сервера ", err)
	}
}
