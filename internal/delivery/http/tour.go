package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"silkroad/m/internal/domain/tour"
	"strconv"
)

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
	Tours []tour.Tour `json:"tours"`
}

// getAllTour Returns a list of tours with optional filters
// @Summary Returns a list of tours
// @Description This method returns a list of tours. You can filter by price range, tour place, tour date, quantity, duration, and search by title. Pagination is also supported via limit and offset.
// @Tags tour
// @Accept  json
// @Produce  json
// @Param priceRange query string false "Filter by price range, example: '100-500'"
// @Param tour_place query string false "Filter by tour place"
// @Param quantity query int false "Filter by quantity of people"
// @Param duration query int false "Filter by duration of the tour"
// @Param tour_date query string false "Filter by date of the tour, format: YYYY-MM-DDT00:00:00+00:00"
// @Param search query string false "Search tours by title"
// @Param limit query int false "Limit the number of returned tours"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} getAllToursResponse "List of tours"
// @Failure 500 {object} errorResponse "Internal Server Error"
// @Router /tour [get]
func (h *Handler) getAllTour(c *gin.Context) {
	priceRange := c.Query("priceRange")
	tourPlace := c.Query("tour_place")

	quantity, err := strconv.Atoi(c.Query("quantity"))
	if err != nil {
		quantity = 0
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

	tours, err := h.services.Tour.GetAll(priceRange, tourPlace, tourDate, searchTitle, quantity, duration, limit, offset)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllToursResponse{
		Tours: tours,
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
func (h *Handler) getTour(c *gin.Context) {
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
