package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	sqlc "tp3/db"

	_ "github.com/lib/pq"
)

// INICIO CON MINUSCULA = NO PUBLICO
// INICIO CON MAYUSCULA = PUBLICO
func ValidateCreateTarjeta(p sqlc.CreateTarjetaParams) error {
	if p.Pregunta == "" {
		return fmt.Errorf("La pregunta no puede estar vacía")
	}
	if p.Respuesta == "" {
		return fmt.Errorf("La respuesta no puede estar vacía")
	}
	if p.OpcionA == "" {
		return fmt.Errorf("La opción A no puede estar vacía")
	}
	if p.OpcionB == "" {
		return fmt.Errorf("La opción B no puede estar vacía")
	}
	if p.OpcionC == "" {
		return fmt.Errorf("La opción C no puede estar vacía")
	}
	if p.IDTema <= 0 {
		return fmt.Errorf("ID de tema inválido")
	}
	return nil
}
func ValidateUpdateTarjeta(p sqlc.UpdateTarjetaParams) error {
	if p.IDTarjeta <= 0 {
		return fmt.Errorf("ID de tarjeta %d inválido", p.IDTarjeta)
	}
	if p.Pregunta == "" {
		return fmt.Errorf("La pregunta no puede estar vacía")
	}
	if p.Respuesta == "" {
		return fmt.Errorf("La respuesta no puede estar vacía")
	}
	if p.OpcionA == "" {
		return fmt.Errorf("La opción A no puede estar vacía")
	}
	if p.OpcionB == "" {
		return fmt.Errorf("La opción B no puede estar vacía")
	}
	if p.OpcionC == "" {
		return fmt.Errorf("La opción C no puede estar vacía")
	}
	if p.IDTema <= 0 {
		return fmt.Errorf("ID de tema inválido")
	}
	return nil
}

func TarjetasHandler(w http.ResponseWriter, r *http.Request) {
	db := connectDB()
	defer db.Close()
	queries := sqlc.New(db)
	switch r.Method {
	case http.MethodGet:
		temaStr := r.URL.Query().Get("tema") // obtiene el parámetro "tema"
		if temaStr == "" {
			getTarjetas(w, r, queries)
			return
		}

		tema, err := strconv.Atoi(temaStr)
		if err != nil {
			http.Error(w, "ID de tarjeta inválido", http.StatusBadRequest)
			return
		}

		getTarjetasByTema(w, r, tema, queries)

	case http.MethodPost:
		createTarjeta(w, r, queries)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func getTarjetas(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	w.Header().Set("Content-Type", "application/json")

	ctx := context.Background()
	tarjetas, err := queries.ListTarjetas(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Enviar los datos como JSON válido
	if err := json.NewEncoder(w).Encode(tarjetas); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func getTarjetasByTema(w http.ResponseWriter, r *http.Request, tema int, queries *sqlc.Queries) {
	w.Header().Set("Content-Type", "application/json")

	ctx := context.Background()
	tarjetas, err := queries.ListTarjetasByTema(ctx, int32(tema))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Enviar los datos como JSON válido
	if err := json.NewEncoder(w).Encode(tarjetas); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func createTarjeta(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	var p sqlc.CreateTarjetaParams

	// decodificar JSON
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// validar tarjeta
	if err := ValidateCreateTarjeta(p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// crear tarjeta en la base de datos
	ctx := context.Background()
	newTarjeta, err := queries.CreateTarjeta(ctx, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// Enviar los datos como JSON válido
	if err := json.NewEncoder(w).Encode(newTarjeta); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func TarjetasByIDHandler(w http.ResponseWriter, r *http.Request) {
	db := connectDB()
	defer db.Close()
	queries := sqlc.New(db)
	// Extraer el ID del usuario de la URL
	var id int
	_, err := fmt.Sscanf(r.URL.Path, "/tarjetas/%d", &id)
	if err != nil {
		http.Error(w, "ID de tarjeta inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getTarjetaByID(w, r, id, queries)
	case http.MethodPut:
		putTarjetaByID(w, r, id, queries)
	case http.MethodDelete:
		deleteTarjetaByID(w, r, id, queries)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}
func getTarjetaByID(w http.ResponseWriter, r *http.Request, id int, queries *sqlc.Queries) {
	ctx := context.Background()
	tarjeta, err := queries.GetTarjetaById(ctx, int32(id))
	if err != nil {
		http.Error(w, "Tarjeta no encontrada", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// Enviar los datos como JSON válido
	if err := json.NewEncoder(w).Encode(tarjeta); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func putTarjetaByID(w http.ResponseWriter, r *http.Request, id int, queries *sqlc.Queries) {
	var p sqlc.UpdateTarjetaParams
	p.IDTarjeta = int32(id)
	// decodificar JSON
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	// validar usuario
	if err := ValidateUpdateTarjeta(p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// actualizar tarjeta en la base de datos
	ctx := context.Background()
	err = queries.UpdateTarjeta(ctx, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Enviar los datos como JSON válido
	if err := json.NewEncoder(w).Encode(p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func deleteTarjetaByID(w http.ResponseWriter, r *http.Request, id int, queries *sqlc.Queries) {
	ctx := context.Background()
	err := queries.DeleteTarjeta(ctx, int32(id))
	if err != nil {
		http.Error(w, "Tarjeta no encontrada", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
