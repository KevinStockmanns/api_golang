package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/KevinStockmanns/api_golang/internal/constants"
	"github.com/KevinStockmanns/api_golang/internal/db"
	"github.com/KevinStockmanns/api_golang/internal/dtos"
	"github.com/KevinStockmanns/api_golang/internal/encryptor"
	"github.com/KevinStockmanns/api_golang/internal/models"
	"github.com/KevinStockmanns/api_golang/internal/services"
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
		Password: encryptor.EncryptPassword(userDto.Password),
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

	var rol models.Rol
	if err := services.GetOrCreateRol(&rol, string(constants.User)); err != nil {
		return c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Message: "ocurrio un error al crear el usuario",
		})
	}
	user.RolId = rol.ID

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Message: "Ocurrio un error al crear el usuario en la base de datos",
		})
	}

	tx.Commit()

	token, err := encryptor.GenerateJWT(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error al crear token de seguridad")
	}

	userResponse := dtos.UserWithTokenResponseDTO{
		Token: token,
		UserResponseDTO: dtos.UserResponseDTO{
			ID:       user.ID,
			Name:     user.Name,
			LastName: user.LastName,
			Email:    user.Email,
			Birthday: user.Birthday,
			Status:   user.Status,
			Phone:    user.Phone,
			Rol:      rol.Name,
		},
	}

	c.Response().Header().Set("Location", fmt.Sprintf("/user/%d", user.ID))
	return c.JSON(http.StatusCreated, userResponse)
}
