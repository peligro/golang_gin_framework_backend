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
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	//go get github.com/gosimple/slug
)

/*
	func Receta_get(c *gin.Context) {
		datos := modelos.Recetas{}
		database.Database.Order("id desc").Find(&datos)
		c.JSON(http.StatusOK, datos)
	}
*/

func Receta_get(c *gin.Context) {
	datos := modelos.Recetas{}
	database.Database.Order("id desc").Preload("Categoria").Find(&datos)
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
		arreglo = append(arreglo, dto.RecetaResponse{Id: dato.Id, Nombre: dato.Nombre, Slug: dato.Slug, CategoriaDtoId: dato.CategoriaID, Categoria: dato.Categoria.Nombre, Tiempo: dato.Tiempo, Descripcion: dato.Descripcion, Foto: foto, Fecha: fecha})

	}

	c.JSON(http.StatusOK, arreglo)
}
func Receta_get_con_parametro(c *gin.Context) {
	id := c.Param("id")

	dato := modelos.Receta{}

	if err := database.Database.Preload("Categoria").First(&dato, id); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"estado":        "error",
			"mensaje":       "Recurso no disponible",
			"errorOpcional": err.Error.Error(),
		})
	} else {
		//obtenemos la URL del proyecto
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		//formateamos fecha
		fecha := fmt.Sprintf("%d/%02d/%02d", dato.Fecha.Day(), dato.Fecha.Month(), dato.Fecha.Year())
		//formeatemos la url de la foto
		foto := scheme + "://" + c.Request.Host + "/public/uploads/recetas/" + dato.Foto
		//llenamos el arreglo de tipo struct
		c.JSON(http.StatusOK, dto.RecetaResponse{Id: dato.Id, Nombre: dato.Nombre, Slug: dato.Slug, CategoriaDtoId: dato.CategoriaID, Categoria: dato.Categoria.Nombre, Tiempo: dato.Tiempo, Descripcion: dato.Descripcion, Foto: foto, Fecha: fecha})
	}

}

func Receta_post(c *gin.Context) {
	file, errFoto := c.FormFile("foto")
	//validamos los campos
	errosValidation := map[string][]string{}
	if strings.TrimSpace((c.PostForm(("nombre")))) == "" {
		errosValidation["nombre"] = append(errosValidation["nombre"], "El nombre es obligatorio")
	}
	if utf8.RuneCountInString((c.PostForm(("nombre")))) > 100 { //solo para que veas que se puede hacer
		errosValidation["nombre"] = append(errosValidation["nombre"], "El nombre no debe tener más de 100 caracteres")
	}
	if strings.TrimSpace((c.PostForm(("tiempo")))) == "" {
		errosValidation["tiempo"] = append(errosValidation["tiempo"], "El tiempo es obligatorio")
	}
	if strings.TrimSpace((c.PostForm(("descripcion")))) == "" {
		errosValidation["descripcion"] = append(errosValidation["descripcion"], "El descripcion es obligatorio")
	}
	if strings.TrimSpace((c.PostForm(("categoria_id")))) == "" {
		errosValidation["categoria_id"] = append(errosValidation["categoria_id"], "El categoria_id es obligatorio")
	}
	if len(errosValidation) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"estado":  "error",
			"mensaje": errosValidation,
		})
		return
	}
	//validamos que exista la categoría por id
	catExiste := modelos.Categoria{}
	if err := database.Database.First(&catExiste, c.PostForm("categoria_id")); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"estado":        "error",
			"mensaje":       "Recurso no disponible",
			"errorOpcional": "La categoría informada no existe", //opcional para ver el error real
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
	//validamos que no existe el nombre de receta
	existe := modelos.Recetas{}
	database.Database.Where(&modelos.Receta{Nombre: c.PostForm("nombre")}).Find(&existe)
	if len(existe) >= 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"estado":        "error",
			"mensaje":       "Recurso no disponible",
			"errorOpcional": "Ya existe el nombre: " + c.PostForm("nombre"),
		})
		return
	}

	//nombramos el archivo
	var extension = strings.Split(file.Filename, ".")[1]
	tiempo := strings.Split(time.Now().String(), " ")
	foto := string(tiempo[4][6:14]) + "." + extension
	var archivo string = "public/uploads/recetas/" + foto
	// subimos el archivo
	c.SaveUploadedFile(file, archivo)
	//creamos el registro
	categoria_id, _ := strconv.ParseUint(c.PostForm("categoria_id"), 10, 64)

	datos := modelos.Receta{
		CategoriaID: uint(categoria_id),
		Nombre:      c.PostForm("nombre"),
		Slug:        slug.Make(c.PostForm("nombre")),
		Tiempo:      c.PostForm("tiempo"),
		Foto:        foto,
		Descripcion: c.PostForm("descripcion"),
		Fecha:       time.Now()}
	database.Database.Save(&datos)
	//retornamos
	c.JSON(http.StatusCreated, gin.H{
		"estado":  "ok",
		"mensaje": "Se creó el registro exitosamente",
	})

}

func Receta_put(c *gin.Context) {
	//obtenemos los datos del json request
	var body dto.RecetaDto
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
	//validamos que exista la receta por id
	datos := modelos.Receta{}
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
	datos.Tiempo = body.Tiempo
	datos.Descripcion = body.Descripcion
	datos.CategoriaID = body.CategoriaId
	database.Database.Save(&datos)
	//retornamos
	c.JSON(http.StatusOK, gin.H{
		"estado":  "ok",
		"mensaje": "Se modificó el registro exitosamente",
	})
}
func Receta_delete(c *gin.Context) {
	//obtenemos el id path
	id := c.Param("id")
	//validamos que exista la receta por id
	datos := modelos.Receta{}
	if err := database.Database.First(&datos, id); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"estado":        "error",
			"mensaje":       "Recurso no disponible",
			"errorOpcional": err.Error.Error(), //opcional para ver el error real
		})
		return
	}
	// borramos foto
	borrar := "public/uploads/recetas/" + datos.Foto
	e := os.Remove(borrar)
	if e != nil {
		log.Fatal(e)
	}
	//borramos el registro
	database.Database.Delete(&datos)
	//retornamos
	c.JSON(http.StatusOK, gin.H{
		"estado":  "ok",
		"mensaje": "Se eliminó el registro exitosamente",
	})
}
