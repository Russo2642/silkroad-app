package http

import (
	_ "silkroad/m/docs"
	"silkroad/m/internal/service"
	"silkroad/m/internal/telegram"

	"github.com/gin-gonic/gin"

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
			form.POST("/contact_form", h.createContactForm)
			form.POST("/help_with_tour_form", h.createHelpWithTourForm)
		}

		getCountries := api.Group("/countries")
		{
			getCountries.GET("/", h.getCountries)
		}

		tour := api.Group("/tours")
		{
			tour.POST("/", h.createTour)
			tour.GET("/", h.getAllTour)
			tour.GET("/prices", h.getMinMaxPrice)
			tour.GET("/:id", h.getTourById)
			tour.GET("/title/:slug", h.getTourBySlug)
			tour.POST("/:id/photos/upload", h.uploadTourPhotosNew)
			tour.GET("/:id/photos", h.getTourPhotos)
		}

		photos := api.Group("/photos")
		{
			photos.GET("/", h.getPhotosByFilter)
			photos.PUT("/:photoID", h.updateTourPhoto)
			photos.DELETE("/:photoID", h.deleteTourPhoto)
		}
	}

	return router
}
