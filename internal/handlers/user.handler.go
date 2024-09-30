package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/KevinStockmanns/api_golang/internal/db"
	"github.com/KevinStockmanns/api_golang/internal/dtos"
	"github.com/KevinStockmanns/api_golang/internal/models"
	"github.com/KevinStockmanns/api_golang/internal/validators"
	"github.com/labstack/echo/v4"
)

func UserPost(c echo.Context) error {
	var userDto dtos.UserPostDTO

	if err := c.Bind(&userDto); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "no se pudo leer el cuerpo enviado",
		})
	}
	userDto.Normalize()
	if errors, ok := validators.ValidateDTOs(userDto); !ok {
		return c.JSON(http.StatusBadRequest, errors)
	}
	bDay, _ := time.Parse("2006-01-02", userDto.Birthday)
	user := models.User{
		Name:     userDto.Name,
		LastName: userDto.LastName,
		Email:    userDto.Email,
		Birthday: bDay,
		Password: userDto.Password,
		Status:   true,
		Phone:    userDto.Phone,
	}
	if status, errs := validators.UserValidations(user, userDto); status != http.StatusOK {
		return c.JSON(status, errs)
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	c.Response().Header().Set("Location", fmt.Sprintf("/user/%d", user.ID))
	return c.JSON(http.StatusCreated, user)
}
