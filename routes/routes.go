//routes/routes.go
package routes

import (
	"github.com/gorilla/mux"
	"server/controllers"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Rutas existentes
	r.HandleFunc("/clientes", controllers.GetClientes).Methods("GET")
	r.HandleFunc("/clientes/buscar", controllers.BuscarClientes).Methods("GET")
	r.HandleFunc("/clientes/{clave}", controllers.GetClientePorClave).Methods("GET")
	r.HandleFunc("/clientes", controllers.CrearCliente).Methods("POST")
	r.HandleFunc("/clientes/{clave}", controllers.ActualizarCliente).Methods("PUT")
	r.HandleFunc("/clientes/{clave}", controllers.EliminarCliente).Methods("DELETE")
	
	// Nuevas rutas para validación
	r.HandleFunc("/clientes/validacion/todos", controllers.GetClientesConValidacion).Methods("GET")
	r.HandleFunc("/clientes/validacion/errores", controllers.GetClientesConErrores).Methods("GET")
	r.HandleFunc("/clientes/validacion/buscar", controllers.BuscarClientesConValidacion).Methods("GET")
	
	return r
}