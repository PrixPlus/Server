package auth

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/prixplus/server/models"
	"github.com/prixplus/server/settings"
)

// Create a new token for a given User
func NewToken(user models.User) (*models.Token, error) {

	sets, err := settings.Get()
	if err != nil {
		return nil, errors.Wrap(err, "getting setting")
	}

	// Create the JWT token
	jwtToken := jwt.New(jwt.GetSigningMethod(sets.JWT.Algorithm))

	timeout := time.Hour * sets.JWT.Expiration

	expire := time.Now().Add(timeout)
	jwtToken.Claims["uid"] = user.Id
	jwtToken.Claims["exp"] = expire.Unix()

	// I could use some key id to identify what secret key are we using
	// but it is optional and isn't utilized in this package
	// jwtToken.Header["kid"] = "Id of the Secret Key used to encrypt this token"

	raw, err := jwtToken.SignedString([]byte(sets.JWT.SecretKey))
	if err != nil {
		return nil, errors.Wrap(err, "creating jwt token")
	}

	token := &models.Token{
		Raw:       raw, // == jwtToken.Raw  ??
		Expire:    expire,
		Method:    jwtToken.Method.Alg(),
		Header:    jwtToken.Header,
		Claims:    jwtToken.Claims,
		Signature: jwtToken.Signature,
		Valid:     jwtToken.Valid,
	}

	return token, nil
}

// Create a new token for a given User
func ParseToken(raw string) (*models.Token, error) {

	sets, err := settings.Get()
	if err != nil {
		return nil, errors.Wrap(err, "getting settings")
	}

	jwtToken, err := jwt.Parse(raw, func(token *jwt.Token) (interface{}, error) {
		// Check if encryption algorithm in token is the same
		if jwt.GetSigningMethod(sets.JWT.Algorithm) != token.Method {
			return nil, errors.New("invalid signing algorithm")
		}

		// I could use some key id to identify whay secret key are we using
		// but it is optional and isn't utilized in this package
		// secretKey, err := myLookupForSecretKey(token.Header["kid"])

		return []byte(sets.JWT.SecretKey), nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "error parsing jwt token")
	}

	if !jwtToken.Valid {
		return nil, errors.New("ivalid token provided or it's expired")
	}

	token := &models.Token{
		Raw:       jwtToken.Raw, // Dunno if it works well  ??
		Method:    jwtToken.Method.Alg(),
		Header:    jwtToken.Header,
		Claims:    jwtToken.Claims,
		Signature: jwtToken.Signature,
		Valid:     jwtToken.Valid,
		//Expire:    expire,
	}

	return token, nil
}

// Get user logged from context
func GetUserIdFromContext(c *gin.Context) (int64, error) {
	id, ok := c.Get("uid")
	if !ok {
		return 0, errors.New("user not logged")
	}

	userId, ok := id.(int64)
	if !ok {
		return 0, errors.New("casting claims")
	}

	return userId, nil
}

// Get user logged from context
func GetUserFromContext(c *gin.Context) (*models.User, error) {
	userId, err := GetUserIdFromContext(c)
	if err != nil {
		return nil, errors.Wrap(err, "getting users from context")
	}

	user := &models.User{Id: userId}
	err = user.Get(nil)
	if err != nil {
		return nil, errors.Wrapf(err, "getting users with your with id %d", userId)
	}

	return user, nil
}
