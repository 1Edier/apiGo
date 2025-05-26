//utils/validation.go
package utils

import (
	
	"regexp"
	"strings"
	"server/models"
)

// ErrorCampo representa un error en un campo específico
type ErrorCampo struct {
	Campo   string `json:"campo"`
	Mensaje string `json:"mensaje"`
}

// ClienteConErrores representa un cliente con sus errores de validación
type ClienteConErrores struct {
	Cliente models.Cliente  `json:"cliente"`
	Errores []ErrorCampo    `json:"errores"`
	TieneErrores bool        `json:"tieneErrores"`
}

// ValidarClave valida que la clave contenga solo números
func ValidarClave(clave string) string {
	if clave == "" {
		return "La clave es requerida"
	}
	
	// Validar que solo contenga números
	matched, _ := regexp.MatchString(`^\d+$`, clave)
	if !matched {
		return "La clave debe contener únicamente números"
	}
	
	return ""
}

// ValidarNombreContacto valida que el nombre no contenga números
func ValidarNombreContacto(nombre string) string {
	if nombre == "" {
		return "El nombre de contacto es requerido"
	}
	
	// Validar que no contenga números
	matched, _ := regexp.MatchString(`\d`, nombre)
	if matched {
		return "El nombre de contacto no debe contener números"
	}
	
	return ""
}

// ValidarCorreo valida que el correo tenga un dominio conocido
func ValidarCorreo(correo string) string {
	if correo == "" {
		return "El correo es requerido"
	}
	
	// Validar formato básico de correo
	emailRegex := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	if !emailRegex.MatchString(correo) {
		return "Formato de correo inválido"
	}
	
	// Dominios conocidos más populares
	dominiosConocidos := []string{
		"gmail.com", "hotmail.com", "outlook.com", "yahoo.com", 
		"live.com", "icloud.com", "aol.com", "protonmail.com", 
		"zoho.com", "mail.com", "yandex.com", "gmx.com",
		"tutanota.com", "fastmail.com", "rocketmail.com",
	}
	
	// Extraer dominio
	parts := strings.Split(correo, "@")
	if len(parts) != 2 {
		return "Formato de correo inválido"
	}
	
	dominio := strings.ToLower(parts[1])
	
	// Verificar si es un dominio conocido
	for _, dominioConocido := range dominiosConocidos {
		if dominio == dominioConocido {
			return ""
		}
	}
	
	// Si no es un dominio conocido, verificar si parece ser un dominio empresarial válido
	if strings.Contains(dominio, ".") {
		partesDominio := strings.Split(dominio, ".")
		if len(partesDominio) >= 2 && len(partesDominio[len(partesDominio)-1]) >= 2 {
			// Podría ser un dominio empresarial válido, pero advertimos
			return "Se recomienda usar un proveedor de correo conocido"
		}
	}
	
	return "Por favor utiliza un servicio de correo conocido"
}

// ValidarTelefono valida que el teléfono tenga exactamente 10 dígitos
func ValidarTelefono(telefono string) string {
	if telefono == "" {
		return "El teléfono es requerido"
	}
	
	// Eliminar espacios, guiones y paréntesis para validar
	telefonoLimpio := regexp.MustCompile(`[\s\-\(\)]+`).ReplaceAllString(telefono, "")
	
	// Validar que solo contenga números
	matched, _ := regexp.MatchString(`^\d+$`, telefonoLimpio)
	if !matched {
		return "El teléfono debe contener únicamente números (se permiten espacios, guiones y paréntesis)"
	}
	
	// Validar que tenga exactamente 10 dígitos (incluyendo lada)
	if len(telefonoLimpio) != 10 {
		return "El teléfono debe tener exactamente 10 dígitos incluyendo lada"
	}
	
	return ""
}

// ValidarCliente valida todos los campos de un cliente
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

// GetClientesConErrores obtiene todos los clientes y sus errores de validación
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

// GetClientesConErroresSolo obtiene solo los clientes que tienen errores
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

// BuscarClientesConErrores busca clientes y retorna con información de errores
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