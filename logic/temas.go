package logic

import (
	"fmt"
	sqlc "tp3/db" // Importa tus tipos generados por sqlc
)

// --- ESTE ES EL CÓDIGO "LIMPIO" ---
// No hay http.ResponseWriter. Solo lógica de validación pura.

func ValidateCreateTema(nombre string) error {
	if nombre == "" {
		return fmt.Errorf("el nombre del tema no puede estar vacío")
	}
	return nil
}

func ValidateUpdateTema(p sqlc.UpdateTemaParams) error {
	if p.IDTema <= 0 {
		return fmt.Errorf("ID de tema %d inválido", p.IDTema)
	}
	if p.NombreTema == "" {
		return fmt.Errorf("el nombre del tema no puede estar vacío")
	}

	return nil
}
