package rutas

import (
	"backend/database"
	"backend/dto"
	"backend/jwt"
	"backend/modelos"
	"backend/utilidades"
	"backend/validaciones"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Seguridad_registro(c *gin.Context) {
	var body dto.UsuarioDto
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"estado":         "error",
			"mensaje":        "Ocurrió un error inesperado",
			"estadoOpcional": err.Error(),
		})
		return
	}
	//validamos
	if len(body.Nombre) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"estado":  "error",
			"mensaje": "El nombre es obligatorio",
		})
		return
	}
	if len(body.Correo) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"estado":  "error",
			"mensaje": "El E-Mail es obligatorio",
		})
		return
	}

	if validaciones.Regex_correo.FindStringSubmatch(body.Correo) == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"estado":  "error",
			"mensaje": "El E-Mail ingresado no es válido",
		})
		return
	}
	if len(body.Password) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"estado":  "error",
			"mensaje": "El Password es obligatorio",
		})
		return
	}
	//preguntamos si existe correo
	existe := modelos.Usuarios{}
	database.Database.Where(&modelos.Usuario{Correo: body.Correo}).Find(&existe)
	if len(existe) >= 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"estado":        "error",
			"mensaje":       "Recurso no disponible",
			"errorOpcional": "Ya existe el correo: " + body.Correo,
		})
		return
	}

	//generamos hash con bcrypt
	costo := 10
	bytes, _ := bcrypt.GenerateFromPassword([]byte(body.Password), costo)
	//insertamos
	uuid := uuid.New()
	save := modelos.Usuario{
		Nombre:   body.Nombre,
		Correo:   body.Correo,
		Password: string(bytes),
		EstadoID: 2,
		Token:    uuid.String(),
		Fecha:    time.Now()}
	database.Database.Save(&save)

	//obtenemos la URL del proyecto
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	url := scheme + "://" + c.Request.Host + "/api/v1/seguridad/verificacion/" + uuid.String()
	//enviar mail de verificación
	var mensaje = "<h1>Verificación de cuenta</h1>Hola " + body.Nombre + " te haz registrado exitosamente, por favor haz clic en <a href='" + url + "'>" + url + "</a> para verificar tu mail<br/>o copia y pega el siguiente enlace en tu navegador favorito " + url + " "
	utilidades.EnviarCorreo(body.Correo, "Verificación cuenta", mensaje)
	//retornamos
	c.JSON(http.StatusCreated, gin.H{
		"estado":  "ok",
		"mensaje": "Se creó el registro exitosamente",
	})
}
func Seguridad_verificacion(c *gin.Context) {
	token := c.Param("token")

	//preguntamos si existe correo
	datos := modelos.Usuarios{}
	database.Database.Where(&modelos.Usuario{Token: token, EstadoID: 2}).Find(&datos)
	if len(datos) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"estado":        "error",
			"mensaje":       "Recurso no disponible",
			"errorOpcional": "No existe el token " + token + " para un usuario no activo",
		})
		return
	}
	//modificamos el registro
	datos[0].Token = ""
	datos[0].EstadoID = 1
	database.Database.Save(&datos)

	//retornamos
	c.Redirect(http.StatusMovedPermanently, os.Getenv("RUTA_FRONTEND"))
}
func Seguridad_login(c *gin.Context) {
	var body dto.LoginDto
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"estado":         "error",
			"mensaje":        "Ocurrió un error inesperado",
			"estadoOpcional": err.Error(),
		})
		return
	}
	//validamos

	if len(body.Correo) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"estado":  "error",
			"mensaje": "El E-Mail es obligatorio",
		})
		return
	}

	if validaciones.Regex_correo.FindStringSubmatch(body.Correo) == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"estado":  "error",
			"mensaje": "El E-Mail ingresado no es válido",
		})
		return
	}
	if len(body.Password) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"estado":  "error",
			"mensaje": "El Password es obligatorio",
		})
		return
	}
	//preguntamos si existe correo

	usuario := modelos.Usuarios{}
	database.Database.Where(&modelos.Usuario{Correo: body.Correo}).Where(&modelos.Usuario{EstadoID: 1}).Find(&usuario)

	if len(usuario) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"estado":        "error",
			"mensaje":       "Recurso no disponible",
			"errorOpcional": "La credenciales ingresadas con inválidas",
		})
		return
	}
	//usamos bcrypt para comparar password
	passwordBytes := []byte(body.Password)
	passwordBD := []byte(usuario[0].Password)

	errPassword := bcrypt.CompareHashAndPassword(passwordBD, passwordBytes)

	if errPassword != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"estado":        "error",
			"mensaje":       "Recurso no disponible",
			"errorOpcional": "La credenciales ingresadas con inválidas" + errPassword.Error(),
		})
		return
	} else {

		jwtKey, errJWT := jwt.GenerarJWT(usuario[0].Correo, usuario[0].Nombre, usuario[0].Id)
		if errJWT != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"estado":        "error",
				"mensaje":       "Recurso no disponible",
				"errorOpcional": "Ocurrió un error al intentar generar el token: " + errJWT.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"id":     usuario[0].Id,
			"nombre": usuario[0].Nombre,
			"token":  jwtKey,
		})
		return
	}

}
