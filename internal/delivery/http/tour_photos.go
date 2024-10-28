package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// uploadTourPhotos godoc
// @Summary Upload photos for a tour
// @Description Uploads multiple photos for a specific tour by tourID.
// @Tags tours
// @Accept multipart/form-data
// @Produce json
// @Param tourID path int true "Tour ID"
// @Param photos formData file true "Photos to upload" multiple
// @Success 200 {object} map[string]interface{} "status: OK, message: Photos uploaded successfully"
// @Failure 400 {object} map[string]interface{} "status: Bad Request, message: Invalid tourID or form data"
// @Failure 500 {object} map[string]interface{} "status: Internal Server Error, message: Error message"
// @Router /tours/photos/{tourID} [post]
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
