package main

import (
	"net/http"
	"tp3/logic"

	_ "github.com/lib/pq"
)

func main() {
	http.HandleFunc("/users", logic.UsersHandler) // Mapea la ruta /users a la función UsersHandler del paquete logic
	http.HandleFunc("/users/{id}", logic.UsersByIDHandler)

	initServer()
}
