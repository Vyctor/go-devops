package main

import (
	"alert-manager/email"
)

func main() {
	email.SendEmail([]string{"vyctorguimaraes@gmail.com"}, "Alerta de servidor down", "Google", "Erro ao conectar no servidor", "06/06/2023 14:00", "./template.html")
}
