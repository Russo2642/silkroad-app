package http

import (
	"net/http"

	"silkroad/m/internal/domain/country"

	"github.com/gin-gonic/gin"
)

// @Summary Get countries
// @Description Get list of countries with optional filtering
// @Tags countries
// @Accept json
// @Produce json
// @Param is_active query bool false "Filter by active status"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} errorResponse
// @Router /countries [get]
func (h *Handler) getCountries(c *gin.Context) {
	var filter country.CountryFilter

	if err := c.ShouldBindQuery(&filter); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	countries, err := h.services.Country.GetAll(filter)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": countries,
	})
}
