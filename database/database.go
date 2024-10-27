package database

//go get gorm.io/driver/mysql
//go get gorm.io/gorm
import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Database = func() (db *gorm.DB) {
	//se valida existencia .env y variables de entorno
	errorVariables := godotenv.Load()
	if errorVariables != nil {

		panic(errorVariables)
		return
	}
	//dsn := "root123456@tcp(localhost:3306)/golang1?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
	if db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		fmt.Println("Error de conexión")
		panic(err)
	} else {
		fmt.Println("Conexión exitosa")
		return db
	}
}()
