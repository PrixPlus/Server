package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/prixplus/server/models"
	"github.com/prixplus/server/settings"
)

// Create a new token for a given User
func NewToken(user models.User) (*models.Token, error) {

	sets, err := settings.Get()
	if err != nil {
		return nil, err
	}

	// Create the JWT token
	jwtToken := jwt.New(jwt.GetSigningMethod(sets.JWT.Algorithm))

	timeout := time.Hour * sets.JWT.Expiration

	expire := time.Now().Add(timeout)
	jwtToken.Claims["id"] = user.Id
	jwtToken.Claims["exp"] = expire.Unix()

	// I could use some key id to identify what secret key are we using
	// but it is optional and isn't utilized in this package
	// jwtToken.Header["kid"] = "Id of the Secret Key used to encrypt this token"

	raw, err := jwtToken.SignedString([]byte(sets.JWT.SecretKey))
	if err != nil {
		return nil, err
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
		return nil, err
	}

	jwtToken, err := jwt.Parse(raw, func(token *jwt.Token) (interface{}, error) {
		// Check if encryption algorithm in token is the same
		if jwt.GetSigningMethod(sets.JWT.Algorithm) != token.Method {
			return nil, errors.New("Invalid signing algorithm")
		}

		// I could use some key id to identify whay secret key are we using
		// but it is optional and isn't utilized in this package
		// secretKey, err := myLookupForSecretKey(token.Header["kid"])

		return []byte(sets.JWT.SecretKey), nil
	})

	if err != nil || !jwtToken.Valid {
		return nil, errors.New("Ivalid token provided or it's expired")
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
