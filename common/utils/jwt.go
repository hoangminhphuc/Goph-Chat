package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hoangminhphuc/goph-chat/common"
)

type TokenPayLoad interface {
	UserId() int
	Role() string
}

type Payload struct {
	UID   int    `json:"user_id"`
	URole string `json:"role"`
}

func (t Payload) UserId() int {
	return t.UID
}

func (t Payload) Role() string {
	return t.URole
}

// embeds jwt.StandardClaims, which already satisfies the jwt.Claims interface
type customClaims struct {
	Payload Payload `json:"payload"`  
	jwt.RegisteredClaims
}

func GenerateToken(payload Payload, expiry int, secret string) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		customClaims{
			Payload{
				UID:   payload.UID,
				URole: payload.URole,
			},
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(now.Local().Add(time.Second * time.Duration(expiry))),
				IssuedAt: jwt.NewNumericDate(now.Local()),
				ID: fmt.Sprintf("%d", now.Local().UnixNano()),
			},
		})

	//return final token with 3 parts (signed with hash algo + secret key)
		tokenString, err := token.SignedString([]byte(secret))

		if err != nil {
			return "", common.NewError("Failed to generate token", http.StatusInternalServerError)
		}

		return tokenString, nil
}

func VerifyToken(tokenString string, secret string) (TokenPayLoad, error) {
	t, err := jwt.ParseWithClaims(tokenString, 
		&customClaims{}, 
		// Keyfunc
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
	})

	if err != nil || !t.Valid {
		return nil, common.NewError("Invalid token", http.StatusUnauthorized)
	}



	claim, ok :=t.Claims.(*customClaims)

	if !ok {
		return nil, common.NewError("Invalid token", http.StatusUnauthorized)
	}

	return claim.Payload, nil

}