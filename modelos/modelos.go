package modelos

import (
	"backend/database"
	"time"
)

type Categoria struct {
	Id     uint   `json:"id"`
	Nombre string `gorm:"type:varchar(100)" json:"nombre"`
	Slug   string `gorm:"type:varchar(100)" json:"slug"`
}
type Categorias []Categoria

type Receta struct {
	Id          uint      `json:"id"`
	CategoriaID uint      `json:"categoria_id"`
	Categoria   Categoria `json:"categoria"`
	Nombre      string    `gorm:"type:varchar(100)" json:"nombre"`
	Slug        string    `gorm:"type:varchar(100)" json:"slug"`
	Tiempo      string    `gorm:"type:varchar(100)" json:"tiempo"`
	Foto        string    `gorm:"type:varchar(100)" json:"foto"`
	Descripcion string    `json:"descripcion"`
	Fecha       time.Time `json:"fecha"`
}
type Recetas []Receta

func Migraciones() {
	database.Database.AutoMigrate(&Categoria{}, &Receta{})

}
