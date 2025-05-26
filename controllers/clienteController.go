
package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"server/models"
	"server/utils"
)


func GetClientes(w http.ResponseWriter, r *http.Request) {
	clientes, err := utils.GetClientes()
	if err != nil {
		http.Error(w, "Error al obtener clientes: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clientes)
}


func GetClientesConValidacion(w http.ResponseWriter, r *http.Request) {
	clientesConErrores, err := utils.GetClientesConErrores()
	if err != nil {
		http.Error(w, "Error al obtener clientes: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clientesConErrores)
}


func GetClientesConErrores(w http.ResponseWriter, r *http.Request) {
	clientesConErrores, err := utils.GetClientesConErroresSolo()
	if err != nil {
		http.Error(w, "Error al obtener clientes con errores: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clientesConErrores)
}


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


func CrearCliente(w http.ResponseWriter, r *http.Request) {
	var cliente models.Cliente
	if err := json.NewDecoder(r.Body).Decode(&cliente); err != nil {
		http.Error(w, "Error al decodificar cliente: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	
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


func ActualizarCliente(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clave := vars["clave"]
	
	var cliente models.Cliente
	if err := json.NewDecoder(r.Body).Decode(&cliente); err != nil {
		http.Error(w, "Error al decodificar cliente: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	
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