package main

import (
	"net/http"
	"tp3/logic"

	_ "github.com/lib/pq"
)

func main() {
	http.HandleFunc("/users", logic.UsersHandler) // Mapea la ruta /users a la función UsersHandler del paquete logic
	http.HandleFunc("/users/{id}", logic.UsersByIDHandler)
	http.HandleFunc("/tarjetas", logic.TarjetasHandler) // Mapea la ruta /tarjetas a la función TarjetasHandler del paquete logic
	http.HandleFunc("/tarjetas/{id}", logic.TarjetasByIDHandler)
	http.HandleFunc("/temas", logic.TemasHandler) // Mapea la ruta /temas a la función TemasHandler del paquete logic
	http.HandleFunc("/temas/{id}", logic.TemasByIDHandler)
	logic.InitServer()
}
