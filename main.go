package main

import (
	"log"
	"net/http"

	"github.com/a-h/templ"

	// IMPORTANTE:
	// 1. Importa tus tipos de sqlc (para que coincidan)
	sqlc "tp3/db"
	// 2. Importa tus vistas generadas por templ
	views "tp3/Views"
)

func main() {
	// --- 1. CREA TUS DATOS FALSOS ---
	// Simplemente creamos slices de las structs que sqlc generó.

	fakeTemas := []sqlc.Tema{
		{IDTema: 1, NombreTema: "mate"},
		{IDTema: 2, NombreTema: "ing"},
		{IDTema: 3, NombreTema: "leng"},
	}

	fakeUsuarios := []sqlc.Usuario{
		{IDUsuario: 1, NombreUsuario: "Usuario de Prueba", Email: "fake@test.com"},
		{IDUsuario: 2, NombreUsuario: "Otro Usuario", Email: "otro@test.com"},
	}

	fakeTarjetas := []sqlc.Tarjetum{
		{IDTarjeta: 1, IDTema: 1, Pregunta: "¿Funciona?"},
		{IDTarjeta: 2, IDTema: 2, Pregunta: "¿entiendo?"},
	}

	// --- 2. CREA UN HANDLER BÁSICO ---
	// Este handler no va a la BBDD, solo usa los datos falsos.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// Llama a tu componente IndexPage exactamente como lo haría tu
		// handler real, pero pasándole los datos falsos.
		componente := views.IndexPage(fakeTemas, fakeUsuarios, fakeTarjetas)

		// Renderiza el componente
		templ.Handler(componente).ServeHTTP(w, r)
	})

	// --- 3. INICIA EL SERVIDOR DE PRUEBAS ---
	// Lo iniciamos en un puerto DIFERENTE (ej: 8081) para
	// que no choque con tu servidor real (que está en 8080).
	log.Println("Servidor de PRUEBAS DE UI iniciado en http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
