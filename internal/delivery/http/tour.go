package http

import (
	"net/http"
	"silkroad/m/internal/domain/tour"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

	if !tour.IsValidTourType(input.Type) {
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

// @Description Получение списка туров с возможностью фильтрации по местоположению, количеству участников, цене, продолжительности, популярности и другим параметрам
// @Tags tours
// @Accept json
// @Produce json
// @Param tour_place query string false "Tour place (country)"
// @Param quantity query int false "Number of participants - filters tours that can accommodate this number of people"
// @Param type query []string false "Tour type (array of strings)"
// @Param difficulty query []int false "Difficulty level (array of integers from 1 to 5)"
// @Param activities query []string false "Activities (array of strings)"
// @Param categories query []string false "Categories (array of strings)"
// @Param priceMin query int false "Minimum price"
// @Param priceMax query int false "Maximum price"
// @Param duration query int false "Duration in days"
// @Param search query string false "Search by title and description"
// @Param limit query int false "Limit for pagination" default(4)
// @Param offset query int false "Offset for pagination" default(0)
// @Param popular query bool false "Filter by popular tours"
// @Success 200 {object} getAllToursResponse
// @Failure 500 {object} errorResponse
// @Router /tours [get]
func (h *Handler) getAllTour(c *gin.Context) {
	tourPlace := c.Query("tour_place")

	var quantity *int
	if quantityStr := c.Query("quantity"); quantityStr != "" {
		if val, err := strconv.Atoi(quantityStr); err == nil && val > 0 {
			quantity = &val
		}
	}

	var priceMin, priceMax *int
	if priceMinStr := c.Query("priceMin"); priceMinStr != "" {
		if val, err := strconv.Atoi(priceMinStr); err == nil && val > 0 {
			priceMin = &val
		}
	}
	if priceMaxStr := c.Query("priceMax"); priceMaxStr != "" {
		if val, err := strconv.Atoi(priceMaxStr); err == nil && val > 0 {
			priceMax = &val
		}
	}

	var durationFilter *tour.RangeFilter
	if durationStr := c.Query("duration"); durationStr != "" {
		if duration, err := strconv.Atoi(durationStr); err == nil && duration > 0 {
			durationFilter = &tour.RangeFilter{Min: &duration, Max: &duration}
		}
	}

	tourTypes := c.QueryArray("type")
	var types []tour.TourType
	for _, t := range tourTypes {
		if tour.IsValidTourType(tour.TourType(t)) {
			types = append(types, tour.TourType(t))
		}
	}

	difficultyStr := c.QueryArray("difficulty")
	var difficulties []tour.Difficulty
	for _, d := range difficultyStr {
		if diffInt, err := strconv.Atoi(d); err == nil {
			if tour.IsValidDifficulty(tour.Difficulty(diffInt)) {
				difficulties = append(difficulties, tour.Difficulty(diffInt))
			}
		}
	}

	activities := c.QueryArray("activities")
	categories := c.QueryArray("categories")

	searchTitle := c.Query("search")

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		limit = 4
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	var popular *bool
	if popularParam := c.Query("popular"); popularParam != "" {
		if popularParam == "true" {
			val := true
			popular = &val
		} else if popularParam == "false" {
			val := false
			popular = &val
		}
	}

	filter := tour.TourFilter{
		SearchQuery: searchTitle,
		Limit:       limit,
		Offset:      offset,
		PriceMin:    priceMin,
		PriceMax:    priceMax,
		Duration:    durationFilter,
		Popular:     popular,
		Quantity:    quantity,
		Type:        types,
		Difficulty:  difficulties,
		Activities:  activities,
		Categories:  categories,
	}

	if tourPlace != "" {
		filter.Country = []string{tourPlace}
	}

	tours, totalItems, err := h.services.Tour.GetAll(filter)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	currentPage := (offset / limit) + 1
	totalPages := (totalItems + limit - 1) / limit

	c.JSON(http.StatusOK, getAllToursResponse{
		Tours:        tours,
		CurrentPage:  currentPage,
		ItemsPerPage: limit,
		TotalItems:   totalItems,
		TotalPages:   totalPages,
		TourPlaces:   []string{}, // TODO: implement if needed
	})
}

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

	_tour, err := h.services.Tour.GetByID(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, _tour)
}

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
