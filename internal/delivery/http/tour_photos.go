package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) uploadTourPhotos(c *gin.Context) {
	tourID, err := strconv.Atoi(c.Param("tourID"))
	if err != nil || tourID < 1 {
		newErrorResponse(c, http.StatusBadRequest, "Invalid tourID")
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid form data")
		return
	}
	files := form.File["photos"]

	err = h.services.Tour.AddPhotos(tourID, files)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status":  http.StatusOK,
		"message": "Photos uploaded successfully",
	})
}
