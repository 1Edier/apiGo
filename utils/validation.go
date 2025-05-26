
package utils

import (
	"regexp"
	"strings"
	"server/models"
)


type ErrorCampo struct {
	Campo   string `json:"campo"`
	Mensaje string `json:"mensaje"`
}


type ClienteConErrores struct {
	Cliente models.Cliente  `json:"cliente"`
	Errores []ErrorCampo    `json:"errores"`
	TieneErrores bool        `json:"tieneErrores"`
}


func ValidarClave(clave string) string {
	if clave == "" {
		return "La clave es requerida"
	}
	
	
	matched, _ := regexp.MatchString(`^\d+$`, clave)
	if !matched {
		return "La clave debe contener únicamente números"
	}
	
	return ""
}


func ValidarNombreContacto(nombre string) string {
	if nombre == "" {
		return "El nombre de contacto es requerido"
	}
	
	
	matched, _ := regexp.MatchString(`\d`, nombre)
	if matched {
		return "El nombre de contacto no debe contener números"
	}
	
	return ""
}


func ValidarCorreo(correo string) string {
	if correo == "" {
		return "El correo es requerido"
	}
	
	
	emailRegex := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	if !emailRegex.MatchString(correo) {
		return "Formato de correo inválido"
	}
	

	dominiosConocidos := []string{
		"gmail.com", "hotmail.com", "outlook.com", "yahoo.com", 
		"live.com", "icloud.com", "aol.com", "protonmail.com", 
		"zoho.com", "mail.com", "yandex.com", "gmx.com",
		"tutanota.com", "fastmail.com", "rocketmail.com",
	}
	

	parts := strings.Split(correo, "@")
	if len(parts) != 2 {
		return "Formato de correo inválido"
	}
	
	dominio := strings.ToLower(parts[1])
	

	for _, dominioConocido := range dominiosConocidos {
		if dominio == dominioConocido {
			return ""
		}
	}
	
	
	if strings.Contains(dominio, ".") {
		partesDominio := strings.Split(dominio, ".")
		if len(partesDominio) >= 2 && len(partesDominio[len(partesDominio)-1]) >= 2 {
			
			return "Se recomienda usar un proveedor de correo conocido"
		}
	}
	
	return "Por favor utiliza un servicio de correo conocido"
}


func ValidarTelefono(telefono string) string {
	if telefono == "" {
		return "El teléfono es requerido"
	}
	
	
	telefonoLimpio := regexp.MustCompile(`[\s\-\(\)]+`).ReplaceAllString(telefono, "")
	
	
	matched, _ := regexp.MatchString(`^\d+$`, telefonoLimpio)
	if !matched {
		return "El teléfono debe contener únicamente números (se permiten espacios, guiones y paréntesis)"
	}
	
	
	if len(telefonoLimpio) < 10 {
		return "El teléfono debe tener al menos 10 dígitos incluyendo lada"
	}
	
	
	codigosChiapas := []string{
		"961", // Tuxtla Gutiérrez
		"962", // San Cristóbal de las Casas, Comitán
		"963", // Tapachula
		"964", // Palenque
		"965", // Tonalá
		"966", // Arriaga, Pijijiapan
		"967", // Las Margaritas, Altamirano
		"968", // Villaflores, Villa Corzo
		"992", // Reforma (parte de Chiapas)
		"994", // Mapastepec
	}
	
	
	tieneCodigoChiapas := false
	codigoEncontrado := ""
	
	for _, codigo := range codigosChiapas {
		if strings.HasPrefix(telefonoLimpio, codigo) {
			tieneCodigoChiapas = true
			codigoEncontrado = codigo
			break
		}
	}
	
	if !tieneCodigoChiapas {
		return "El teléfono debe ser de Chiapas. Códigos válidos: 961 (Tuxtla Gutiérrez), 962 (San Cristóbal), 963 (Tapachula), 964 (Palenque), 965 (Tonalá), 966 (Arriaga), 967 (Las Margaritas), 968 (Villaflores), 992 (Reforma), 994 (Mapastepec)"
	}
	
	
	longitudEsperada := 10
	if len(telefonoLimpio) != longitudEsperada {
		return "El teléfono debe tener exactamente 10 dígitos (3 de lada + 7 de número local)"
	}
	
	
	numeroLocal := telefonoLimpio[3:]
	if len(numeroLocal) != 7 {
		return "Después del código de área (" + codigoEncontrado + ") debe haber exactamente 7 dígitos"
	}
	
	return ""
}


func ValidarCliente(cliente models.Cliente) []ErrorCampo {
	var errores []ErrorCampo
	
	if err := ValidarClave(cliente.Clave); err != "" {
		errores = append(errores, ErrorCampo{Campo: "clave", Mensaje: err})
	}
	
	if err := ValidarNombreContacto(cliente.NombreContacto); err != "" {
		errores = append(errores, ErrorCampo{Campo: "nombreContacto", Mensaje: err})
	}
	
	if err := ValidarCorreo(cliente.Correo); err != "" {
		errores = append(errores, ErrorCampo{Campo: "correo", Mensaje: err})
	}
	
	if err := ValidarTelefono(cliente.TelefonoContacto); err != "" {
		errores = append(errores, ErrorCampo{Campo: "telefonoContacto", Mensaje: err})
	}
	
	return errores
}


func GetClientesConErrores() ([]ClienteConErrores, error) {
	clientes, err := GetClientes()
	if err != nil {
		return nil, err
	}
	
	var clientesConErrores []ClienteConErrores
	
	for _, cliente := range clientes {
		errores := ValidarCliente(cliente)
		clienteConError := ClienteConErrores{
			Cliente:      cliente,
			Errores:      errores,
			TieneErrores: len(errores) > 0,
		}
		clientesConErrores = append(clientesConErrores, clienteConError)
	}
	
	return clientesConErrores, nil
}


func GetClientesConErroresSolo() ([]ClienteConErrores, error) {
	clientesConErrores, err := GetClientesConErrores()
	if err != nil {
		return nil, err
	}
	
	var soloConErrores []ClienteConErrores
	for _, cliente := range clientesConErrores {
		if cliente.TieneErrores {
			soloConErrores = append(soloConErrores, cliente)
		}
	}
	
	return soloConErrores, nil
}


func BuscarClientesConErrores(termino string) ([]ClienteConErrores, error) {
	clientes, err := BuscarClientes(termino)
	if err != nil {
		return nil, err
	}
	
	var clientesConErrores []ClienteConErrores
	
	for _, cliente := range clientes {
		errores := ValidarCliente(cliente)
		clienteConError := ClienteConErrores{
			Cliente:      cliente,
			Errores:      errores,
			TieneErrores: len(errores) > 0,
		}
		clientesConErrores = append(clientesConErrores, clienteConError)
	}
	
	return clientesConErrores, nil
}