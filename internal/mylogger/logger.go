package mylogger

import (
	"github.com/mamoru777/foundation2/loginit"
	"github.com/sirupsen/logrus"

	"os"
)

var Logger *logrus.Logger

func Init() *os.File {
	Logger = loginit.LogInit("tcp", "localhost:5044", "ChatService")
	_ = Logger
	file, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Logger.Out = file
	} else {
		logrus.Error("Не удалось логировать в файл, использую логирование в консоль")
	}
	return file
}
