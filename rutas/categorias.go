package rutas

import (
	"backend/database"
	"backend/dto"
	"backend/modelos"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug" //go get github.com/gosimple/slug
)

func Categoria_get(c *gin.Context) {
	datos := modelos.Categorias{}
	database.Database.Order("id desc").Find(&datos)
	c.JSON(http.StatusOK, datos)
}

func Categoria_get_con_parametro(c *gin.Context) {
	id := c.Param("id")

	datos := modelos.Categoria{}

	if err := database.Database.First(&datos, id); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"estado":        "error",
			"mensaje":       "Recurso no disponible",
			"errorOpcional": err.Error.Error(), //opcional para ver el error real
		})
	} else {
		c.JSON(http.StatusOK, datos)
	}

}

func Categoria_post(c *gin.Context) {
	//aplicamos middleware de protección de rutas

	var body dto.CategoriaDto
	//validamos que venga el json
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"estado":        "error",
			"mensaje":       "Recurso no disponible",
			"errorOpcional": err.Error(),
		})
		return
	}
	//validamos que no existe el nombre de categoría
	existe := modelos.Categorias{}
	database.Database.Where(&modelos.Categoria{Nombre: body.Nombre}).Find(&existe)
	if len(existe) >= 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"estado":        "error",
			"mensaje":       "Recurso no disponible",
			"errorOpcional": "Ya existe el nombre: " + body.Nombre,
		})
		return
	}
	//creamos el registro
	datos := modelos.Categoria{Nombre: body.Nombre, Slug: slug.Make(body.Nombre)}

	database.Database.Save(&datos)
	//retornamos
	c.JSON(http.StatusCreated, gin.H{
		"estado":  "ok",
		"mensaje": "Se creó el registro exitosamente",
	})

}
func Categoria_put(c *gin.Context) {

	//obtenemos los datos del json request
	var body dto.CategoriaDto
	//validamos que venga el json
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"estado":        "error",
			"mensaje":       "Recurso no disponible",
			"errorOpcional": err.Error(),
		})
		return
	}
	//obtenemos el id path
	id := c.Param("id")
	//validamos que exista la categoría por id
	datos := modelos.Categoria{}
	if err := database.Database.First(&datos, id); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"estado":        "error",
			"mensaje":       "Recurso no disponible",
			"errorOpcional": err.Error.Error(), //opcional para ver el error real
		})
		return
	}
	//editamos el registro
	datos.Nombre = body.Nombre
	datos.Slug = slug.Make(body.Nombre)
	database.Database.Save(&datos)
	//retornamos
	c.JSON(http.StatusOK, gin.H{
		"estado":  "ok",
		"mensaje": "Se modificó el registro exitosamente",
	})
}

func Categoria_delete(c *gin.Context) {
	//obtenemos el id path
	id := c.Param("id")
	//validamos que exista la categoría por id
	datos := modelos.Categoria{}
	if err := database.Database.First(&datos, id); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"estado":        "error",
			"mensaje":       "Recurso no disponible",
			"errorOpcional": err.Error.Error(), //opcional para ver el error real
		})
		return
	}
	//validamos que exista la receta por categoria_id
	categoria_id, _ := strconv.ParseUint(id, 10, 64)
	existe := modelos.Recetas{}
	database.Database.Where(&modelos.Receta{CategoriaID: uint(categoria_id)}).Find(&existe)
	if len(existe) >= 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"estado":        "error",
			"mensaje":       "Ocurrió un error inesperado",
			"errorOpcional": "No se puede borrar la categoría porque existe en una receta",
		})
		return
	}
	//borramos el registro
	database.Database.Delete(&datos)
	//retornamos
	c.JSON(http.StatusOK, gin.H{
		"estado":  "ok",
		"mensaje": "Se eliminó el registro exitosamente",
	})
}
