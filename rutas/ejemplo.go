package rutas

import (
	"backend/dto"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Ejemplo_get(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"estado": "ok", "mensaje": "Método GET"})
}
func Ejemplo_get_con_parametro(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Método GET | id=" + id,
	})
}

/*
	func Ejemplo_post(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"estado": "ok", "mensaje":"Método POST"})
	}
*/
func Ejemplo_post(c *gin.Context) {
	var body dto.EjemploDto

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje":  "Método post",
		"correo":   body.Correo,
		"password": body.Password,
	})
}

func Ejemplo_put(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusCreated, gin.H{
		"estado":  "ok",
		"mensaje": "Método PUT | id=" + id,
	})
}

func Ejemplo_delete(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusCreated, gin.H{
		"mensaje": "Método delete | id=" + id,
	})
}

func Ejemplo_query_string(c *gin.Context) {
	id := c.Query("id")
	slug := c.Query("slug")
	c.JSON(200, gin.H{
		"mensaje": "query string | id=" + id + " | slug=" + slug,
	})
}

func Ejemplo_upload(c *gin.Context) {
	file, err := c.FormFile("foto")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"estado":  "error",
			"mensaje": "Ocurrió un error inesperado",
		})
		return
	}
	//nombramos el archivo
	var extension = strings.Split(file.Filename, ".")[1]
	time := strings.Split(time.Now().String(), " ")
	foto := string(time[4][6:14]) + "." + extension
	var archivo string = "public/uploads/fotos/" + foto
	// Upload the file to specific dst.
	c.SaveUploadedFile(file, archivo)

	c.JSON(http.StatusCreated, gin.H{
		"estado":  "ok",
		"message": "Se creó el registro exitosamente",
		"foto":    foto,
	})
}
