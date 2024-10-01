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
	"gorm.io/gorm"
)

func UserSignUp(c echo.Context) error {
	var userDto dtos.UserSignUpDTO

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
	user.Rol = rol

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

func UserLogin(c echo.Context) error {
	var userDto dtos.UserLoginDTO
	if err := c.Bind(&userDto); err != nil {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: "ocurrio un error al leer el cuerpo de la petición"})
	}

	if errs, ok := validators.ValidateDTOs(userDto); !ok {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Message: "error de validación",
			Errors:  errs.Errors,
		})
	}

	var user models.User
	if err := db.DB.Model(user).Preload("Rol").Where("email = ?", userDto.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, dtos.ErrorResponse{
				Message: fmt.Sprintf("no se encontro un usuario con el correo %s en la base de datos", userDto.Email),
			})
		} else {
			return c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
				Message: "ocurrio un error al buscar el usuario",
			})
		}
	}
	if err := encryptor.VerifyPassword(userDto.Password, user.Password); err != nil {
		return c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
			Message: "credenciales inválidas",
		})
	}

	token, _ := encryptor.GenerateJWT(user.Email)
	return c.JSON(http.StatusOK, dtos.UserWithTokenResponseDTO{
		Token: token,
		UserResponseDTO: dtos.UserResponseDTO{
			ID:       user.ID,
			Name:     user.Name,
			LastName: user.LastName,
			Email:    user.Email,
			Birthday: user.Birthday,
			Status:   user.Status,
			Phone:    user.Phone,
			Rol:      user.Rol.Name,
		},
	})
}
