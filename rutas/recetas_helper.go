package rutas

import (
	"backend/database"
	"backend/dto"
	"backend/modelos"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Receta_Helper_Por_Usuario_get(c *gin.Context) {
	usuario_id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	//validamos que exista el usuario_id
	user := modelos.Usuario{}
	if err := database.Database.First(&user, usuario_id); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"estado":        "error",
			"mensaje":       "Recurso no disponible",
			"errorOpcional": err.Error.Error(), //opcional para ver el error real
		})
		return
	}
	datos := modelos.Recetas{}
	database.Database.Where(&modelos.Receta{UsuarioID: uint(usuario_id)}).Order("id desc").Preload("Categoria").Preload("Usuario").Find(&datos)
	var arreglo = dto.RecetasResponse{}
	//obtenemos la URL del proyecto
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	for _, dato := range datos {
		//formateamos fecha
		fecha := fmt.Sprintf("%d/%02d/%02d", dato.Fecha.Day(), dato.Fecha.Month(), dato.Fecha.Year())
		//formeatemos la url de la foto
		foto := scheme + "://" + c.Request.Host + "/public/uploads/recetas/" + dato.Foto
		//llenamos el arreglo de tipo struct
		arreglo = append(arreglo, dto.RecetaResponse{Id: dato.Id, Nombre: dato.Nombre, Slug: dato.Slug, CategoriaDtoId: dato.CategoriaID, Categoria: dato.Categoria.Nombre, UsuarioId: dato.UsuarioID, Usuario: dato.Usuario.Nombre, Tiempo: dato.Tiempo, Descripcion: dato.Descripcion, Foto: foto, Fecha: fecha})

	}

	c.JSON(http.StatusOK, arreglo)
}
func Receta_Helper_Home_get(c *gin.Context) {
	datos := modelos.Recetas{}
	database.Database.Order("id desc").Preload("Categoria").Preload("Usuario").Limit(3).Find(&datos)
	var arreglo = dto.RecetasResponse{}
	//obtenemos la URL del proyecto
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	for _, dato := range datos {
		//formateamos fecha
		fecha := fmt.Sprintf("%d/%02d/%02d", dato.Fecha.Day(), dato.Fecha.Month(), dato.Fecha.Year())
		//formeatemos la url de la foto
		foto := scheme + "://" + c.Request.Host + "/public/uploads/recetas/" + dato.Foto
		//llenamos el arreglo de tipo struct
		arreglo = append(arreglo, dto.RecetaResponse{Id: dato.Id, Nombre: dato.Nombre, Slug: dato.Slug, CategoriaDtoId: dato.CategoriaID, Categoria: dato.Categoria.Nombre, UsuarioId: dato.UsuarioID, Usuario: dato.Usuario.Nombre, Tiempo: dato.Tiempo, Descripcion: dato.Descripcion, Foto: foto, Fecha: fecha})

	}

	c.JSON(http.StatusOK, arreglo)
}
func Receta_Helper_Slug_get(c *gin.Context) {
	datos := modelos.Recetas{}
	database.Database.Where(&modelos.Receta{Slug: c.Param("slug")}).Order("id desc").Preload("Categoria").Preload("Usuario").Limit(1).Find(&datos)
	//validamos que si no hay datos por el slug retorne 404
	if len(datos) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"estado":        "error",
			"mensaje":       "Recurso no disponible",
			"errorOpcional": "No existe la receta por el slug " + c.Param("slug"),
		})
		return
	}
	//obtenemos la URL del proyecto
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	//formateamos fecha
	fecha := fmt.Sprintf("%d/%02d/%02d", datos[0].Fecha.Day(), datos[0].Fecha.Month(), datos[0].Fecha.Year())
	//formeatemos la url de la foto
	foto := scheme + "://" + c.Request.Host + "/public/uploads/recetas/" + datos[0].Foto

	c.JSON(http.StatusOK, dto.RecetaResponse{Id: datos[0].Id, Nombre: datos[0].Nombre, Slug: datos[0].Slug, CategoriaDtoId: datos[0].CategoriaID, Categoria: datos[0].Categoria.Nombre, UsuarioId: datos[0].UsuarioID, Usuario: datos[0].Usuario.Nombre, Tiempo: datos[0].Tiempo, Descripcion: datos[0].Descripcion, Foto: foto, Fecha: fecha})
}
func Receta_Helper_Buscador_get(c *gin.Context) {
	datos := modelos.Recetas{}
	database.Database.Where("nombre LIKE ?", "%"+c.Query("search")+"%").Order("id desc").Preload("Categoria").Preload("Usuario").Limit(3).Find(&datos)
	var arreglo = dto.RecetasResponse{}
	//obtenemos la URL del proyecto
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	for _, dato := range datos {
		//formateamos fecha
		fecha := fmt.Sprintf("%d/%02d/%02d", dato.Fecha.Day(), dato.Fecha.Month(), dato.Fecha.Year())
		//formeatemos la url de la foto
		foto := scheme + "://" + c.Request.Host + "/public/uploads/recetas/" + dato.Foto
		//llenamos el arreglo de tipo struct
		arreglo = append(arreglo, dto.RecetaResponse{Id: dato.Id, Nombre: dato.Nombre, Slug: dato.Slug, CategoriaDtoId: dato.CategoriaID, Categoria: dato.Categoria.Nombre, UsuarioId: dato.UsuarioID, Usuario: dato.Usuario.Nombre, Tiempo: dato.Tiempo, Descripcion: dato.Descripcion, Foto: foto, Fecha: fecha})

	}

	c.JSON(http.StatusOK, arreglo)
}
