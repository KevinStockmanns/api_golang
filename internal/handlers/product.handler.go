package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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

func ProductCreate(c echo.Context) error {
	var productDto dtos.ProductCreateDTO
	if err := c.Bind(&productDto); err != nil {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: "ocurrio un error al leer el cuerpo de la petición"})
	}

	if errs, ok := validators.ValidateDTOs(productDto); !ok {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Message: "error de validación",
			Errors:  errs.Errors,
		})
	}
	var product models.Product
	if status, errs := validators.ProductCreate(product, productDto); status != http.StatusOK {
		return c.JSON(http.StatusOK, dtos.ErrorResponse{
			Message: "error de validación",
			Errors:  errs.Errors,
		})
	}
	product.Create(productDto)

	tx := db.DB.Begin()

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{Message: "ocurrio un error al guardar el producto"})
	}

	tx.Commit()

	var productResponse dtos.ProductResponseDTO
	productResponse.Init(product)

	c.Response().Header().Set("Location", fmt.Sprintf("/product/%d", product.ID))
	return c.JSON(http.StatusOK, productResponse)
}

func ProductGet(c echo.Context) error {
	idParam := c.Param("id")
	if !utils.IsInteger(idParam) {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: "el id debe ser un número entero"})
	}

	var product models.Product
	if err := db.DB.Preload("Versions").First(&product, idParam).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, dtos.ErrorResponse{Message: "no se encontro el producto"})
		} else {
			return c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{Message: "ocurrio un error al obtener el producto"})
		}
	}

	var productResponse dtos.ProductResponseDTO
	productResponse.Init(product)

	return c.JSON(http.StatusOK, productResponse)
}

func ProductList(c echo.Context) error {
	page := 1
	size := 10
	limit := 10
	sort := "name"
	order := "ASC"
	status := true

	if statusParam := strings.ToLower(c.QueryParam("status")); statusParam == "false" {
		claims, ok := c.Get("tokenClaims").(*encryptor.Claims)
		if !ok || claims == nil {
			return c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{Message: "autorización requerida para ver productos desactivados"})
		} else {
			if !encryptor.IsAdmin(*claims) {
				return c.JSON(http.StatusForbidden, dtos.ErrorResponse{Message: "no tienes los permisos necesarios para ver los productos desactivados"})
			}
		}

		status = false
	}
	if strings.ToUpper(c.QueryParam("order")) == "ASC" || strings.ToUpper(c.QueryParam("order")) == "DESC" {
		order = strings.ToUpper(c.QueryParam("order"))
	}
	if sortParam := strings.ToLower(c.QueryParam("sort")); sortParam != "" {
		if sortParam == "id" || sortParam == "views" || sortParam == "stock" || sortParam == "date" {
			sort = sortParam
		}
	}
	if c.QueryParam("size") != "" {
		if s, err := strconv.Atoi(c.QueryParam("size")); err == nil && s > 0 {
			size = s
		}
	}
	if c.QueryParam("page") != "" {
		if p, err := strconv.Atoi(c.QueryParam("page")); err == nil && p > 0 {
			page = p
		}
	}
	if size > limit {
		size = limit
	}

	pagination := services.NewPagination[models.Product](page, size, limit)
	pagination.RunQuery(db.DB, "status = ?", []interface{}{status}, fmt.Sprintf("%s %s", sort, order), []string{"Versions"})

	response := make([]dtos.ProductResponseDTO, len(pagination.Content))
	for i, p := range pagination.Content {
		var pResponse dtos.ProductResponseDTO
		pResponse.Init(p)
		response[i] = pResponse
	}

	return c.JSON(http.StatusOK, services.Pagination[dtos.ProductResponseDTO]{
		Page:       pagination.Page,
		Size:       pagination.Size,
		Total:      pagination.Total,
		TotalPages: pagination.TotalPages,
		Content:    response,
	})
}

func ProductUpViews(c echo.Context) error {
	var idVersions dtos.VersionUpViewDTO
	if err := c.Bind(&idVersions); err != nil {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: "ocurrio un errro al leer el cuerpo de la petición"})
	}
	if errs, ok := validators.ValidateDTOs(idVersions); !ok {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Message: "error de validación",
			Errors:  errs.Errors,
		})
	}
	var versions []models.Version
	tx := db.DB.Begin()
	if err := tx.Where("id IN ?", idVersions.IDVersions).Find(&versions); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Message: "ocurrio un error al obtener las versiones",
		})
	}
	for i := range versions {
		versions[i].Views++
	}

	if err := tx.Save(&versions).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Message: "Ocurrió un error al guardar las versiones",
		})
	}

	tx.Commit()

	return c.NoContent(http.StatusNoContent)
}

func ProductUpdate(c echo.Context) error {
	idParam := c.Param("id")
	if !utils.IsInteger(idParam) {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: "el id ingresado debe ser un número entero"})
	}
	var productDto dtos.ProductUpdateDTO
	if c.Bind(&productDto) != nil {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: "ocurrio un error al leer el cuerpo de la petición"})
	}
	if errs, ok := validators.ValidateDTOs(productDto); !ok {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Message: "error de validación",
			Errors:  errs.Errors,
		})
	}

	var product models.Product
	tx := db.DB.Begin()

	if err := tx.Preload("Versions").First(&product, idParam).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, dtos.ErrorResponse{Message: "no se encontro el producto en la base de datos"})
		} else {
			return c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{Message: "ocurrio un error al buscar el producto"})
		}
	}

	if status, errs := validators.ProductUpdate(product, productDto); status != http.StatusOK {
		return c.JSON(status, dtos.ErrorResponse{
			Message: "error de validación",
			Errors:  errs.Errors,
		})
	}

	services.ProductUpdate(&product, productDto, tx)

	tx.Commit()
	var productResponse dtos.ProductResponseDTO
	productResponse.Init(product)

	return c.JSON(http.StatusOK, productResponse)
}

func ProductDelete(c echo.Context) error {
	idParam := c.Param("id")
	if !utils.IsInteger(idParam) {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: "el id ingresado debe ser un número entero"})
	}

	var product models.Product
	if err := db.DB.First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: "no se encontro el producto en la base de datos"})
		} else {
			return c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{Message: "ocurrio un error al buscar el producto"})
		}
	}

	product.Status = false
	db.DB.Save(&product)
	return c.NoContent(http.StatusNoContent)
}

func ProductPriceHistoyList(c echo.Context) error {
	idParam := c.Param("id")
	if !utils.IsInteger(idParam) {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: "el id ingresado debe ser un número entero"})
	}

	var dto dtos.ProductPriceHistoryDTO
	if c.Bind(&dto) != nil {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: "ocurrio un error al leer el cuerpo de la petición"})
	}
	if errs, ok := validators.ValidateDTOs(dto); !ok {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: "error de validación", Errors: errs.Errors})
	}

	var initTime time.Time
	var endTime time.Time
	if time, err := time.Parse("2006-01-02", dto.InitTime); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: "ocurrio un error al leer la fecha ingresada"})
	} else {
		initTime = time
	}
	if time, err := time.Parse("2006-01-02", dto.EndTime); err != nil {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: "ocurrio un error al leer la fecha ingresada"})
	} else {
		endTime = time
	}

	var history []models.PriceHistory
	if err := services.GetHistory(db.DB, &history, initTime, endTime, idParam); err != nil {
		return c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{Message: "ocurrio un error al buscar el historial de precios"})
	}

	return c.JSON(http.StatusOK, history)
}

func ProductChangePrice(c echo.Context) error {
	var dto dtos.ProductChangePrice
	if c.Bind(&dto) != nil {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: "ocurrio un error al leer el cuerpo de la petición"})
	}
	if errs, ok := validators.ValidateDTOs(dto); !ok {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: "error de validación", Errors: errs.Errors})
	}
	if status, errs := validators.ProductChangePrices(dto); status != http.StatusOK {
		return c.JSON(status, dtos.ErrorResponse{Message: "error de validación", Errors: errs.Errors})
	}

	if err := services.ChangeProductsPrice(dto, db.DB); err != nil {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
