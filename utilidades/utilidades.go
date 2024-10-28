package utilidades

import (
	"os"

	"github.com/joho/godotenv"
	gomail "gopkg.in/gomail.v2"
)

func EnviarCorreo(correo string, asunto string, mensaje string) {
	//variables globales
	errorVariables := godotenv.Load()
	if errorVariables != nil {

		panic(errorVariables)

	}
	//enviamos el mail
	msg := gomail.NewMessage()
	msg.SetHeader("From", "test@tusitio.com")
	msg.SetHeader("To", correo)
	msg.SetHeader("Subject", asunto)
	msg.SetBody("text/html", mensaje)
	//ahora se debe configurar la conexión con el SMTP
	n := gomail.NewDialer(os.Getenv("SMTP_SERVER"), 587, os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASSWORD"))

	if err := n.DialAndSend(msg); err != nil {
		panic(err)
	}

}
