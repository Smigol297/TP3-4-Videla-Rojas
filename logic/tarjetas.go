package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	sqlc "tp3/db"

	_ "github.com/lib/pq"
)

// INICIO CON MINUSCULA = NO PUBLICO
// INICIO CON MAYUSCULA = PUBLICO

type Resultado struct {
	Pregunta          string
	RespuestaUsuario  string
	RespuestaCorrecta string
	EsCorrecta        bool
}

func ValidateCreateTarjeta(p sqlc.CreateTarjetaParams) error {
	if p.Pregunta == "" {
		return fmt.Errorf("la pregunta no puede estar vacía")
	}
	if p.Respuesta == "" {
		return fmt.Errorf("la respuesta no puede estar vacía")
	}
	if p.OpcionA == "" {
		return fmt.Errorf("la opción A no puede estar vacía")
	}
	if p.OpcionB == "" {
		return fmt.Errorf("la opción B no puede estar vacía")
	}
	if p.OpcionC == "" {
		return fmt.Errorf("la opción C no puede estar vacía")
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
		return fmt.Errorf("la pregunta no puede estar vacía")
	}
	if p.Respuesta == "" {
		return fmt.Errorf("la respuesta no puede estar vacía")
	}
	if p.OpcionA == "" {
		return fmt.Errorf("la opción A no puede estar vacía")
	}
	if p.OpcionB == "" {
		return fmt.Errorf("la opción B no puede estar vacía")
	}
	if p.OpcionC == "" {
		return fmt.Errorf("la opción C no puede estar vacía")
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
	temaStr := r.URL.Query().Get("tema") // obtiene el parámetro "tema"
	switch r.Method {
	case http.MethodGet:
		//GET /tarjetas
		if temaStr == "" {
			getTarjetas(w, r, queries)
			return
		}

		//GET /tarjetas?tema=1
		TarjetasHTMLHandler(w, r, queries)
		//getTarjetasByTema(w, r, tema, queries)
	//POST /tarjetas
	case http.MethodPost:
		if temaStr != "" {
			TarjetasHTMLHandler(w, r, queries)
			return
		}
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
func TarjetasHTMLHandler(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	temaStr := r.URL.Query().Get("tema")
	tema, err := strconv.Atoi(temaStr)
	switch r.Method {
	case http.MethodGet:
		if err != nil {
			http.Error(w, "ID de tarjeta inválido", http.StatusBadRequest)
			return
		}
		getTarjetasByTema(w, r, tema, queries)

	case http.MethodPost:
		verificarRespuestasHandler(w, r, queries, tema)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}
func getTarjetasByTema(w http.ResponseWriter, r *http.Request, tema int, queries *sqlc.Queries) {
	ctx := context.Background()
	tarjetas, err := queries.ListTarjetasByTema(ctx, int32(tema))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cambiar el content type a HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Generar el HTML directamente
	html := generarHTMLTarjetas(tarjetas)
	w.Write([]byte(html))
}

// Función para generar el HTML de las tarjetas
// Función para generar el HTML completo de las tarjetas
func generarHTMLTarjetas(tarjetas []sqlc.Tarjetum) string {
	var html strings.Builder

	html.WriteString(`<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Tarjetas de Estudio</title>
    <style>
        .tarjeta {
            display: none;
            border: 1px solid #333;
            padding: 20px;
            margin: 20px;
            width: 300px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .visible {
            display: block;
        }
        .campo-texto {
            display: block;
            margin-top: 10px;
            width: 100%;
            padding: 8px;
            border: 1px solid #ccc;
            border-radius: 4px;
        }
        .boton {
            margin-top: 15px;
            padding: 10px 15px;
            background-color: #007bff;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
        }
        .boton:hover {
            background-color: #0056b3;
        }
        .pregunta {
            font-weight: bold;
            margin-bottom: 15px;
        }
        .opcion {
            margin: 5px 0;
        }
    </style>
</head>
<body>
    <h1>Tarjetas de estudio</h1>
    <hr>
    <form id="formTarjetas" method="post">
        <input type="hidden" name="tema" value="` + strconv.Itoa(int(tarjetas[0].IDTema)) + `">`)

	// Generar cada tarjeta
	if len(tarjetas) == 0 {
		html.WriteString(`<p>No hay tarjetas para este tema.</p>`)
	} else {
		for i, tarjeta := range tarjetas {
			visible := ""
			if i == 0 {
				visible = "visible"
			}

			html.WriteString(fmt.Sprintf(`
				<div class="tarjeta %s" id="tarjeta-%d">
					<div class="pregunta">%s</div>
					<div class="opcion">A. %s</div>
					<div class="opcion">B. %s</div>
					<div class="opcion">C. %s</div>
					<input type="hidden" name="pregunta%d" value="%s">
					<input type="hidden" name="respuestaCorrecta%d" value="%s">
					<input type="text" name="respuesta%d" class="campo-texto" placeholder="Escribe A, B o C" required>
					<button type="button" class="boton" onclick="siguiente(%d)">Siguiente</button>
				</div>
			`,
				visible,
				i+1,
				tarjeta.Pregunta,
				tarjeta.OpcionA,
				tarjeta.OpcionB,
				tarjeta.OpcionC,
				i+1,
				tarjeta.Pregunta,
				i+1,
				tarjeta.Respuesta,
				i+1,
				i+1))
		}
	}

	html.WriteString(`
    </form>
    <script>
        function siguiente(n) {
            const actual = document.getElementById('tarjeta-' + n);
            const siguiente = document.getElementById('tarjeta-' + (n + 1));
            
            // Validar que se haya ingresado una respuesta
            const respuestaInput = actual.querySelector('input[type="text"]');
            if (!respuestaInput.value.trim()) {
                alert('Por favor, ingresa tu respuesta antes de continuar.');
                return;
            }
            
            if (!siguiente) {
                // Última tarjeta - enviar formulario
                document.getElementById('formTarjetas').submit();
                return;
            }
            
            actual.classList.remove('visible');
            siguiente.classList.add('visible');
        }
    </script>
</body>
</html>`)

	return html.String()
}

// Función auxiliar para determinar qué tarjeta debe ser visible
func getVisibleClass(index int) string {
	if index == 0 {
		return "visible"
	}
	return ""
}

func verificarRespuestasHandler(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries, temaId int) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	var resultados []Resultado
	correctas := 0
	total := 0

	// Procesar cada respuesta
	for key, values := range r.Form {
		if strings.HasPrefix(key, "respuesta") && !strings.HasPrefix(key, "respuestaCorrecta") {
			idxStr := strings.TrimPrefix(key, "respuesta")
			_, err := strconv.Atoi(idxStr)
			if err != nil {
				continue
			}

			respuestaUsuario := strings.TrimSpace(strings.ToUpper(values[0]))
			respuestaCorrecta := strings.TrimSpace(strings.ToUpper(r.Form.Get("respuestaCorrecta" + idxStr)))
			pregunta := r.Form.Get("pregunta" + idxStr)

			esCorrecta := respuestaUsuario == respuestaCorrecta
			if esCorrecta {
				correctas++
			}
			total++

			resultados = append(resultados, Resultado{
				Pregunta:          pregunta,
				RespuestaUsuario:  values[0],
				RespuestaCorrecta: respuestaCorrecta,
				EsCorrecta:        esCorrecta,
			})
		}
	}

	// Generar HTML de resultados
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := generarHTMLResultados(r, resultados, correctas, total, temaId)
	w.Write([]byte(html))
}

// Función para generar HTML de resultados
func generarHTMLResultados(r *http.Request, resultados []Resultado, correctas int, total int, temaId int) string {
	var html strings.Builder

	porcentaje := 0
	if total > 0 {
		porcentaje = (correctas * 100) / total
	}

	html.WriteString(`<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Resultados</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        .resultado { border: 1px solid #ddd; padding: 15px; margin: 10px 0; border-radius: 5px; }
        .correcta { background-color: #d4edda; border-color: #c3e6cb; }
        .incorrecta { background-color: #f8d7da; border-color: #f5c6cb; }
        .estadisticas { background-color: #e9ecef; padding: 15px; border-radius: 5px; margin-bottom: 20px; text-align: center; }
        .pregunta { font-weight: bold; margin-bottom: 10px; }
        .respuesta { margin: 5px 0; }
        .correcta-text { color: #155724; font-weight: bold; }
        .incorrecta-text { color: #721c24; font-weight: bold; }
        .volver { display: inline-block; margin: 5px; padding: 10px 20px; background-color: #007bff; color: white; text-decoration: none; border-radius: 5px; }
    </style>
</head>
<body>
    <h1>Resultados del Cuestionario</h1>
    
    <div class="estadisticas">
        <h2>Estadísticas</h2>
        <p><strong>Puntuación:</strong> ` + strconv.Itoa(correctas) + `/` + strconv.Itoa(total) + `</p>
        <p><strong>Porcentaje:</strong> ` + strconv.Itoa(porcentaje) + `%</p>
    </div>
    
    <h2>Detalle de Respuestas:</h2>`)

	for i, resultado := range resultados {
		clase := "correcta"
		textoEstado := "✓ Correcta"
		if !resultado.EsCorrecta {
			clase = "incorrecta"
			textoEstado = "✗ Incorrecta"
		}

		html.WriteString(fmt.Sprintf(`
		<div class="resultado %s">
			<div class="pregunta">Pregunta %d: %s</div>
			<div class="respuesta"><strong>Tu respuesta:</strong> %s</div>
			<div class="respuesta"><strong>Respuesta correcta:</strong> %s</div>
			<div class="respuesta %s-text">%s</div>
		</div>
		`,
			clase,
			i+1,
			resultado.Pregunta,
			resultado.RespuestaUsuario,
			resultado.RespuestaCorrecta,
			clase,
			textoEstado))
	}

	html.WriteString(`
    <br>
    <a href="/tarjetas?tema=` + strconv.Itoa(temaId) + `" class="volver">Volver a intentar</a>
    <a href="/" class="volver">Volver al inicio</a>
</body>
</html>`)

	return html.String()
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
	w.Header().Set("Content-Type", "application/json") // <- Servidor dice: "Te voy a enviar JSON"
	w.WriteHeader(http.StatusCreated)                  // <- Servidor dice: "Operación exitosa, recurso creado (201)"

	// Enviar los datos como JSON válido
	if err := json.NewEncoder(w).Encode(newTarjeta); err != nil { // <- Servidor ENVÍA los datos al cliente
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
	//GET/tarjetas=1
	case http.MethodGet:
		getTarjetaByID(w, r, id, queries)
	//PUT/tarjetas=1
	case http.MethodPut:
		putTarjetaByID(w, r, id, queries)
	//DELETE/tarjetas=1
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
