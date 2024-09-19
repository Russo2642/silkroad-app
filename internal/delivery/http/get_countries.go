package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// getCountries Returns a list of countries
// @Summary Returns a list of countries
// @Description This method returns a list of all available countries
// @Tags countries
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{} "List of countries"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /countries [get]
func (h *Handler) getCountries(c *gin.Context) {
	countries, err := h.services.GetCountries.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"countries": countries,
	})
}
