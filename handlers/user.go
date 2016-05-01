package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/prixplus/server/models"
	"golang.org/x/crypto/bcrypt"
)

// Get the User of the current session
func GetMe() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Identify if user is or isn't logged
		// with a valid auth token
		sessionId, ok := c.Get("id")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"messages": []string{"This resource is just for authenticated users"},
			})
			return
		}
		uId := sessionId.(int64)

		user := models.User{Id: uId}
		err := user.Get(nil) // Not using any transaction
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("Error getting users with your Id: "+err.Error()))
			return
		}

		// Not sending password, it will be omitted
		user.Password = ""

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
			c.JSON(http.StatusBadRequest, gin.H{
				"messages": []string{"Email or Password can not be empty"},
			})
			return
		}

		// Test if already exists an user with this email
		userSaved := models.User{Email: login.Email}
		users, err := userSaved.GetAll(nil) // Not using any transaction
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("Error getting users with email "+login.Email+": "+err.Error()))
			return
		}

		if len(users) > 0 {
			c.JSON(http.StatusConflict, gin.H{
				"messages": []string{"Sorry, email " + login.Email + " already in use"},
			})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(login.Password), bcrypt.DefaultCost)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("Error encrypting password: "+err.Error()))
			return
		}

		user := models.User{Email: login.Email, Password: string(hashedPassword)}
		err = user.Insert(nil) // Not using any transaction
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("Error inserting your user: "+err.Error()))
			return
		}

		c.Header("Location", fmt.Sprintf("/api/users/%d", user.Id))

		// Not sending password, it will be omitted
		user.Password = ""

		c.JSON(http.StatusCreated, gin.H{
			"results": []models.User{user},
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
			c.JSON(http.StatusUnauthorized, gin.H{
				"messages": []string{"This resource is just for authenticated users"},
			})
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

		// Checks if this user really exists
		userSaved := models.User{Id: uId}
		err = userSaved.Get(nil) // Not using any transaction
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

		// If user is trying to change the his email
		// so validate the new data
		if len(user.Email) != 0 {
			// Validate email
			if !govalidator.IsEmail(user.Email) {
				c.AbortWithError(http.StatusBadRequest, errors.New("Error setting this email: "+user.Email))
				return
			}
		} else {
			// Or just use the saved email instead
			user.Email = userSaved.Email
		}

		// If user is trying to change his password
		// so we have to encrypt this new password
		if len(user.Password) != 0 {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, errors.New("Error encrypting password: "+err.Error()))
			}
			user.Password = string(hashedPassword)
		} else {
			// Or just use the saved password instead
			user.Password = userSaved.Password
		}

		err = user.Update(nil) // Not using any transaction
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("Error trying to update your user: "+err.Error()))
			return
		}

		// Not sending password, it will be omitted
		user.Password = ""

		c.JSON(http.StatusOK, gin.H{
			"results": []models.User{user},
		})
	}
}
