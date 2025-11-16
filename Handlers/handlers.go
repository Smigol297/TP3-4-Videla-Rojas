package handlers

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strconv"
	views "tp3/Views"
	sqlc "tp3/db"
	"tp3/logic"

	"github.com/a-h/templ"
)

// Inyección de Dependencias
type Application struct {
	DB      *sql.DB //conector a la BD
	Queries *sqlc.Queries
}

func (app *Application) HandleGetIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		if r.Method == http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		http.NotFound(w, r)
		return
	}

	ctx := context.Background()
	temas, err := app.Queries.ListTemas(ctx)
	if err != nil {
		log.Printf("Error listando temas: %v", err)
	}
	usuarios, err := app.Queries.ListUsuarios(ctx)
	if err != nil {
		log.Printf("Error listando usuarios: %v", err)
	}
	tarjetas, err := app.Queries.ListTarjetas(ctx)
	if err != nil {
		log.Printf("Error listando tarjetas: %v", err)
	}

	component := views.IndexPage(temas, usuarios, tarjetas)
	templ.Handler(component).ServeHTTP(w, r)
}

// --- Handlers de Creación (POST) ---
func (app *Application) HandleCreateTema(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear formulario", http.StatusBadRequest)
		return
	}
	nombreTema := r.FormValue("nombre_tema")
	if err := logic.ValidateCreateTema(nombreTema); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err := app.Queries.CreateTema(context.Background(), nombreTema)
	if err != nil {
		http.Error(w, "Error al crear tema", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) HandleCreateUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear formulario", http.StatusBadRequest)
		return
	}
	params := sqlc.CreateUsuarioParams{
		NombreUsuario: r.FormValue("nombre_usuario"),
		Email:         r.FormValue("email"),
		Contrasena:    r.FormValue("contrasena"),
	}
	if err := logic.ValidateCreateUser(params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err := app.Queries.CreateUsuario(context.Background(), params)
	if err != nil {
		http.Error(w, "Error al crear usuario", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) HandleCreateTarjeta(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear formulario", http.StatusBadRequest)
		return
	}
	idTema, _ := strconv.Atoi(r.FormValue("id_tema"))
	params := sqlc.CreateTarjetaParams{
		Pregunta:  r.FormValue("pregunta"),
		Respuesta: r.FormValue("respuesta"),
		OpcionA:   r.FormValue("opcion_a"),
		OpcionB:   r.FormValue("opcion_b"),
		OpcionC:   r.FormValue("opcion_c"),
		IDTema:    int32(idTema),
	}
	if err := logic.ValidateCreateTarjeta(params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err := app.Queries.CreateTarjeta(context.Background(), params)
	if err != nil {
		http.Error(w, "Error al crear tarjeta", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// --- Handlers de Eliminación (POST) ---
func (app *Application) HandleDeleteTema(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear formulario", http.StatusBadRequest)
		return
	}
	id, _ := strconv.Atoi(r.FormValue("id"))
	if id <= 0 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	app.Queries.DeleteTema(context.Background(), int32(id))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) HandleDeleteUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear formulario", http.StatusBadRequest)
		return
	}
	id, _ := strconv.Atoi(r.FormValue("id"))
	if id <= 0 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	app.Queries.DeleteUsuario(context.Background(), int32(id))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) HandleDeleteTarjeta(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear formulario", http.StatusBadRequest)
		return
	}
	id, _ := strconv.Atoi(r.FormValue("id"))
	if id <= 0 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	app.Queries.DeleteTarjeta(context.Background(), int32(id))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
