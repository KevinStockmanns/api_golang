package encryptor

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/KevinStockmanns/api_golang/internal/constants"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("tu_secreta_clave") // Cambia esto por una clave segura

type Claims struct {
	UserID uint   `json:"userId"`
	Rol    string `json:"rol,omitempty"`
	jwt.RegisteredClaims
}

func GenerateJWT(id uint, username string, rol string) (string, error) {
	expiratedToken := time.Now().Add(4 * time.Hour)

	claims := &Claims{
		UserID: id,
		Rol:    rol,
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
		if strings.Contains(err.Error(), "token is expired") {
			claims := token.Claims.(*Claims)
			expirationTime := claims.ExpiresAt.Time
			timeExpired := time.Since(expirationTime)

			return nil, errors.New("el token ha expirado hace " + timeExpired.String())
		}
		return nil, errors.New("token inválido")
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token no válido")
}

func IsAdmin(claims Claims) bool {
	return claims.Rol == string(constants.Admin) || claims.Rol == string(constants.SuperAdmin)
}
