package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"silkroad/m/internal/domain/tour"
	"strconv"
)

// createTour creates a new tour
// @Summary Creates a new tour
// @Description This method creates a new tour with the given input data
// @Tags tours
// @Accept  json
// @Produce  json
// @Param input body tour.Tour true "Tour input"
// @Success 201 {object} map[string]interface{} "Created"
// @Failure 400 {object} errorResponse "Bad Request"
// @Failure 422 {object} errorResponse "Invalid Tour Type"
// @Failure 500 {object} errorResponse "Internal Server Error"
// @Router /tours [post]
func (h *Handler) createTour(c *gin.Context) {
	var input tour.Tour
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if !tour.IsValidTourType(input.TourType) {
		newErrorResponse(c, http.StatusUnprocessableEntity, "Invalid Tour Type. Please choose the correct name: "+
			"Однодневный тур / Многодневный тур / Сити-тур / Эксклюзивный тур / Инфо-тур / Авторский тур")
		return
	}

	id, err := h.services.Tour.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":      id,
		"status":  http.StatusCreated,
		"message": "Created",
	})

}

type getAllToursResponse struct {
	Tours        []tour.Tour `json:"tours"`
	CurrentPage  int         `json:"currentPage"`
	ItemsPerPage int         `json:"itemsPerPage"`
	TotalItems   int         `json:"totalItems"`
	TotalPages   int         `json:"totalPages"`
	TourPlaces   []string    `json:"tourPlaces"`
}

// getAllTour godoc
// @Summary Get all tours
// @Description Получение списка туров с возможностью фильтрации по местоположению, дате, названию, количеству, цене, продолжительности, популярности и другим параметрам
// @Tags tours
// @Accept json
// @Produce json
// @Param tour_place query string false "Tour place"
// @Param quantity query []int false "Quantity (array of integers)"
// @Param priceMin query int false "Minimum price"
// @Param priceMax query int false "Maximum price"
// @Param duration query int false "Duration"
// @Param tour_date query string false "Tour date"
// @Param search query string false "Search by title"
// @Param limit query int false "Limit for pagination" default(4)
// @Param offset query int false "Offset for pagination" default(0)
// @Param popular query bool false "Filter by popular tours"
// @Success 200 {object} getAllToursResponse
// @Failure 500 {object} errorResponse
// @Router /tours [get]
func (h *Handler) getAllTour(c *gin.Context) {
	tourPlace := c.Query("tour_place")

	quantityStr := c.QueryArray("quantity")
	var quantity []int
	for _, q := range quantityStr {
		qInt, err := strconv.Atoi(q)
		if err == nil {
			quantity = append(quantity, qInt)
		}
	}

	priceMin, err := strconv.Atoi(c.Query("priceMin"))
	if err == nil {
		priceMin = 0
	}

	priceMax, err := strconv.Atoi(c.Query("priceMax"))
	if err == nil {
		priceMax = 0
	}

	duration, err := strconv.Atoi(c.Query("duration"))
	if err != nil {
		duration = 0
	}

	tourDate := c.Query("tour_date")
	searchTitle := c.Query("search")

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 4
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = 0
	}

	popularParam := c.DefaultQuery("popular", "false")
	popular := popularParam == "true"

	tours, currentPage, limit, totalItems, totalPages, tourPlaces, err := h.services.Tour.GetAll(tourPlace, tourDate, searchTitle, quantity, priceMin, priceMax, duration, limit, offset, popular)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllToursResponse{
		Tours:        tours,
		CurrentPage:  currentPage,
		ItemsPerPage: limit,
		TotalItems:   totalItems,
		TotalPages:   totalPages,
		TourPlaces:   tourPlaces,
	})
}

// getTour Returns tour by ID
// @Summary Returns tour by ID
// @Description This method returns the details of a specific tour by its ID
// @Tags tours
// @Accept  json
// @Produce  json
// @Param id path int true "Tour ID"
// @Success 200 {object} tour.Tour "Tour details"
// @Failure 500 {object} errorResponse "Internal Server Error"
// @Router /tours/{id} [get]
func (h *Handler) getTourById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	_tour, err := h.services.Tour.GetById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, _tour)
}

// getTourBySlug Returns tour by slug
// @Summary Returns tour details by slug
// @Description This method returns the details of a specific tour based on its slug
// @Tags tours
// @Accept  json
// @Produce  json
// @Param slug path string true "Tour Slug"
// @Success 200 {object} tour.Tour "Tour details"
// @Failure 500 {object} errorResponse "Internal Server Error"
// @Router /tours/title/{slug} [get]
func (h *Handler) getTourBySlug(c *gin.Context) {
	_slug := c.Param("slug")
	_tour, err := h.services.Tour.GetBySlug(_slug)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, _tour)
}

// getMinMaxPrice Returns the minimum and maximum tour prices
// @Summary Returns the minimum and maximum prices of all tours
// @Description This method returns the minimum and maximum prices of all available tours
// @Tags tours
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{} "Min and Max tour prices"
// @Failure 500 {object} errorResponse "Internal Server Error"
// @Router /tours/prices [get]
func (h *Handler) getMinMaxPrice(c *gin.Context) {
	minPrice, maxPrice, err := h.services.Tour.GetMinMaxPrice()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"min_price": minPrice,
		"max_price": maxPrice,
	})
}
