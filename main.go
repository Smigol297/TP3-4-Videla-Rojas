package main

import (
	"net/http"
	"tp3/logic"

	_ "github.com/lib/pq"
)

func main() {
	// Servidor de archivos estáticos
	htmlContent := "./static"
	fileServer := http.FileServer(http.Dir(htmlContent))
	http.Handle("/", fileServer)

	// Los endpoints de API mantienen su funcionalidad JSON
	http.HandleFunc("/users", logic.UsersHandler)
	http.HandleFunc("/users/", logic.UsersByIDHandler)
	http.HandleFunc("/tarjetas", logic.TarjetasHandler) // Maneja tanto JSON como HTML
	http.HandleFunc("/tarjetas/", logic.TarjetasByIDHandler)
	http.HandleFunc("/temas", logic.TemasHandler)
	http.HandleFunc("/temas/", logic.TemasByIDHandler)

	logic.InitServer()
}
