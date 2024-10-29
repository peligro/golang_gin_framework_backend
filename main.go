package main

//snap install go --classic
//go get github.com/pilu/fresh
//go run github.com/pilu/fresh
//go run main.go
import (
	"backend/middleware"
	"backend/modelos"
	"backend/rutas"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var prefijo = "/api/v1/"

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	//migrar la bd
	modelos.Migraciones()
	//archivos estáticos
	router.Static("/public", "./public")

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"estado": "ok", "mensaje": "Hola mundo desde Golang con Gin Framework con GORM ORM"})
	})

	//custom error 404
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"estado": "error", "message": "Recurso no disponible"})
	})
	//rutas
	router.GET(prefijo+"ejemplo", rutas.Ejemplo_get)
	router.POST(prefijo+"ejemplo", rutas.Ejemplo_post)
	router.GET(prefijo+"ejemplo/:id", rutas.Ejemplo_get_con_parametro)
	router.PUT(prefijo+"ejemplo/:id", rutas.Ejemplo_put)
	router.DELETE(prefijo+"ejemplo/:id", rutas.Ejemplo_delete)
	router.GET(prefijo+"query-string", rutas.Ejemplo_query_string)
	router.POST(prefijo+"upload", rutas.Ejemplo_upload)

	router.GET(prefijo+"categorias", rutas.Categoria_get)
	router.GET(prefijo+"categorias/:id", rutas.Categoria_get_con_parametro)
	router.POST(prefijo+"categorias", middleware.ValidarJWTMiddleware, rutas.Categoria_post)
	router.PUT(prefijo+"categorias/:id", middleware.ValidarJWTMiddleware, rutas.Categoria_put)
	router.DELETE(prefijo+"categorias/:id", middleware.ValidarJWTMiddleware, rutas.Categoria_delete)

	router.GET(prefijo+"recetas", rutas.Receta_get)
	router.GET(prefijo+"recetas/:id", rutas.Receta_get_con_parametro)
	router.POST(prefijo+"recetas", middleware.ValidarJWTMiddleware, rutas.Receta_post)
	router.PUT(prefijo+"recetas/:id", middleware.ValidarJWTMiddleware, rutas.Receta_put)
	router.DELETE(prefijo+"recetas/:id", middleware.ValidarJWTMiddleware, rutas.Receta_delete)

	router.POST(prefijo+"contactanos", rutas.Contactanos_post)

	router.POST(prefijo+"seguridad/registro", rutas.Seguridad_registro)
	router.GET(prefijo+"seguridad/verificacion/:token", rutas.Seguridad_verificacion)
	router.POST(prefijo+"seguridad/login", rutas.Seguridad_login)

	//variables globales
	errorVariables := godotenv.Load()
	if errorVariables != nil {

		panic(errorVariables)

	}
	//inicio servidor
	router.Run(":" + os.Getenv("PORT"))
}

/*
import "github.com/gin-gonic/gin"
var prefijo = "/api/v1/"
func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"estado": "ok", "mensaje":"Hola mundo desde Golang con Gin Framework"})
	})
	router.Run(":8085")
}*/
