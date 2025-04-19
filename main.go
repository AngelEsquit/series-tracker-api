// @title Series Tracker API
// @version 1.0
// @description API para gestionar series de televisión
// @host localhost:8080
// @schemes http
package main

import (
	"fmt"
	"log"
	"net/http"

	_ "series-tracker-backend/docs" // Importa la documentación generada por swaggo

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	initDB()        // Inicializa la base de datos
	defer closeDB() // Asegúrate de cerrar la base de datos al final

	router := mux.NewRouter() // Crea un nuevo router de mux
	setupRoutes(router)       // Configura las rutas de la API

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler) // Configura la ruta para la documentación Swagger

	handler := enableCORS(router)         // Habilita CORS para todas las rutas
	http.ListenAndServe(":8080", handler) // Inicia el servidor en el puerto 8080

	fmt.Println("Servidor escuchando en el puerto 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
