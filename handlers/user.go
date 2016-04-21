package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/prixplus/server/models"
)

// Get the User of the current session
func GetMe() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Identify if user is or isn't logged
		// with a valid auth token
		sessionId, ok := c.Get("id")
		if !ok {
			c.AbortWithError(http.StatusUnauthorized, errors.New("This resource is just for authenticated users"))
			return
		}
		uId := sessionId.(int64)

		user := models.User{Id: uId}
		err := user.Get(nil) // Not using any transaction
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("Error getting users with your Id: "+err.Error()))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"results": []models.User{user},
		})
	}
}

// Create an User
// Return an Location header with the location of the new content
// in body returns the location and the new user as its results
func PostUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		var login models.Login

		err := c.BindJSON(&login)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("Error parsing JSON: "+err.Error()))
			return
		}

		if len(login.Email) == 0 || len(login.Password) == 0 {
			c.AbortWithError(http.StatusBadRequest, errors.New("Email or Password can not be empty"))
			return
		}

		// Test if already exists an user with this email
		u := models.User{Email: login.Email}
		users, err := u.GetAll(nil) // Not using any transaction
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("Error getting users with email "+login.Email+": "+err.Error()))
			return
		}

		if len(users) > 0 {
			c.AbortWithError(http.StatusConflict, errors.New("Sorry, email "+login.Email+" already taken"))
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(login.Password), bcrypt.DefaultCost)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("Error encrypting password: "+err.Error()))
		}

		user := models.User{Email: login.Email, Password: string(hashedPassword)}
		err = user.Insert(nil) // Not using any transaction
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("Error inserting your user: "+err.Error()))
			return
		}

		c.Header("Location", fmt.Sprintf("/api/users/%d", user.Id))

		c.JSON(http.StatusCreated, gin.H{
			"location": fmt.Sprintf("/api/users/%d", user.Id),
			"results":  []models.User{user},
		})
	}
}

// Update an User
func PutUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Identify if user is or isn't logged
		// with a valid auth token
		sessionId, ok := c.Get("id")
		if !ok {
			c.AbortWithError(http.StatusUnauthorized, errors.New("This resource is just for authenticated users"))
			return
		}
		uId := sessionId.(int64)

		id, err := strconv.ParseInt(c.Param("id"), 10, 32)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("Error converting user Id: "+err.Error()))
			return
		}

		// Check if user isn't trying to modify another user
		if id != uId {
			c.AbortWithError(http.StatusUnauthorized, errors.New("You can't update other users info"))
			return
		}

		// Get info received in json
		var user models.User

		err = c.BindJSON(&user)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("Error parsing JSON: "+err.Error()))
			return
		}

		u := models.User{Id: uId}
		err = u.Get(nil) // Not using any transaction
		if err != nil {
			c.AbortWithError(http.StatusNoContent, errors.New("Error getting user: "+err.Error()))
			return
		}

		// Check if user isn't trying to
		// Change his Id or
		// to change his email or password to ant empty value
		if user.Id != 0 && user.Id != uId {
			c.AbortWithError(http.StatusBadRequest, errors.New("You can't change your Id"))
			return
		}

		// If user didn't send his Id
		user.Id = uId

		// Check if he isn't trying to change his email
		if len(user.Email) == 0 {
			user.Email = u.Email
		}

		// Check if he isn't trying to change his password
		if len(user.Password) == 0 {
			user.Password = u.Password
		}

		err = user.Update(nil) // Not using any transaction
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("Error trying to update your user: "+err.Error()))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"results": []models.User{user},
		})
	}
}
