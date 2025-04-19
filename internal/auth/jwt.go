package auth

import "github.com/golang-jwt/jwt/v5"

type JWTAuthenticator struct {
	secret, aud, iss string
}

func NewJWTAuthenticator(secret, aud, iss string) *JWTAuthenticator {
	return &JWTAuthenticator{
		secret, aud, iss,
	}
}

func (j *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secret))

	if err != nil {
		return "", err
	}

	return tokenString, err
}

func (j *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return nil, nil
}
