package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"silkroad/m/internal/domain/forms"
)

// createHelpWithTourForm Creates a helpWithTour form
// @Summary Create a help with tour form
// @Description This method creates a new helpWithTour form
// @Tags forms
// @Accept  json
// @Produce  json
// @Param input body forms.HelpWithTourForm true "helpWithTour form data"
// @Success 201 {object} map[string]interface{} "Created"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /help_with_tour_form [post]
func (h *Handler) createHelpWithTourForm(c *gin.Context) {
	var input forms.HelpWithTourForm
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.HelpWithTourForm.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	message := fmt.Sprintf(
		"ФОРМА С ПОМОЩЬЮ ВЫБОРА ТУРА\n*ID формы*: %d\n*Имя*: %s\n*Телефон*: %s\n*Место*: %s\n*Дата поездки*: %s",
		id, input.Name, input.Phone, input.Place, input.WhenDate,
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
