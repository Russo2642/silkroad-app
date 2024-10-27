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
// @Tags tour
// @Accept  json
// @Produce  json
// @Param input body tour.Tour true "Tour input"
// @Success 201 {object} map[string]interface{} "Created"
// @Failure 400 {object} errorResponse "Bad Request"
// @Failure 422 {object} errorResponse "Invalid Tour Type"
// @Failure 500 {object} errorResponse "Internal Server Error"
// @Router /tour [post]
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

// getAllTour Returns a list of tours with optional filters
// @Summary Returns a list of tours
// @Description This method returns a list of tours. You can filter by price range, tour place, tour date, quantity, duration, and search by title. Pagination is also supported via limit and offset.
// @Tags tour
// @Accept  json
// @Produce  json
// @Param tour_place query string false "Filter by tour place"
// @Param quantity query int false "Filter by quantity of people"
// @Param priceMin query int false "Filter by price range (from min), example: 100"
// @Param priceMax query int false "Filter by price range (to max), example: 500"
// @Param duration query int false "Filter by duration of the tour"
// @Param tour_date query string false "Filter by date of the tour, format: YYYY-MM-DDT00:00:00+00:00"
// @Param search query string false "Search tours by title"
// @Param limit query int false "Limit the number of returned tours"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} getAllToursResponse "List of tours / CurrentPage / ItemsPerPage / TotalItems / TotalPages / TourPlaces"
// @Failure 500 {object} errorResponse "Internal Server Error"
// @Router /tour [get]
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

	tours, currentPage, limit, totalItems, totalPages, tourPlaces, err := h.services.Tour.GetAll(tourPlace, tourDate, searchTitle, quantity, priceMin, priceMax, duration, limit, offset)
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
// @Tags tour
// @Accept  json
// @Produce  json
// @Param id path int true "Tour ID"
// @Success 200 {object} tour.Tour "Tour details"
// @Failure 500 {object} errorResponse "Internal Server Error"
// @Router /tour/{id} [get]
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
// @Tags tour
// @Accept  json
// @Produce  json
// @Param slug path string true "Tour Slug"
// @Success 200 {object} tour.Tour "Tour details"
// @Failure 500 {object} errorResponse "Internal Server Error"
// @Router /tour/title/{slug} [get]
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
// @Tags tour
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{} "Min and Max tour prices"
// @Failure 500 {object} errorResponse "Internal Server Error"
// @Router /tour/prices [get]
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
