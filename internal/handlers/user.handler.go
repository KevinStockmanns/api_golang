package handlers

import (
	"fmt"
	"log"
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
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Message: "no se puedo leer el cuerpo de la petición",
		})
	}
	if errors, ok := validators.ValidateDTOs(userDto); !ok {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Message: "errores de validación",
			Errors:  errors.Errors,
		})
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
	user.Normalize()
	if status, errs := validators.UserValidations(user, userDto); status != http.StatusOK {
		return c.JSON(status, dtos.ErrorResponse{
			Message: "errores de validación",
			Errors:  errs.Errors,
		})
	}

	tx := db.DB.Begin()

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Message: "Ocurrio un error al crear el usuario en la base de datos",
		})
	}

	tx.Commit()

	c.Response().Header().Set("Location", fmt.Sprintf("/user/%d", user.ID))
	return c.JSON(http.StatusCreated, user)
}
