package modelos

import (
	"backend/database"
	"time"
)

type Categoria struct {
	Id     uint   `json:"id"` //bigint
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

type Contacto struct {
	Id       uint      `json:"id"`
	Nombre   string    `gorm:"type:varchar(100)" json:"nombre"`
	Correo   string    `gorm:"type:varchar(100)" json:"correo"`
	Telefono string    `gorm:"type:varchar(50)" json:"telefono"`
	Mensaje  string    `json:"descripcion"`
	Fecha    time.Time `json:"fecha"`
}
type Contactos []Contacto

type Estado struct {
	Id     uint   `json:"id"`
	Nombre string `gorm:"type:varchar(50)" json:"nombre"`
}

type Estados []Estado

type Usuario struct {
	Id       uint      `json:"id"`
	EstadoID uint      `json:"estado_id"`
	Estado   Estado    `json:"estado"`
	Nombre   string    `gorm:"type:varchar(100)" json:"nombre"`
	Correo   string    `gorm:"type:varchar(100)" json:"correo"`
	Password string    `gorm:"type:varchar(160)" json:"password"`
	Token    string    `gorm:"type:varchar(100)" json:"token"`
	Fecha    time.Time `json:"fecha"`
}

type Usuarios []Usuario

func Migraciones() {
	database.Database.AutoMigrate(&Categoria{}, &Receta{}, &Contacto{}, &Estado{}, &Usuario{})

}
