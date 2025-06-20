package http

import (
	"net/http"
	"strconv"

	"silkroad/m/internal/domain/tour"

	"github.com/gin-gonic/gin"
)

// @Summary Upload photos for a tour (new version)
// @Description Uploads multiple photos for a specific tour by tourID with advanced metadata support
// @Tags tours
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Tour ID"
// @Param photoType formData string true "Type of photo (preview, gallery, route, booking)"
// @Param photos formData file true "Photos to upload" multiple
// @Param title formData string false "Photo title"
// @Param description formData string false "Photo description"
// @Param alt_text formData string false "Photo alt text"
// @Param display_order formData int false "Display order"
// @Success 200 {object} map[string]interface{} "status: OK, message: Photos uploaded successfully"
// @Failure 400 {object} map[string]interface{} "status: Bad Request, message: Invalid parameters"
// @Failure 500 {object} map[string]interface{} "status: Internal Server Error, message: Error message"
// @Router /tours/{id}/photos/upload [post]
func (h *Handler) uploadTourPhotosNew(c *gin.Context) {
	tourID, err := strconv.Atoi(c.Param("id"))
	if err != nil || tourID < 1 {
		newErrorResponse(c, http.StatusBadRequest, "Invalid tourID")
		return
	}

	photoTypeStr := c.PostForm("photoType")
	if photoTypeStr == "" {
		newErrorResponse(c, http.StatusBadRequest, "Photo type is required")
		return
	}

	photoType := tour.TourPhotoType(photoTypeStr)
	if !tour.IsValidPhotoType(photoType) {
		newErrorResponse(c, http.StatusBadRequest, "Invalid photo type")
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid form data")
		return
	}

	files := form.File["photos"]
	if len(files) == 0 {
		newErrorResponse(c, http.StatusBadRequest, "No photos provided")
		return
	}

	metadata := tour.TourPhotoInput{
		TourID:    tourID,
		PhotoType: photoType,
	}

	if title := c.PostForm("title"); title != "" {
		metadata.Title = &title
	}

	if description := c.PostForm("description"); description != "" {
		metadata.Description = &description
	}

	if altText := c.PostForm("alt_text"); altText != "" {
		metadata.AltText = &altText
	}

	if displayOrderStr := c.PostForm("display_order"); displayOrderStr != "" {
		if displayOrder, err := strconv.Atoi(displayOrderStr); err == nil {
			metadata.DisplayOrder = displayOrder
		}
	}

	err = h.services.TourPhoto.UploadPhotos(tourID, files, photoType, metadata)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Photos uploaded successfully",
	})
}

// @Summary Get photos for a tour
// @Description Get all photos for a specific tour grouped by type
// @Tags tours
// @Accept json
// @Produce json
// @Param id path int true "Tour ID"
// @Success 200 {object} tour.TourPhotosGrouped "Tour photos"
// @Failure 400 {object} map[string]interface{} "status: Bad Request, message: Invalid tourID"
// @Failure 404 {object} map[string]interface{} "status: Not Found, message: Tour not found"
// @Failure 500 {object} map[string]interface{} "status: Internal Server Error, message: Error message"
// @Router /tours/{id}/photos [get]
func (h *Handler) getTourPhotos(c *gin.Context) {
	tourID, err := strconv.Atoi(c.Param("id"))
	if err != nil || tourID < 1 {
		newErrorResponse(c, http.StatusBadRequest, "Invalid tourID")
		return
	}

	photos, err := h.services.TourPhoto.GetByTourID(tourID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, photos)
}

// @Summary Get photos by filter
// @Description Get photos with filtering and pagination
// @Tags photos
// @Accept json
// @Produce json
// @Param tour_id query int false "Tour ID"
// @Param photo_type query string false "Photo type" Enums(preview, gallery, route, booking)
// @Param is_active query bool false "Is active"
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{} "photos: array of photos, total: total count"
// @Failure 400 {object} map[string]interface{} "status: Bad Request, message: Invalid parameters"
// @Failure 500 {object} map[string]interface{} "status: Internal Server Error, message: Error message"
// @Router /photos [get]
func (h *Handler) getPhotosByFilter(c *gin.Context) {
	var filter tour.TourPhotoFilter

	if err := c.ShouldBindQuery(&filter); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid filter parameters")
		return
	}

	photos, total, err := h.services.TourPhoto.GetByFilter(filter)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"photos": photos,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// @Summary Update tour photo
// @Description Update photo metadata
// @Tags photos
// @Accept json
// @Produce json
// @Param photoID path int true "Photo ID"
// @Param input body tour.TourPhotoInput true "Photo update data"
// @Success 200 {object} map[string]interface{} "status: success, message: Photo updated successfully"
// @Failure 400 {object} map[string]interface{} "status: Bad Request, message: Invalid parameters"
// @Failure 500 {object} map[string]interface{} "status: Internal Server Error, message: Error message"
// @Router /photos/{photoID} [put]
func (h *Handler) updateTourPhoto(c *gin.Context) {
	photoID, err := strconv.Atoi(c.Param("photoID"))
	if err != nil || photoID < 1 {
		newErrorResponse(c, http.StatusBadRequest, "Invalid photoID")
		return
	}

	var input tour.TourPhotoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.TourPhoto.Update(photoID, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Photo updated successfully",
	})
}

// @Summary Delete tour photo
// @Description Delete photo (soft delete)
// @Tags photos
// @Accept json
// @Produce json
// @Param photoID path int true "Photo ID"
// @Success 200 {object} map[string]interface{} "status: success, message: Photo deleted successfully"
// @Failure 400 {object} map[string]interface{} "status: Bad Request, message: Invalid photoID"
// @Failure 500 {object} map[string]interface{} "status: Internal Server Error, message: Error message"
// @Router /photos/{photoID} [delete]
func (h *Handler) deleteTourPhoto(c *gin.Context) {
	photoID, err := strconv.Atoi(c.Param("photoID"))
	if err != nil || photoID < 1 {
		newErrorResponse(c, http.StatusBadRequest, "Invalid photoID")
		return
	}

	err = h.services.TourPhoto.Delete(photoID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Photo deleted successfully",
	})
}
