package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"silkroad/m/internal/domain/forms"
)

// createContactForm Creates a contact form
// @Summary Create a contact form
// @Description This method creates a new contact form and, if TourID is provided, sends tour details
// @Tags forms
// @Accept  json
// @Produce  json
// @Param input body forms.ContactForm true "Contact form data"
// @Success 201 {object} map[string]interface{} "Created"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /form/contact_form [post]
func (h *Handler) createContactForm(c *gin.Context) {
	var input forms.ContactForm
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.ContactForm.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	message := fmt.Sprintf(
		"ФОРМА ОБРАТНОЙ СВЯЗИ\n*ID формы*: %d\n*Имя*: %s\n*Телефон*: %s\n*Электронная почта*: %s\n*Текст*: %s",
		id, input.Name, input.Phone, input.Email, input.Description,
	)

	if input.TourID != nil {
		tour, err := h.services.Tour.GetById(*input.TourID)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		tourInfo := fmt.Sprintf(
			"\n\n*ID тура*: %d\n*Тип тура*: %s\n*Название тура*: %s\n*Место тура*: %s\n*Дата тура*: %s"+
				"\n*Количество человек*: %d\n*Активности*: %s\n*Тариф*: %s",
			*input.TourID, tour.TourType, tour.Title, tour.TourPlace, tour.TourDate.Format("02-01-2006"), tour.Quantity, tour.Activity, tour.Tariff)
		message += tourInfo
	}

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
