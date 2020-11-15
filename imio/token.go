package imio

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"log"
	"net/http"
	"time"
)

type IBClaims struct {
	id string
	jwt.StandardClaims
}

type TokenError struct {
	message string
}

func (err *TokenError) Error() string {
	return err.message
}

func tokenVerify(r *http.Request) *AppError {
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(Secret), nil
		})
	if err == nil {
		if !token.Valid {
			return &AppError{message: "Token is not valid", statusCode: 400}
		}
	} else {
		return &AppError{error: err, message: err.Error(), statusCode: 500}
	}
	return nil
}

func createToken(key string) (string, *AppError) {
	claims := IBClaims{
		key,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(72)).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(Secret)
	if err != nil {
		log.Fatal(err)
		return "", &AppError{error: err, message: "生成token失败", statusCode: 500}
	}
	return signedToken, nil
}
