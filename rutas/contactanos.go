package rutas

import (
	"backend/database"
	"backend/dto"
	"backend/modelos"
	"backend/utilidades"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

/*
	{
	    "nombre":"Juan Pérez",
	    "correo":"info@tamila.cl",
	    "telefono":"",
	    "mensaje":"hola quiero contactarlos con ñandú"
	}
*/
func Contactanos_post(c *gin.Context) {

	var body dto.ContactanosDto

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"estado":        "error",
			"mensaje":       "Se produjo un error inesperado",
			"errorOpcional": err.Error(),
		})
		return
	}

	// creamos el registro
	datos := modelos.Contacto{Nombre: body.Nombre, Correo: body.Correo, Telefono: body.Telefono, Mensaje: body.Mensaje, Fecha: time.Now()}
	database.Database.Save(&datos)
	//enviamos el correo
	var mensaje = "<h1>Mensaje recibido</h1><ul> <li>Nombre: " + body.Nombre + "</li><li>E-Mail: " + body.Correo + "</li><li>Teléfono: " + body.Telefono + "</li><li>Mensaje: " + body.Mensaje + "</li></ul>"
	utilidades.EnviarCorreo(body.Correo, "Prueba Golang", mensaje)
	//retornamos
	c.JSON(http.StatusOK, gin.H{
		"estado":  "ok",
		"mensaje": "Se creó el registro exitosamente",
	})
}
