//main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"server/routes"
	"server/utils"
)

func main() {
	
	if err := utils.InitExcel(); err != nil {
		log.Fatalf("Error al inicializar el archivo Excel: %v", err)
	}
	
	
	r := routes.SetupRoutes()
	

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"}) // Permite cualquier origen, ajustar en producción
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	

	port := 8080
	fmt.Printf("Servidor iniciado en http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), 
		handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}