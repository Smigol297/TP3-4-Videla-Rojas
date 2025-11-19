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

// Esta función maneja la ruta / (el Home). pagina principal
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

	//1. Creamos el componente de la página principal,
	component := views.IndexPage(temas, usuarios, tarjetas)
	//ivertir ese componente en HTML y enviarlo al navegador.
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
	//cambio HTMX
	//http.Redirect(w, r, "/", http.StatusSeeOther)
	//obtengo la lista actualizada de temas
	temas, err := app.Queries.ListTemas(context.Background())
	if err != nil {
		log.Printf("Error listando temas post-creación: %v", err)
		http.Error(w, "Error al obtener lista de temas", http.StatusInternalServerError)
		return
	}

	// 2. Renderizamos y devolvemos SÓLO el componente de la lista.
	// HTMX lo recibirá y lo pondrá en el hx-target="#lista-temas".
	views.TemaList(temas).Render(r.Context(), w)
	//renderizamos el FORMULARIO VACÍO (OOB)
	// HTMX detectará el atributo hx-swap-oob="true" y buscará dónde ponerlo
	views.TemaFormOOB().Render(r.Context(), w)
	//  Actualizamos el FORMULARIO DE TARJETAS (OOB 2)
	// Al renderizar esto, HTMX buscará el id="form-tarjetas" y lo reemplazará.
	// Como le pasamos la lista 'temas' nueva, el <select> incluirá el nuevo tema.
	views.TarjetaFormOOB(temas).Render(r.Context(), w)
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
	//borro redirect para aplicar HTMX
	//http.Redirect(w, r, "/", http.StatusSeeOther)
	usuarios, err := app.Queries.ListUsuarios(context.Background())
	if err != nil {
		log.Printf("Error listando temas post-creación: %v", err)
		http.Error(w, "Error al obtener lista de temas", http.StatusInternalServerError)
		return
	}

	// 2. Renderizamos y devolvemos SÓLO el componente de la lista.
	// HTMX lo recibirá y lo pondrá en el hx-target="#usuarios-temas".
	views.UsuarioList(usuarios).Render(r.Context(), w)
	//renderizamos el FORMULARIO VACÍO (OOB)
	// HTMX detectará el atributo hx-swap-oob="true" y buscará dónde ponerlo
	views.UsuarioFormOOB().Render(r.Context(), w)

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
	//http.Redirect(w, r, "/", http.StatusSeeOther)

	//obtengo la lista actualizada de tarjetas
	tarjetas, err := app.Queries.ListTarjetas(context.Background())
	if err != nil {
		log.Printf("Error listando tarjetas post-creación: %v", err)
		http.Error(w, "Error al obtener lista de tarjetas", http.StatusInternalServerError)
		return
	}

	// 2. ¡FALTABA ESTO! Obtenemos la lista de TEMAS (Para el select del formulario)
	temas, err := app.Queries.ListTemas(context.Background())
	if err != nil {
		log.Printf("Error listando temas: %v", err)
		http.Error(w, "Error al obtener lista de temas", http.StatusInternalServerError)
		return
	}
	// Renderizamos y devolvemos SÓLO el componente de la lista.
	views.TarjetaList(tarjetas).Render(r.Context(), w)
	//renderizamos el FORMULARIO VACÍO (OOB)
	// HTMX detectará el atributo hx-swap-oob="true" y buscará dónde ponerlo
	views.TarjetaFormOOB(temas).Render(r.Context(), w)

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
