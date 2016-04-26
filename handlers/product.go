package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/prixplus/server/errs"

	"github.com/gin-gonic/gin"
	"github.com/prixplus/server/models"
)

// Get a list of products like the given product info
func GetProductList() gin.HandlerFunc {
	return func(c *gin.Context) {

		var product models.Product

		err := c.BindJSON(&product)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("Error parsing JSON: "+err.Error()))
			return
		}

		products, err := product.GetAll(nil) // Not using any transaction
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("Error getting the products: "+err.Error()))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"results": products,
		})
	}
}

// Get a list of products like the given product info
func GetProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

		id, err := strconv.ParseInt(c.Param("id"), 10, 32)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("Error converting product Id: "+err.Error()))
			return
		}

		product := &models.Product{Id: id}
		products, err := product.GetAll(nil) // Not using any transaction
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("Error getting the products: "+err.Error()))
			return
		}

		if len(products) == 0 {
			c.AbortWithError(errs.Status[errs.ElementNotFound], errs.ElementNotFound)
			return
		}

		if len(products) > 1 {
			c.AbortWithError(http.StatusConflict, errors.New("Sorry but we've found more than one product with this Id.."))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"results": products,
		})
	}
}

// Create a Product
// Returns a Location header with the location of the brand new content
func PostProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

		var product models.Product

		err := c.BindJSON(&product)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("Error parsing JSON: "+err.Error()))
			return
		}

		// Check if there is no other product with this Gtin
		p := models.Product{Gtin: product.Gtin}
		products, err := p.GetAll(nil) // Not using any transaction
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("Error getting the products: "+err.Error()))
			return
		}

		if len(products) > 0 {
			c.AbortWithError(http.StatusConflict, errors.New("Sorry, product with GTIN "+product.Gtin+" already exist"))
			return
		}

		err = product.Insert(nil) // Not using any transaction
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("Error inserting the product: "+err.Error()))
			return
		}

		c.Header("Location", fmt.Sprintf("/api/products/%d", product.Id))

		c.JSON(http.StatusCreated, gin.H{
			"results": []models.Product{product},
		})
	}
}

// Update an User
func PutProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

		var product models.Product

		err := c.BindJSON(&product)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("Error parsing JSON: "+err.Error()))
			return
		}

		id, err := strconv.ParseInt(c.Param("id"), 10, 32)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("Error converting product Id: "+err.Error()))
			return
		}

		// Check if user isn't trying to modify another product
		if product.Id != 0 && product.Id != id {
			c.AbortWithError(http.StatusUnauthorized, errors.New("You can't update other users info"))
			return
		}

		productSaved := models.Product{Id: id}
		err = productSaved.Get(nil) // Not using any transaction
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("Error getting the product like "+productSaved.String()+": "+err.Error()))
			return
		}

		// Make changes in current product
		if len(product.Gtin) > 0 {
			productSaved.Gtin = product.Gtin
		}

		// Make changes in current product
		if len(product.Description) > 0 {
			productSaved.Description = product.Description
		}

		err = productSaved.Update(nil) // Not using any transaction
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("Error updating product: "+err.Error()))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"results": []models.Product{productSaved},
		})
	}
}
