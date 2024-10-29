package middleware

import (
	"backend/database"
	"backend/modelos"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func ValidarJWTMiddleware(c *gin.Context) {
	errorVariables := godotenv.Load()
	if errorVariables != nil {

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"estado":         "error",
			"mensaje":        "No autorizado",
			"estadoOpcional": "Error con las variables locales",
		})

	}
	miClave := []byte(os.Getenv("SECRET_JWT"))
	var header = c.GetHeader("Authorization")
	//validamos si el header tiene valor
	if len(header) == 0 {

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"estado":         "error",
			"mensaje":        "No autorizado",
			"estadoOpcional": "Error porque no viene el token",
		})
		return
	}
	// separamos el valor del header
	splitBearer := strings.Split(header, " ")
	if len(splitBearer) != 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"estado":         "error",
			"mensaje":        "No autorizado",
			"estadoOpcional": "Error porque el token trae un solo texto",
		})
		return
	}
	// validamos que el header traiga el valor del Bearer
	splitToken := strings.Split(splitBearer[1], ".")
	if len(splitToken) != 3 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"estado":         "error",
			"mensaje":        "No autorizado",
			"estadoOpcional": "Error porque el header no trae el valor Bearer",
		})
		return
	}
	tk := strings.TrimSpace(splitBearer[1])
	//validamos que el jwt esté correctamente firmado
	token, err := jwt.Parse(tk, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"estado":         "error",
				"mensaje":        "No autorizado",
				"estadoOpcional": "Error con la firma del token",
			})

		}

		return miClave, nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"estado":         "error",
			"mensaje":        "No autorizado",
			"estadoOpcional": "Error con el formato del token",
		})
		return
	}
	//obtenemos los datos del JWT
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		//preguntamos si existe correo
		datos := modelos.Usuario{}

		if err := database.Database.First(&datos, claims["id"]); err.Error != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"estado":         "error",
				"mensaje":        "No autorizado",
				"estadoOpcional": "Error con el id del usuario informado en el token",
			})
			return
		} else {

			c.Next()
		}
	} else {

		return
	}

}
