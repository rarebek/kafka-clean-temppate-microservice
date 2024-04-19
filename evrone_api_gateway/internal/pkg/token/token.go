package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrValidationErrorMalformed  = errors.New("Token is malformed")
	ErrTokenExpiredOrNotValidYet = errors.New("Token is either expired or not active yet")
)

func GenerateJwtToken(jwtsecret string, claims *jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(jwtsecret))
}

func GenerateToken(sub, token_type, jwtsecret string, access_ttl, refresh_ttl time.Duration, optionalFields ...map[string]interface{}) (string, string, error) {
	accessClaims := jwt.MapClaims{
		"sub":  sub,
		"type": token_type,
		"exp":  time.Now().Add(access_ttl).Unix(),
	}

	for _, fields := range optionalFields {
		for key, value := range fields {
			accessClaims[key] = value
		}
	}

	// generate access token
	access_token, err := GenerateJwtToken(jwtsecret, &accessClaims)
	if err != nil {
		return "", "", err
	}

	// generate refresh token
	refresh_token, err := GenerateJwtToken(jwtsecret, &jwt.MapClaims{
		"exp": time.Now().Add(refresh_ttl).Unix(),
		"sub": sub,
	})
	if err != nil {
		return "", "", err
	}
	return access_token, refresh_token, err
}

func ParseJwtToken(tokenStr, jwtsecret string) (map[string]interface{}, error) {
	var claims map[string]interface{}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtsecret), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return claims, ErrValidationErrorMalformed
			}
			if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return claims, ErrTokenExpiredOrNotValidYet
			}
		}
		return claims, fmt.Errorf("Couldn't handle this token: %w", err)
	}
	// get claims
	if mapClaims, ok := token.Claims.(jwt.MapClaims); ok {
		claims = mapClaims
	}
	return claims, nil
}
