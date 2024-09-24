package main

import (
	"github.com/gin-contrib/cors"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"os"
	"os/signal"
	"silkroad/m/internal/config"
	"silkroad/m/internal/delivery/http"
	"silkroad/m/internal/repository"
	"silkroad/m/internal/repository/pg"
	"silkroad/m/internal/server"
	"silkroad/m/internal/service"
	"silkroad/m/internal/telegram"
	"syscall"
)

// @title SilkRoad App API
// @version 1.0
// description API Server for SilkRoad Application
// @host 178.128.123.250:80
// @basePath /api
func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := config.InitConfig("configs", "config"); err != nil {
		logrus.Fatalf("init config err: %s", err.Error())
	}

	db, err := pg.NewPostgresDB(pg.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})
	if err != nil {
		logrus.Fatalf("init db err: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	telegramClient := telegram.NewTelegramClient()
	handlers := http.NewHandler(services, telegramClient)

	srv := new(server.Server)
	go func() {
		router := handlers.InitRoutes()

		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept"},
			AllowCredentials: true,
		}))

		if err := srv.Run(viper.GetString("port"), router); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Print("App Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("App Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Error occured while shutting down http server: %s", err.Error())
	}

	if db != nil {
		db.Close()
	}
}
