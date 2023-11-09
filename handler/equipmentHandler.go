package handler

import (
	"miniproject/config"
	"miniproject/entity"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetAllEquipments(c echo.Context) error {
	var equipments []entity.Equipment
	if err := config.DB.Table("equipments").Find(&equipments).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Failed to fetch equipments"})
	}

	return c.JSON(http.StatusOK, equipments)
}
