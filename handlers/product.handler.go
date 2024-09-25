package handlers

import (
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/KevinStockmanns/api_golang/db"
	"github.com/KevinStockmanns/api_golang/models"
	"github.com/KevinStockmanns/api_golang/models/wrapper"
	"github.com/KevinStockmanns/api_golang/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetProduct(c echo.Context) error {
	var product models.Product
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "El id debe ser un número")
	}
	if err := db.DB.Model(&models.Product{}).First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, "Producto no encontrado")
		} else {
			return c.JSON(http.StatusInternalServerError, "Error al buscar el producto")
		}
	}
	return c.JSON(200, &product)
}

func PostProduct(c echo.Context) error {
	var product models.Product
	var productPost models.PostProduct

	if err := c.Bind(&productPost); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := utils.Validate.Struct(productPost); err != nil {
		errors := utils.ValidateErrors(err.(validator.ValidationErrors))
		if len(errors) > 0 {
			return c.JSON(http.StatusBadRequest, map[string][]models.ErrorWrapper{
				"errors": errors,
			})
		}
	}

	product.Name = productPost.Name
	product.Status = productPost.Status
	for _, v := range productPost.Versions {
		product.Versions = append(product.Versions, models.Version{
			Name:        v.Name,
			Price:       v.Price,
			ResalePrice: v.ResalePrice,
			Status:      v.Status,
			Date:        time.Now().UTC(),
			Stock:       v.Stock,
			Vistas:      0,
		})
	}

	result := db.DB.Create(&product)
	if result.Error != nil {
		log.Println(result.Error)
		if strings.Contains(result.Error.Error(), "Error 1062") {
			return c.JSON(http.StatusBadRequest, "El nombre del produccto ya esta en uso")
		} else {
			return c.JSON(http.StatusInternalServerError, "error al agregar")
		}
	}

	return c.JSON(http.StatusOK, product)
}

func GetProducts(c echo.Context) error {
	var products []models.Product

	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")

	page := 1
	limit := 10

	if pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil {
			page = p
		}
	}
	if limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil {
			limit = l
		}
	}

	offset := (page - 1) * limit

	result := db.DB.Limit(limit).Offset(offset).Find(&products)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error.Error())
	}

	var total int64
	db.DB.Model(&models.Product{}).Count(&total)

	return c.JSON(http.StatusOK, wrapper.PageResponse{
		Page:          page,
		Size:          limit,
		TotalPage:     int(math.Ceil(float64(total) / float64(limit))),
		TotalElements: total,
		Content:       products,
	})
}
