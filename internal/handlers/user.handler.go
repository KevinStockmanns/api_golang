package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/KevinStockmanns/api_golang/internal/constants"
	"github.com/KevinStockmanns/api_golang/internal/db"
	"github.com/KevinStockmanns/api_golang/internal/dtos"
	"github.com/KevinStockmanns/api_golang/internal/encryptor"
	"github.com/KevinStockmanns/api_golang/internal/models"
	"github.com/KevinStockmanns/api_golang/internal/services"
	"github.com/KevinStockmanns/api_golang/internal/utils"
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
	if status, errs := validators.UserSignUp(user, userDto); status != http.StatusOK {
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
		// log.Println(err)
		return c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Message: "Ocurrio un error al crear el usuario en la base de datos",
		})
	}

	tx.Commit()

	token, err := encryptor.GenerateJWT(user.ID, user.Email, user.Rol.Name)
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

	if stats, errs := validators.UserLogin(user, userDto); stats != http.StatusOK {
		return c.JSON(stats, dtos.ErrorResponse{
			Message: "errores de validación",
			Errors:  errs.Errors,
		})
	}

	token, _ := encryptor.GenerateJWT(user.ID, user.Email, user.Rol.Name)

	var userResponse dtos.UserResponseDTO
	userResponse.Init(user)
	return c.JSON(http.StatusOK, dtos.UserWithTokenResponseDTO{
		Token:           token,
		UserResponseDTO: userResponse,
	})
}

func UserUpdate(c echo.Context) error {
	var userDto dtos.UserUpdateDTO
	id := c.Param("id")
	if !utils.IsInteger(id) {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Message: "el id debe ser un número entero",
		})
	}
	if err := c.Bind(&userDto); err != nil {
		c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Message: "ocurrio un probelma al leer el cuerpo de la petición",
		})
	}
	if errs, ok := validators.ValidateDTOs(userDto); !ok {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Message: "error de validación",
			Errors:  errs.Errors,
		})
	}
	var user models.User
	if err := db.DB.Model(user).Preload("Rol").First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, dtos.ErrorResponse{
				Message: "el usuario no se encuentra en la base de datos",
			})
		} else {
			return c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
				Message: "ocurrio un error al obtener el usuairo",
			})
		}
	}
	idNum, _ := strconv.ParseUint(id, 10, 32)

	if status, errs := validators.UserUpdate(user, userDto, uint(idNum)); status != http.StatusOK {
		return c.JSON(status, dtos.ErrorResponse{
			Message: "error de validación",
			Errors:  errs.Errors,
		})
	}

	tx := db.DB.Begin()
	user.Update(userDto)
	if userDto.Rol != nil {
		var rol models.Rol
		if err := services.GetOrCreateRol(&rol, strings.ToUpper(*userDto.Rol)); err != nil {
			return c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{Message: "ocurrio un error al actualizar el rol"})
		}
		user.Rol = rol
	}

	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Message: "ocurrio un error al guardar el usuario",
		})
	}

	tx.Commit()

	var userResponse dtos.UserResponseDTO
	userResponse.Init(user)

	return c.JSON(http.StatusOK, userResponse)
}

func GetUser(c echo.Context) error {
	id := c.Param("id")
	if !utils.IsInteger(id) {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: "el id debe ser un número entero"})
	}

	tokenClaims := c.Get("tokenClaims").(*encryptor.Claims)
	idInt, _ := strconv.ParseUint(id, 10, 32)
	if tokenClaims.Rol == "USER" && tokenClaims.UserID != uint(idInt) {
		return c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{Message: "no tienes los permisios requeridos para obtener la información de otro usuario"})
	}

	var user models.User
	if err := db.DB.Preload("Rol").First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, dtos.ErrorResponse{Message: "no se encontro ninugn usuario con el id ingresado"})
		} else {
			return c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{Message: "ocurrio un error al obtener el usuario"})
		}
	}

	var userResponse dtos.UserResponseDTO
	userResponse.Init(user)
	return c.JSON(http.StatusOK, userResponse)
}

func UserChangePassword(c echo.Context) error {
	var userDto dtos.UserChangePassword
	if err := c.Bind(&userDto); err != nil {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: "ocurrio un error al leer el cuerpo de la petición"})
	}
	if errs, ok := validators.ValidateDTOs(userDto); !ok {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Message: "error de validación",
			Errors:  errs.Errors,
		})
	}

	if userDto.ActualPassword == userDto.NewPassword {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: "la clave nueva debe ser diferente a la actual"})
	}
	idUser := c.Get("tokenClaims").(*encryptor.Claims)

	var user models.User
	if err := db.DB.First(&user, idUser.UserID).Error; err != nil {
		return c.JSON(http.StatusNotFound, dtos.ErrorResponse{Message: "usuario no encontrado"})
	}
	if err := encryptor.VerifyPassword(userDto.ActualPassword, user.Password); err != nil {
		return c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{Message: "credenciales inválidas"})
	}
	newPass := encryptor.EncryptPassword(userDto.NewPassword)
	user.Password = newPass
	if err := db.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{Message: "ocurrio un error al guardar la clave"})
	}
	return c.NoContent(http.StatusNoContent)
}

func UserList(c echo.Context) error {
	sizeParam := c.QueryParam("size")
	pageParam := c.QueryParam("page")

	size := 10
	page := 1

	if sizeParam != "" {
		if s, err := strconv.Atoi(sizeParam); err == nil && s > 0 {
			size = s
		}
	}
	if pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}
	// var users []models.User

	pagination := services.NewPagination[models.User](page, size, 10)
	pagination.RunQuery(db.DB, "status = ?", []interface{}{true}, "name ASC", []string{"Rol"})

	userDTOs := make([]dtos.UserResponseDTO, len(pagination.Content))
	for i, user := range pagination.Content {
		var userDto dtos.UserResponseDTO
		userDto.Init(user)
		userDTOs[i] = userDto
	}

	return c.JSON(http.StatusOK, services.Pagination[dtos.UserResponseDTO]{
		Page:       pagination.Page,
		Size:       pagination.Size,
		Total:      pagination.Total,
		TotalPages: pagination.TotalPages,
		Content:    userDTOs,
	})
}
