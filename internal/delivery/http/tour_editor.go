package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"silkroad/m/internal/domain/tour"
)

// createTourEditor Creates a custom tour
// @Summary Create a custom tour
// @Description This method allows users to create a custom tour by submitting their details and tour preferences. The tour data will also be sent to Telegram.
// @Tags tour
// @Accept  json
// @Produce  json
// @Param input body tour.TourEditor true "Custom tour data"
// @Success 200 {object} map[string]interface{} "ID and creation status of the tour"
// @Failure 400 {object} errorResponse "Bad request"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /tour_editor [post]
func (h *Handler) createTourEditor(c *gin.Context) {
	var input tour.TourEditor
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.TourEditor.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	message := fmt.Sprintf(
		"СОЗДАЙ ЭКСКУРСИЮ САМ\n*ID формы*: %d\n*Имя*: %s\n*Телефон*: %s\n*Электронная почта*: %s\n"+
			"*Дата поездки*: %s\n*Активности*: %s\n*Локации*: %s",
		id, input.Name, input.Phone, input.Email, input.TourDate, input.Activity, input.Location,
	)

	err = h.telegramClient.SendTelegramMessage(message)
	if err != nil {
		logrus.Errorf("Error sending message in Telegram: %v", err)
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":      id,
		"status":  http.StatusCreated,
		"message": "Created",
	})
}
