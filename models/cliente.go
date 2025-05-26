//models/cliente.go
package models


type Cliente struct {
	Clave            string `json:"clave"`
	NombreContacto   string `json:"nombreContacto"`
	Correo           string `json:"correo"`
	TelefonoContacto string `json:"telefonoContacto"`
}