package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/prixplus/server/auth"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/prixplus/server/models"
	"golang.org/x/crypto/bcrypt"
)

// Get the User of the current session
func GetMe() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := auth.GetUserFromContext(c)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, errors.Wrapf(err, "user not logged %d"))
			return
		}

		// Not sending password, it will be omitted
		user.Password = ""

		c.JSON(http.StatusOK, gin.H{
			"results": []*models.User{user},
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
			c.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "parsing login JSON"))
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
			c.AbortWithError(http.StatusBadRequest, errors.Wrapf(err, "getting users with email %s", login.Email))
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
			c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "encrypting password"))
			return
		}

		user := &models.User{Email: login.Email, Password: string(hashedPassword)}
		err = user.Insert(nil) // Not using any transaction
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "inserting your user"))
			return
		}

		c.Header("Location", fmt.Sprintf("/api/users/%d", user.Id))

		// Not sending password, it will be omitted
		user.Password = ""

		c.JSON(http.StatusCreated, gin.H{
			"results": []*models.User{user},
		})
	}
}

// Update the user from the session (update me)
func PutUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		me, err := auth.GetUserFromContext(c)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, errors.Wrapf(err, "user not logged %d"))
			return
		}

		id, err := strconv.ParseInt(c.Param("id"), 10, 32)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "converting id from parameters"))
			return
		}

		// Check if user isn't trying to modify another user
		if id != me.Id {
			c.AbortWithError(http.StatusUnauthorized, errors.New("parameter id is different from yours"))
			return
		}

		// Get info received in json
		var user models.User

		err = c.BindJSON(&user)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "parsing user JSON"))
			return
		}

		// Check if user isn't trying to change his id or
		if user.Id != 0 && user.Id != me.Id {
			c.AbortWithError(http.StatusBadRequest, errors.New("can't change your Id"))
			return
		}

		// If user is trying to change the his email
		// so validate the new data
		if len(user.Email) != 0 {
			// Validate email
			if !govalidator.IsEmail(user.Email) {
				c.AbortWithError(http.StatusBadRequest, errors.Errorf("email not valid %s", user.Email))
				return
			}
			me.Email = user.Email
		}

		// If user is trying to change his password
		// so we have to encrypt this new password
		if len(user.Password) != 0 {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "encrypting password"))
			}
			me.Password = string(hashedPassword)
		}

		err = me.Update(nil) // Not using any transaction
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "trying to update your user"))
			return
		}

		// Not sending password, it will be omitted
		me.Password = ""

		c.JSON(http.StatusOK, gin.H{
			"results": []*models.User{me},
		})
	}
}
