package main

import (
	"silkroad/m/internal/app"

	"github.com/sirupsen/logrus"
)

// @title SilkRoad App API
// @version 1.0
// description API Server for SilkRoad Application

// @basePath /api

func main() {
	application := app.New()

	if err := application.Run(); err != nil {
		logrus.Fatalf("Failed to run application: %s", err.Error())
	}
}
