BINARY=myapp

# Compilar
build:
	@echo ">= Building application..."
	@go build -o $(BINARY) .
# Levantar servicios con Docker Compose
up:
	@docker compose up -d
# Generar código con sqlc
sqlc: db/Schema/*.sql db/Queries/*.sql
	@sqlc generate
# Actualizar dependencias
tidy: go.mod go.sum
	@go mod tidy
# Ejecutar servidor
run: up tidy sqlc build
	@air
down:
	@docker compose down
	@rm -rf $(BINARY)
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

# Ejecutar tarjetas
listTarjetas:
	@curl -s http://localhost:8080/tarjetas | jq
listTarjetasByTema:
	@curl -s http://localhost:8080/tarjetas?tema=$(tema) | jq
createTarjeta:
	@curl -X POST http://localhost:8080/tarjetas \
	-H "Content-Type: application/json" \
	-d '{"pregunta":"$(pregunta)","respuesta":"$(respuesta)","opcion_a":"$(opcion_a)","opcion_b":"$(opcion_b)","opcion_c":"$(opcion_c)","id_tema":$(id_tema)}'
getTarjetaByID:
	@curl -X GET http://localhost:8080/tarjetas/$(id) | jq
putTarjetaByID:
	@curl -X PUT http://localhost:8080/tarjetas/$(id) \
	-H "Content-Type: application/json" \
	-d '{"pregunta":"$(pregunta)","respuesta":"$(respuesta)","opcion_a":"$(opcion_a)","opcion_b":"$(opcion_b)","opcion_c":"$(opcion_c)","id_tema":$(id_tema)}'
deleteTarjetaByID:
	@curl -X DELETE http://localhost:8080/tarjetas/$(id)

testTarjetas:
	@echo "=== Listar tarjetas inicial ==="
	@curl -s http://localhost:8080/tarjetas | jq
	@echo "=== Crear tarjeta 1 ==="
	@ID1=$$(curl -s -X POST http://localhost:8080/tarjetas \
	  -H "Content-Type: application/json" \
	  -d '{"pregunta":"¿Capital de Francia?","respuesta":"Paris","opcion_a":"Londres","opcion_b":"Paris","opcion_c":"Madrid","id_tema":1}' | jq -r '.id_tarjeta'); \
	echo "Tarjeta 1 ID: $$ID1"; \
	echo "=== Crear tarjeta 2 ==="; \
	ID2=$$(curl -s -X POST http://localhost:8080/tarjetas \
	  -H "Content-Type: application/json" \
	  -d '{"pregunta":"¿Capital de España?","respuesta":"Madrid","opcion_a":"Barcelona","opcion_b":"Valencia","opcion_c":"Madrid","id_tema":1}' | jq -r '.id_tarjeta'); \
	ID3=$$(curl -s -X POST http://localhost:8080/tarjetas \
      -H "Content-Type: application/json" \
      -d '{"pregunta":"¿Cuál es el río más largo del mundo?","respuesta":"Nilo","opcion_a":"Amazonas","opcion_b":"Nilo","opcion_c":"Misisipi","id_tema":3}' | jq -r '.id_tarjeta'); \
	ID4=$$(curl -s -X POST http://localhost:8080/tarjetas \
      -H "Content-Type: application/json" \
      -d '{"pregunta":"¿Qué planeta es conocido como el planeta rojo?","respuesta":"Marte","opcion_a":"Júpiter","opcion_b":"Venus","opcion_c":"Marte","id_tema":3}' | jq -r '.id_tarjeta'); \
	echo "Tarjeta 2 ID: $$ID2"; \
	echo "Tarjeta 3 ID: $$ID3"; \
	echo "Tarjeta 4 ID: $$ID4"; \
	echo "=== Listar tarjetas después de crear ==="; \
	curl -s http://localhost:8080/tarjetas | jq; \
	echo "=== Listar tarjetas por tema (id_tema=1) ==="; \
	curl -s http://localhost:8080/tarjetas?tema=3 | jq; \
	echo "=== Obtener tarjeta 1 ==="; \
	curl -s http://localhost:8080/tarjetas/$$ID1 | jq; \
	echo "=== Actualizar tarjeta 1 ==="; \
	curl -s -X PUT http://localhost:8080/tarjetas/$$ID1 \
	  -H "Content-Type: application/json" \
	  -d '{"pregunta":"¿Capital de Francia actualizada?","respuesta":"París","opcion_a":"Londres","opcion_b":"Berlín","opcion_c":"Madrid","id_tema":1}' | jq; \
	echo "=== Listar tarjetas después de actualizar ==="; \
	curl -s http://localhost:8080/tarjetas | jq; \
	echo "=== Eliminar tarjeta 1 ==="; \
	curl -s -X DELETE http://localhost:8080/tarjetas/$$ID1; \
	echo "=== Eliminar tarjeta 2 ==="; \
	curl -s -X DELETE http://localhost:8080/tarjetas/$$ID2; \
	echo "=== Eliminar tarjeta 3 ==="; \
	curl -s -X DELETE http://localhost:8080/tarjetas/$$ID3; \
	echo "=== Eliminar tarjeta 4 ==="; \
	curl -s -X DELETE http://localhost:8080/tarjetas/$$ID4; \
	echo "=== Listar tarjetas final ==="; \
	curl -s http://localhost:8080/tarjetas | jq;

# Ejecutar temas
listTemas:
	@curl -s http://localhost:8080/temas | jq
createTema:
	@curl -X POST http://localhost:8080/temas \
	-d '{"nombre_tema":"$(nombre)"}'
getTemaByID:
	@curl -X GET http://localhost:8080/temas/$(id) | jq
putTemaByID:
	@curl -X PUT http://localhost:8080/temas/$(id) \
	-H "Content-Type: application/json" \
	-d '{"nombre_tema":"$(nombre)"}'
deleteTemaByID:
	@curl -X DELETE http://localhost:8080/temas/$(id)
testTemas:
	@echo "=== Listar temas inicial ==="
	@curl -s http://localhost:8080/temas | jq
	@echo "=== Crear tema 1 ==="
	@ID1=$$(curl -s -X POST http://localhost:8080/temas \
	  -H "Content-Type: application/json" \
	  -d '{"nombre_tema":"Geografía"}' | jq -r '.id_tema'); \
	echo "Tema 1 ID: $$ID1"; \
	echo "=== Crear tema 2 ==="; \
	ID2=$$(curl -s -X POST http://localhost:8080/temas \
	  -H "Content-Type: application/json" \
	  -d '{"nombre_tema":"Ciencia"}' | jq -r '.id_tema'); \
	echo "Tema 2 ID: $$ID2"; \
	echo "=== Listar temas después de crear ==="; \
	curl -s http://localhost:8080/temas | jq; \
	echo "=== Obtener tema 1 ==="; \
	curl -s http://localhost:8080/temas/$$ID1 | jq; \
	echo "=== Actualizar tema 1 ==="; \
	curl -s -X PUT http://localhost:8080/temas/$$ID1 \
	  -H "Content-Type: application/json" \
	  -d '{"nombre_tema":"Geografía Actualizada"}' | jq; \
	echo "=== Listar temas después de actualizar ==="; \
	curl -s http://localhost:8080/temas | jq; \
	echo "=== Eliminar tema 1 ==="; \
	curl -s -X DELETE http://localhost:8080/temas/$$ID1; \
	echo "=== Eliminar tema 2 ==="; \
	curl -s -X DELETE http://localhost:8080/temas/$$ID2; \
	echo "=== Listar temas final ==="; \
	curl -s http://localhost:8080/temas | jq
# Ejecutar flujo completo
allTests: testUsers testTarjetas testTemas