package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"push-api-cron/data/repository"
	"push-api-cron/data/service"
	"push-api-cron/server"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	db, err := NewPostgresDb()
	if err != nil {
		logrus.Fatal(err)
	}
	repo := repository.NewRepository(db)
	serv := service.NewService(repo)
	hand := server.NewRouter(serv)
	s := new(Server)
	if err := s.Run("8080", hand.InitRoutes()); err != nil {
		logrus.Fatal(err)
	}
}

func NewPostgresDb() (*sqlx.DB, error) {
	//db, err := sqlx.Open("postgres", fmt.Sprintf("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"))
	//cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password))
	db, err := sqlx.Open("postgres", fmt.Sprintf("postgres://postgres:postgres@192.168.0.104:5432/postgres?sslmode=disable"))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil

}
