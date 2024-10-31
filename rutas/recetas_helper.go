package rutas

import (
	"backend/database"
	"backend/dto"
	"backend/modelos"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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
	categoria_id, _ := strconv.ParseUint(c.Query("categoria_id"), 10, 64)
	//validamos que exista el usuario_id
	cat := modelos.Categoria{}
	if err := database.Database.First(&cat, categoria_id); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"estado":        "error",
			"mensaje":       "Recurso no disponible",
			"errorOpcional": err.Error.Error(), //opcional para ver el error real
		})
		return
	}
	datos := modelos.Recetas{}
	database.Database.Where("nombre LIKE ?", "%"+c.Query("search")+"%").Where(&modelos.Receta{CategoriaID: cat.Id}).Order("id desc").Preload("Categoria").Preload("Usuario").Find(&datos)
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
func Receta_Helper_Editar_Foto(c *gin.Context) {
	//obtenemos el valor de
	file, errFoto := c.FormFile("foto")
	//validamos los campos
	errosValidation := map[string][]string{}

	if strings.TrimSpace((c.PostForm(("receta_id")))) == "" {
		errosValidation["receta_id"] = append(errosValidation["receta_id"], "El campo receta_id es obligatorio")
	}
	if len(errosValidation) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"estado":  "error",
			"mensaje": errosValidation,
		})
		return
	}
	//validamos que exista la receta por id
	receta := modelos.Receta{}
	if err := database.Database.First(&receta, c.PostForm("receta_id")); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"estado":        "error",
			"mensaje":       "Recurso no disponible",
			"errorOpcional": "La receta informada no existe", //opcional para ver el error real
		})
		return
	}

	//validamos el archivo
	if errFoto != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"estado":  "error",
			"mensaje": "Ocurrió un error inesperado",
		})
		return
	}

	//.Println(file.Header["Content-Type"][0])
	if file.Header["Content-Type"][0] == "image/jpeg" || file.Header["Content-Type"][0] == "image/png" {

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"estado":        "error",
			"mensaje":       "Recurso no disponible",
			"errorOpcional": "El archivo de la foto no es compatible, debe ser JPG o PNG",
		})
		return
	}
	// borramos foto
	borrar := "public/uploads/recetas/" + receta.Foto
	e := os.Remove(borrar)
	if e != nil {
		log.Fatal(e)
	}
	//nombramos el archivo
	var extension = strings.Split(file.Filename, ".")[1]
	tiempo := strings.Split(time.Now().String(), " ")
	foto := string(tiempo[4][6:14]) + "." + extension
	var archivo string = "public/uploads/recetas/" + foto
	// subimos el archivo
	c.SaveUploadedFile(file, archivo)
	//creamos el registro
	receta.Foto = foto
	database.Database.Save(&receta)
	//retornamos
	c.JSON(http.StatusOK, gin.H{
		"estado":  "ok",
		"mensaje": "Se modificó el registro exitosamente",
	})

}
