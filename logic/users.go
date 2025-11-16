package logic

import (
	"fmt"
	sqlc "tp3/db"
)

func ValidateCreateUser(p sqlc.CreateUsuarioParams) error {
	if p.NombreUsuario == "" {
		return fmt.Errorf("el nombre del usuario no puede estar vacío")
	}
	if p.Email == "" {
		return fmt.Errorf("el email del usuario no puede estar vacío")
	}
	if p.Contrasena == "" {
		return fmt.Errorf("la contraseña del usuario no puede estar vacía")
	}
	return nil
}

func ValidateUpdateUser(p sqlc.UpdateUsuarioParams) error {
	if p.IDUsuario <= 0 {
		return fmt.Errorf("ID de usuario %d inválido", p.IDUsuario)
	}
	if p.NombreUsuario == "" {
		return fmt.Errorf("el nombre del usuario no puede estar vacío")
	}
	if p.Email == "" {
		return fmt.Errorf("el email del usuario no puede estar vacío")
	}
	if p.Contrasena == "" {
		return fmt.Errorf("la contraseña del usuario no puede estar vacía")
	}
	return nil
}
