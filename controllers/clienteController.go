//controller/clienteController
package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"server/models"
	"server/utils"
)

// GetClientes maneja la petición para obtener todos los clientes
func GetClientes(w http.ResponseWriter, r *http.Request) {
	clientes, err := utils.GetClientes()
	if err != nil {
		http.Error(w, "Error al obtener clientes: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clientes)
}

// GetClientesConValidacion obtiene todos los clientes con información de validación
func GetClientesConValidacion(w http.ResponseWriter, r *http.Request) {
	clientesConErrores, err := utils.GetClientesConErrores()
	if err != nil {
		http.Error(w, "Error al obtener clientes: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clientesConErrores)
}

// GetClientesConErrores obtiene solo los clientes que tienen errores de validación
func GetClientesConErrores(w http.ResponseWriter, r *http.Request) {
	clientesConErrores, err := utils.GetClientesConErroresSolo()
	if err != nil {
		http.Error(w, "Error al obtener clientes con errores: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clientesConErrores)
}

// BuscarClientes maneja la petición para buscar clientes
func BuscarClientes(w http.ResponseWriter, r *http.Request) {
	termino := r.URL.Query().Get("q")
	if termino == "" {
		http.Error(w, "El parámetro 'q' de búsqueda es requerido", http.StatusBadRequest)
		return
	}
	
	clientes, err := utils.BuscarClientes(termino)
	if err != nil {
		http.Error(w, "Error al buscar clientes: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clientes)
}

// BuscarClientesConValidacion busca clientes y retorna con información de validación
func BuscarClientesConValidacion(w http.ResponseWriter, r *http.Request) {
	termino := r.URL.Query().Get("q")
	if termino == "" {
		http.Error(w, "El parámetro 'q' de búsqueda es requerido", http.StatusBadRequest)
		return
	}
	
	clientesConErrores, err := utils.BuscarClientesConErrores(termino)
	if err != nil {
		http.Error(w, "Error al buscar clientes: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clientesConErrores)
}

// GetClientePorClave maneja la petición para obtener un cliente por clave
func GetClientePorClave(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clave := vars["clave"]
	
	cliente, err := utils.GetClientePorClave(clave)
	if err != nil {
		http.Error(w, "Error al obtener cliente: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	if cliente == nil {
		http.Error(w, "Cliente no encontrado", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cliente)
}

// CrearCliente maneja la petición para crear un nuevo cliente
func CrearCliente(w http.ResponseWriter, r *http.Request) {
	var cliente models.Cliente
	if err := json.NewDecoder(r.Body).Decode(&cliente); err != nil {
		http.Error(w, "Error al decodificar cliente: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	// Validar datos del cliente
	errores := utils.ValidarCliente(cliente)
	if len(errores) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Errores de validación",
			"errores": errores,
		})
		return
	}
	
	if err := utils.CrearCliente(cliente); err != nil {
		http.Error(w, "Error al crear cliente: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Cliente creado con éxito"})
}

// ActualizarCliente maneja la petición para actualizar un cliente existente
func ActualizarCliente(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clave := vars["clave"]
	
	var cliente models.Cliente
	if err := json.NewDecoder(r.Body).Decode(&cliente); err != nil {
		http.Error(w, "Error al decodificar cliente: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	// Validar datos del cliente
	errores := utils.ValidarCliente(cliente)
	if len(errores) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Errores de validación",
			"errores": errores,
		})
		return
	}
	
	if err := utils.ActualizarCliente(clave, cliente); err != nil {
		http.Error(w, "Error al actualizar cliente: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Cliente actualizado con éxito"})
}

// EliminarCliente maneja la petición para eliminar un cliente
func EliminarCliente(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clave := vars["clave"]
	
	if err := utils.EliminarCliente(clave); err != nil {
		http.Error(w, "Error al eliminar cliente: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Cliente eliminado con éxito"})
}