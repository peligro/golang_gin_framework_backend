package jwt

import (
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func GenerarJWT(correo string, nombre string, id uint) (string, error) {
	errorVariables := godotenv.Load()
	if errorVariables != nil {
		panic("Error loading .env file")
	}
	miClave := []byte(os.Getenv("SECRET_JWT"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"correo":         nombre,
		"nombre":         nombre,
		"generado_desde": "https://www.cesarcancino.com",
		"id":             id,
		"iat":            time.Now().Unix(),
		"exp":            time.Now().Add(time.Hour * 24).Unix(), //24 HORAS
	})
	tokenString, err := token.SignedString(miClave)
	return tokenString, err
}
