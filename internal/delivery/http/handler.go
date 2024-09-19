package http

import (
	"github.com/gin-gonic/gin"
	_ "silkroad/m/docs"
	"silkroad/m/internal/service"
	"silkroad/m/internal/telegram"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services       *service.Service
	telegramClient *telegram.Client
}

func NewHandler(services *service.Service, telegramClient *telegram.Client) *Handler {
	return &Handler{
		services:       services,
		telegramClient: telegramClient,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		contactForm := api.Group("/contact_form")
		{
			contactForm.POST("/", h.createContactForm)
		}

		helpWithTourForm := api.Group("/help_form")
		{
			helpWithTourForm.POST("/", h.createHelpWithTourForm)
		}

		getCountries := api.Group("/countries")
		{
			getCountries.GET("/", h.getCountries)
		}

		tour := api.Group("/tour")
		{
			tour.POST("/", h.createTour)
			tour.GET("/", h.getAllTour)
			tour.GET("/:id", h.getTour)
		}

		tourEditor := api.Group("/tour_editor")
		{
			tourEditor.POST("/", h.createTourEditor)
		}
	}

	return router
}
