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
		form := api.Group("/form")
		{
			form.POST("/contact-form", h.createContactForm)
			form.POST("/help-with-tour-form", h.createHelpWithTourForm)
		}

		getCountries := api.Group("/countries")
		{
			getCountries.GET("/", h.getCountries)
		}

		tour := api.Group("/tours")
		{
			tour.POST("/", h.createTour)
			tour.GET("/", h.getAllTour)
			tour.GET("/:id", h.getTourById)
			tour.GET("/title/:slug", h.getTourBySlug)
			tour.GET("/prices", h.getMinMaxPrice)
			tour.POST("/photos/:tourID", h.uploadTourPhotos)
		}

		tourEditor := api.Group("/tour_editor")
		{
			tourEditor.POST("/", h.createTourEditor)
		}
	}

	return router
}
