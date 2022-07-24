package main

import (
	"github.com/sirupsen/logrus"
	"push-api-cron/server"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	hand := server.NewRouter()
	serv := new(Server)
	if err := serv.Run("8080", hand.InitRoutes()); err != nil {
		logrus.Fatal(err)
	}
}
