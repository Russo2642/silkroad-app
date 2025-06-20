package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"silkroad/m/internal/aws"
	"silkroad/m/internal/config"
	"silkroad/m/internal/delivery/http"
	"silkroad/m/internal/repository"
	"silkroad/m/internal/repository/pg"
	"silkroad/m/internal/server"
	"silkroad/m/internal/service"
	"silkroad/m/internal/telegram"

	"github.com/gin-contrib/cors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type App struct {
	cfg    *config.Config
	server *server.Server
	db     *sqlx.DB
}

func New() *App {
	return &App{}
}

func (a *App) Run() error {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := godotenv.Load(); err != nil {
		logrus.Printf("Info: .env file not found, using default values: %s", err.Error())
	}

	a.cfg = config.LoadConfig()

	if err := aws.InitS3Client(a.cfg.AWS); err != nil {
		return fmt.Errorf("failed to initialize S3 client: %w", err)
	}

	var err error
	a.db, err = pg.NewPostgresDB(pg.Config{
		Host:     a.cfg.DB.Host,
		Port:     a.cfg.DB.Port,
		Username: a.cfg.DB.Username,
		Password: a.cfg.DB.Password,
		DBName:   a.cfg.DB.DBName,
		SSLMode:  a.cfg.DB.SSLMode,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := a.runMigrations(); err != nil {
		return fmt.Errorf("migration error: %w", err)
	}

	repos := repository.NewRepository(a.db)
	services := service.NewService(repos)
	telegramClient := telegram.NewTelegramClient(a.cfg.Telegram)
	handlers := http.NewHandler(services, telegramClient)

	a.server = server.New(a.cfg.Server)

	go func() {
		router := handlers.InitRoutes()

		router.Use(cors.New(cors.Config{
			AllowOrigins:     a.cfg.CORS.AllowedOrigins,
			AllowMethods:     a.cfg.CORS.AllowedMethods,
			AllowHeaders:     a.cfg.CORS.AllowedHeaders,
			AllowCredentials: a.cfg.CORS.AllowCredentials,
		}))

		if err := a.server.Run(a.cfg.HttpPort, router); err != nil {
			logrus.Fatalf("error occurred while running http server: %s", err.Error())
		}
	}()

	logrus.Print("App Started")

	a.waitForShutdown()

	return nil
}

func (a *App) runMigrations() error {
	driver, err := postgres.WithInstance(a.db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create postgres driver: %w", err)
	}

	migrationsPath := "file://db/migrations"
	if _, err := os.Stat("db/migrations"); os.IsNotExist(err) {
		migrationsPath = "file:///app/db/migrations"
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run up migrations: %w", err)
	}

	logrus.Info("Migrations applied successfully")
	return nil
}

func (a *App) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("App Shutting Down")

	if err := a.server.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Error occurred while shutting down http server: %s", err.Error())
	}

	if a.db != nil {
		if err := a.db.Close(); err != nil {
			logrus.Errorf("Error occurred while closing database connection: %s", err.Error())
		}
	}
}
