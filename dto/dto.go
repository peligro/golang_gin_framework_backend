package dto

//validaciones https://gin-gonic.com/docs/examples/binding-and-validation/

type EjemploDto struct {
	Correo   string `json:"correo"`
	Password string `json:"password"`
}
type CategoriaDto struct {
	Nombre string `json:"nombre" binding:"required"`
}

type RecetaDto struct {
	Nombre      string `json:"nombre" binding:"required"`
	Tiempo      string `json:"tiempo" binding:"required"`
	Descripcion string `json:"descripcion" binding:"required"`
	CategoriaId uint   `json:"categoria_id" binding:"required"`
}

type RecetaResponse struct {
	Id             uint   `json:"id"`
	Nombre         string `json:"nombre"`
	Slug           string `json:"slug"`
	CategoriaDtoId uint   `json:"categoria_id"`
	Categoria      string `json:"categoria"`
	Tiempo         string `json:"tiempo"`
	Descripcion    string `json:"descripcion"`
	Foto           string `json:"foto"`
	Fecha          string `json:"fecha"`
}

type RecetasResponse []RecetaResponse
