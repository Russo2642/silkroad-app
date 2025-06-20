package http

import (
	"fmt"
	"net/http"
	"silkroad/m/internal/domain/forms"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Summary Create a help with tour form
// @Description This method creates a new helpWithTour form with name, phone, country and when_date
// @Tags forms
// @Accept  json
// @Produce  json
// @Param input body forms.HelpWithTourForm true "helpWithTour form data"
// @Success 201 {object} map[string]interface{} "Created"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /form/help_with_tour_form [post]
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
		"ФОРМА ПОМОЩИ С ВЫБОРОМ ТУРА\n*ID формы*: %d\n*Имя*: %s\n*Телефон*: %s\n*Где будем отдыхать*: %s\n*Когда*: %s",
		id, input.Name, input.Phone, input.Country, input.WhenDate.Format("02.01.2006"),
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
