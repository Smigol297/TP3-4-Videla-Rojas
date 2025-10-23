package main

import (
	"context"
	"fmt"
	"net/http"
	sqlc "tp3/db"

	_ "github.com/lib/pq"
)

func main() {
	db := connectDB()
	defer db.Close()

	queries := sqlc.New(db)
	ctx := context.Background()

	http.HandleFunc("/listUsuarios", func(w http.ResponseWriter, r *http.Request) {
		usuarios, err := queries.ListUsuarios(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Renderizar la lista de usuarios en la respuesta
		fmt.Fprintf(w, "Usuarios: %+v", usuarios)
	})
	http.HandleFunc("/createUsuario", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		// Aquí deberías parsear el cuerpo de la solicitud para obtener los datos del nuevo usuario
		newUsuario := sqlc.CreateUsuarioParams{
			NombreUsuario: "Nuevo Usuario",
			Email:         "mail@example.com",
			Contrasena:    "password123",
		}
		_, err := queries.CreateUsuario(ctx, newUsuario)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	initServer()
}
