/*package main

import (
	"net/http"
	"tp3/logic"

	_ "github.com/lib/pq"
)

func main() {
	// Manejo automatico Content-Type para archivos estàticos
	htmlContent := "./static"
	fileServer := http.FileServer(http.Dir(htmlContent))
	// Registra un manejador (handler) para la ruta raíz "/"
	http.Handle("/", fileServer)
	http.HandleFunc("/users", logic.UsersHandler) // Mapea la ruta /users a la función UsersHandler del paquete logic
	http.HandleFunc("/users/", logic.UsersByIDHandler)
	http.HandleFunc("/tarjetas", logic.TarjetasHandler) // Mapea la ruta /tarjetas a la función TarjetasHandler del paquete logic
	http.HandleFunc("/tarjetas/{id}", logic.TarjetasByIDHandler)
	http.HandleFunc("/temas", logic.TemasHandler) // Mapea la ruta /temas a la función TemasHandler del paquete logic
	http.HandleFunc("/temas/", logic.TemasByIDHandler)
	logic.InitServer()
}
*/