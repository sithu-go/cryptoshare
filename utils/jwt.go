package utils

import (
	"cryptoshare/conf"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Create a struct that will be encoded to a JWT.
// We add jwt.RegisteredClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// generate admin access token
func GenerateAccessToken(username string, isAdmin bool) (string, error) {
	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(conf.PrivateKey)

	if err != nil {
		log.Println("Failed to sign id token string")
		return "", err
	}

	return ss, nil
}

func ValidateAccessToken(tokenString string) (*Claims, error) {
	// // Initialize a new instance of `Claims`
	claims := &Claims{}

	// // Parse the JWT string and store the result in `claims`.
	// // Note that we are passing the key in this method as well. This method will return an error
	// // if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// // or if the signature does not match
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return conf.PublicKey, nil
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("access token is invalid")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("access token valid but couldn't parse claims")
	}

	return claims, nil

}

func GenerateRefreshToken(token string) (string, error) {
	claim, err := ValidateAccessToken(token)
	if err != nil {
		return "", err
	}

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Until(claim.ExpiresAt.Time) > 30*time.Second {
		return "", errors.New("token is not expired, yet")
	}

	refreshToken, err := GenerateAccessToken(claim.Username, claim.IsAdmin)
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}
