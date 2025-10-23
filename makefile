BINARY=myapp

# Compilar
build:
	@echo ">= Building application..."
	@go build -o $(BINARY) .

# Ejecutar servidor
run: build
	@air
# Ejecutar getProducts
listUsuarios:
	curl -X GET http://localhost:8080/users
# Crear Usuario nombre="Juan Perez" email="mail@example.com" contrasena="securepassword" por ejemplo
createUsuario:
	curl -X POST http://localhost:8080/users \
	-H "Content-Type: application/json" \
	-d '{"nombre_usuario":"$(nombre)","email":"$(email)","contrasena":"$(contrasena)"}'




add-products:
	# Agregar productos de ejemplo
	curl -X POST http://localhost:8080/products \
	-H "Content-Type: application/json" \
	-d '{"name":"LaptopNueva",
	"description":"Laptop 16GB RAM",
	"price":3500,
	"quantity":10}'

	curl -X POST http://localhost:8080/products \
	-H "Content-Type: application/json" \
	-d '{"name":"MouseNuevo","description":"Mouse inalámbrico","price":50,"quantity":100}'

# Test endpoints con distintos IDs
test-ids:
	# GET por ID
	curl -X GET http://localhost:8080/products/1

	curl -X GET http://localhost:8080/products/2

	# PUT por ID (actualización de ejemplo)
	curl -X PUT http://localhost:8080/products/1 \
	-H "Content-Type: application/json" \
	-d '{"name":"Laptop Updated","description":"Laptop 16GB RAM","price":3600,"quantity":12}'
	

	curl -X PUT http://localhost:8080/products/2 \
	-H "Content-Type: application/json" \
	-d '{"name":"Mouse Updated","description":"Mouse inalámbrico","price":55,"quantity":90}'
	

	# DELETE por ID
	curl -X DELETE http://localhost:8080/products/3
	
	curl -X DELETE http://localhost:8080/products/4
	

# Ejecutar flujo completo
all: run get-products add-productstest-ids get-products stop clean

