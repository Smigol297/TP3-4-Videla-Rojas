BINARY=myapp

# Compilar
build:
	@echo ">= Building application..."
	@go build -o $(BINARY) .

# Ejecutar servidor
run: build
	@air
# Ejecutar getUsers
listUsuarios:
	@curl -s http://localhost:8080/users | jq
# Crear Usuario nombre="Juan Perez" email="mail@example.com" contrasena="securepassword" por ejemplo
createUsuario:
	@curl -X POST http://localhost:8080/users \
	-H "Content-Type: application/json" \
	-d '{"nombre_usuario":"$(nombre)","email":"$(email)","contrasena":"$(contrasena)"}'
getUserByID:
	@curl -X GET http://localhost:8080/users/$(id) | jq
putUserByID:
	@curl -X PUT http://localhost:8080/users/$(id) \
	-H "Content-Type: application/json" \
	-d '{"nombre_usuario":"$(nombre)","email":"$(email)","contrasena":"$(contrasena)"}'
deleteUserByID:
	@curl -X DELETE http://localhost:8080/users/$(id)

testUsers:
	@echo "=== Listar usuarios inicial ==="
	@curl -s http://localhost:8080/users | jq
	@echo "=== Crear usuario 1 ==="
	@ID1=$$(curl -s -X POST http://localhost:8080/users \
	  -H "Content-Type: application/json" \
	  -d '{"nombre_usuario":"Juan Perez","email":"juan@example.com","contrasena":"securepassword"}' | jq -r '.id_usuario'); \
	echo "Usuario 1 ID: $$ID1"; \
	echo "=== Crear usuario 2 ==="; \
	ID2=$$(curl -s -X POST http://localhost:8080/users \
	  -H "Content-Type: application/json" \
	  -d '{"nombre_usuario":"Maria Lopez","email":"maria@example.com","contrasena":"securepassword"}' | jq -r '.id_usuario'); \
	echo "Usuario 2 ID: $$ID2"; \
	echo "=== Listar usuarios después de crear ==="; \
	curl -s http://localhost:8080/users | jq; \
	echo "=== Obtener usuario 1 ==="; \
	curl -s http://localhost:8080/users/$$ID1 | jq; \
	echo "=== Actualizar usuario 1 ==="; \
	curl -s -X PUT http://localhost:8080/users/$$ID1 \
	  -H "Content-Type: application/json" \
	  -d '{"nombre_usuario":"Juan Actualizado","email":"juan@example.com","contrasena":"securepassword"}' | jq; \
	echo "=== Listar usuarios después de actualizar ==="; \
	curl -s http://localhost:8080/users | jq; \
	echo "=== Eliminar usuario 1 ==="; \
	curl -s -X DELETE http://localhost:8080/users/$$ID1; \
	echo "=== Eliminar usuario 2 ==="; \
	curl -s -X DELETE http://localhost:8080/users/$$ID2; \
	echo "=== Listar usuarios final ==="; \
	curl -s http://localhost:8080/users | jq


# Ejecutar flujo completo
all: run testUsers


