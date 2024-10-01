package encryptor

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("tu_secreta_clave") // Cambia esto por una clave segura

type Claims struct {
	Rol string `json:"rol,omitempty"`
	jwt.RegisteredClaims
}

func GenerateJWT(username string, rol string) (string, error) {
	expiratedToken := time.Now().Add(4 * time.Hour)

	claims := &Claims{
		Rol: rol,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiratedToken),
			Subject:   username,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, err
}

func VerifyJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		log.Println(err)
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("firma del token inválida")
		}
		return nil, errors.New("token inválido")
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token no válido")
}
