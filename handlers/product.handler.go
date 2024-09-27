package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/KevinStockmanns/api_golang/db"
	"github.com/KevinStockmanns/api_golang/models"
	"github.com/KevinStockmanns/api_golang/models/wrapper"
	"github.com/KevinStockmanns/api_golang/utils"
	"github.com/KevinStockmanns/api_golang/validators"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetProduct(c echo.Context) error {
	var product models.Product
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "El id debe ser un número")
	}
	if err := db.DB.Model(&models.Product{}).Preload("Versions").First(&product, id).Error; err != nil {
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
	if errors := productPost.NormalizeAndValidate(); len(errors) > 0 {
		return c.JSON(http.StatusBadRequest, map[string][]wrapper.ErrorWrapper{
			"errors": errors,
		})
	}
	if status, err := validators.PostValidations(productPost); status != http.StatusOK {
		errors := wrapper.ErrorsWrapper{}
		errors.AddEror(err)
		return c.JSON(status, errors)
	}

	product.Init(productPost)

	result := db.DB.Create(&product)
	if result.Error != nil {
		log.Println(result.Error)
		if strings.Contains(result.Error.Error(), "Error 1062") {
			return c.JSON(http.StatusBadRequest, "El nombre del produccto ya esta en uso")
		} else {
			return c.JSON(http.StatusInternalServerError, "error al agregar")
		}
	}
	c.Response().Header().Set("Location", fmt.Sprintf("/product/%d", product.ID))

	return c.JSON(http.StatusCreated, product)
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
			if limit > 10 {
				limit = 10
			}
		}
	}

	offset := (page - 1) * limit

	result := db.DB.Limit(limit).Offset(offset).Find(&products)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error.Error())
	}

	var total int64
	db.DB.Model(&models.Product{}).Count(&total)

	return c.JSON(http.StatusOK, wrapper.PageWrapper{
		Page:          page,
		Size:          limit,
		TotalPage:     int(math.Ceil(float64(total) / float64(limit))),
		TotalElements: total,
		Content:       products,
	})
}

func PutProduct(c echo.Context) error {
	var product models.Product
	var productDto models.PutProduct
	id := c.Param("id")
	if isNum := utils.IsInt(id); !isNum {
		return c.JSON(http.StatusBadRequest, wrapper.ErrorWrapper{Error: "el id debe ser un número entero"})
	}

	if err := c.Bind(&productDto); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if errors := productDto.NormalizeAndValidate(); len(errors) > 0 {
		return c.JSON(http.StatusBadRequest, map[string][]wrapper.ErrorWrapper{
			"errors": errors,
		})
	}

	if err := db.DB.Preload("Versions").First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, wrapper.ErrorWrapper{Error: "producto no encontrado"})
		} else {
			return c.JSON(http.StatusInternalServerError, wrapper.ErrorWrapper{Error: "ocuurio un error al obtener el producto"})
		}
	}

	if status, err := validators.PutProductValidation(product, productDto); status != http.StatusOK {
		var errors []wrapper.ErrorWrapper
		errors = append(errors, err)
		responseBody, jsonErr := json.Marshal(map[string][]wrapper.ErrorWrapper{
			"errors": errors,
		})
		if jsonErr != nil {
			return c.JSON(http.StatusInternalServerError, "Error al convertir la respuesta a JSON")
		}
		return c.JSONBlob(status, responseBody)
	}

	tx := db.DB.Begin()

	product.Update(productDto, tx)
	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, "ocurrio un error al guardar los cambios")
	}
	if len(product.Versions) > 6 {
		tx.Rollback()
		errors := wrapper.ErrorsWrapper{}
		errors.AddNewError("versions", "el producto puede tener hasta 6 versiones")
		return c.JSON(http.StatusBadRequest, errors)
	}

	tx.Commit()
	return c.JSON(http.StatusOK, product)
}

func DeleteProduct(c echo.Context) error {
	id := c.Param("id")
	var product models.Product
	if !utils.IsInt(id) {
		errors := wrapper.ErrorsWrapper{}
		errors.AddNewError("", "el id ingresado no es un número entero")
		return c.JSON(http.StatusBadRequest, errors)
	}

	if err := db.DB.First(&product, id).Error; err != nil {
		errors := wrapper.ErrorsWrapper{}
		if err == gorm.ErrRecordNotFound {
			errors.AddNewError("", "el producto no se encontro con el id ingresado")
			return c.JSON(http.StatusNotFound, errors)
		} else {
			errors.AddNewError("", "ocurrio un error al buscar el producto")
			return c.JSON(http.StatusInternalServerError, errors)
		}
	}

	product.Status = false
	if err := db.DB.Save(&product).Error; err != nil {
		errors := wrapper.ErrorsWrapper{}
		errors.AddNewError("", "ocurrio un error al realizar los cambios")
		return c.JSON(http.StatusInternalServerError, errors)
	}

	return c.NoContent(http.StatusNoContent)
}
