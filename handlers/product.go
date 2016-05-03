package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"

	"github.com/prixplus/server/errs"

	"github.com/gin-gonic/gin"
	"github.com/prixplus/server/models"
)

// Get a list of products like the given product info
func GetProductList() gin.HandlerFunc {
	return func(c *gin.Context) {

		q := c.Query("q")
		fmt.Println("QUERYYY:", q)

		products, err := models.QueryProducts(q, nil) // Not using any transaction
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "getting the products"))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"results": products,
		})
	}
}

// Get one Product with the given Id
func GetProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

		id, err := strconv.ParseInt(c.Param("id"), 10, 32)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "id isn't a valid integer"))
			return
		}

		product := &models.Product{Id: id}
		products, err := product.GetAll(nil) // Not using any transaction
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "getting the products"))
			return
		}

		if len(products) == 0 {
			c.AbortWithError(errs.Status[errs.ElemNotFound], errors.Wrapf(errs.ElemNotFound, "product not found with id %s", id))
			return
		}

		if len(products) > 1 {
			c.AbortWithError(errs.Status[errs.ElemNotUnique], errors.Wrapf(errs.ElemNotUnique, "more than one product with id %s", id))
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
			c.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "parsing product JSON"))
			return
		}

		// Check if there is no other product with this Gtin
		p := models.Product{Gtin: product.Gtin}
		products, err := p.GetAll(nil) // Not using any transaction
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "getting the products"))
			return
		}

		if len(products) > 0 {
			c.AbortWithError(http.StatusConflict, errors.Wrap(err, "product with GTIN "+product.Gtin+" already exist"))
			return
		}

		err = product.Insert(nil) // Not using any transaction
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "inserting the product"))
			return
		}

		c.Header("Location", fmt.Sprintf("/api/products/%d", product.Id))

		c.JSON(http.StatusCreated, gin.H{
			"results": []models.Product{product},
		})
	}
}

// Update one Product
// This route just work for Gtin and Description
func PutProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

		id, err := strconv.ParseInt(c.Param("id"), 10, 32)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "converting id from parameters"))
			return
		}

		var product models.Product

		err = c.BindJSON(&product)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "parsing product JSON"))
			return
		}

		// Check if user isn't trying to modify another product
		if product.Id != 0 && product.Id != id {
			c.AbortWithError(http.StatusBadRequest, errors.Wrapf(err, "you can't modify product id %d in the route for product id %d", product.Id, id))
			return
		}

		productSaved := models.Product{Id: id}
		err = productSaved.Get(nil) // Not using any transaction
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "getting the product "+productSaved.String()))
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
			c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "updating product"))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"results": []models.Product{productSaved},
		})
	}
}
