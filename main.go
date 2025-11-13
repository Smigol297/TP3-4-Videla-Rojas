package main

import (
	"log"
	"net/http"
	handlers "tp3/Handlers" // ¡TU NUEVO PAQUETE DE HANDLERS!
	sqlc "tp3/db"           // Tu paquete sqlc
	"tp3/logic"             // Tu paquete de lógica

	_ "github.com/lib/pq"
)

// main.go AHORA SÓLO HACE DOS COSAS:
// 1. Configurar las dependencias (BBDD)
// 2. Registrar las rutas (le dice qué handler va con qué ruta)

func main() {
	// 1. Configurar Dependencias
	dbConn := logic.ConnectDB()
	if dbConn == nil {
		log.Fatal("No se pudo conectar a la BBDD")
	}
	defer dbConn.Close()

	// Creamos la instancia de sqlc.Queries
	queries := sqlc.New(dbConn)

	// Creamos la instancia de Application (que ahora vive en el paquete 'handlers')
	// y le "inyectamos" las dependencias (la BBDD y las queries)
	app := &handlers.Application{
		DB:      dbConn,
		Queries: queries,
	}

	// 2. Registrar Rutas
	//    ¡Observa qué limpio! 'main' solo define las rutas.
	//    Toda la lógica está dentro de 'app'.
	http.HandleFunc("/", app.HandleGetIndex)

	http.HandleFunc("/temas", app.HandleCreateTema)
	http.HandleFunc("/users", app.HandleCreateUsuario)
	http.HandleFunc("/tarjetas", app.HandleCreateTarjeta)

	http.HandleFunc("/temas/delete", app.HandleDeleteTema)
	http.HandleFunc("/users/delete", app.HandleDeleteUsuario)
	http.HandleFunc("/tarjetas/delete", app.HandleDeleteTarjeta)

	// 3. Iniciar el servidor (El antiguo 'InitServer')
	log.Println("Iniciando servidor en http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}
