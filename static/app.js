// Espera a que todo el contenido del DOM (HTML) esté cargado
document.addEventListener('DOMContentLoaded', () => {

    // --- 1. URLs DE API ---
    const TEMA_API = '/temas';
    const USUARIO_API = '/users'; // Tu handler de Go usa '/users'
    const TARJETA_API = '/tarjetas';

    // --- 2. SELECTORES DE ELEMENTOS ---

    // Selectores de Temas
    const temaForm = document.getElementById('tema-form');
    const temaNombreInput = document.getElementById('tema-nombre');
    const temaList = document.getElementById('tema-list');

    // Selectores de Usuarios
    const usuarioForm = document.getElementById('usuario-form');
    const usuarioNombreInput = document.getElementById('usuario-nombre');
    const usuarioEmailInput = document.getElementById('usuario-email');
    const usuarioContrasenaInput = document.getElementById('usuario-contrasena');
    const usuarioList = document.getElementById('usuario-list');

    // Selectores de Tarjetas
    /*const tarjetaForm = document.getElementById('tarjeta-form');
    const tarjetaPreguntaInput = document.getElementById('tarjeta-pregunta');
    const tarjetaRespuestaInput = document.getElementById('tarjeta-respuesta');
    const tarjetaOpcionAInput = document.getElementById('tarjeta-opcion-a');
    const tarjetaOpcionBInput = document.getElementById('tarjeta-opcion-b');
    const tarjetaOpcionCInput = document.getElementById('tarjeta-opcion-c');
    const tarjetaIdTemaInput = document.getElementById('tarjeta-id-tema');
    const tarjetaList = document.getElementById('tarjeta-list');*/


    // --- 3. SECCIÓN: TEMAS ---

    // GET /temas
    async function fetchTemas() {
        try {
            //await: Pausa la función fetchTemas hasta que el servidor responda
            //fetch(): Es la función del navegador para hacer peticiones de red (por defecto una petición GET).
            //response guardará la respuesta HTTP inicial (no los datos en sí, sino el estado de la conexión, cabeceras, etc.).
            const response = await fetch(TEMA_API);
            if (!response.ok) throw new Error('Error al cargar temas');
            //Este método lee el cuerpo de la respuesta y lo transforma (parsea) de texto JSON a un objeto o array de JavaScript
            const temas = await response.json();
            //evito duplicacion de temas 
            temaList.innerHTML = ''; // Limpiar lista
            if (temas && temas.length > 0) {
                temas.forEach(tema => {
                    const li = document.createElement('li');
                    li.textContent = `(ID: ${tema.id_tema}) - ${tema.nombre_tema} `;
                    
                    //cada elemento tiene la opcion de eliminarse.
                    const deleteButton = document.createElement('button');
                    deleteButton.textContent = 'Eliminar';
                    deleteButton.onclick = () => deleteTema(tema.id_tema);
                    
                    li.appendChild(deleteButton);
                    temaList.appendChild(li);
                });
            } else {
                temaList.innerHTML = '<li>No hay temas creados.</li>';
            }
        } catch (error) {
            console.error('Error en fetchTemas:', error);
            temaList.innerHTML = '<li>Error al cargar la lista.</li>';
        }
    }

    // POST /temas
    //Cuando el navegador dispara un evento (como el envío de un formulario), pasa un objeto event a la función que lo maneja.
    async function handleTemaSubmit(event) {
        //preventDefault() evita comportamiento de recarga de página al enviar el formulario.
        event.preventDefault();
        const data = {
            nombre_tema: temaNombreInput.value.trim()
        };
        if (!data.nombre_tema) {
            alert('El nombre del tema es obligatorio');
            return;
        }

        try {
            //URL base (TEMA_API)
            //crear un nuevo recurso en el servidor.
            const response = await fetch(TEMA_API, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                //data es tu objeto JavaScript: { nombre_tema: "Mi Tema" }.
                //convierte en un texto (string) con formato JSON
                body: JSON.stringify(data)
            });
            if (!response.ok) throw new Error(await response.text());
            //.reset() borra automáticamente todos los campos del formulario, dejándolo listo para que el usuario pueda añadir otro tema.
            temaForm.reset();
            await fetchTemas(); // Refrescar lista
        } catch (error) {
            console.error('Error al crear tema:', error);
            alert(`Error al crear tema: ${error.message}`);
        }
    }

    // DELETE /temas/{id}
    async function deleteTema(id) {
        //¿Eliminar tema con ese ID?". Si le doy a "Cancelar", la función se detiene (return;).
        if (!confirm(`¿Eliminar tema con ID ${id}?`)) return;
        try {
            //Si le doy a aceptar, se hace la petición DELETE al servidor.
            //TEMA_API/id construye la URL.
            const response = await fetch(`${TEMA_API}/${id}`, { method: 'DELETE' });
            if (!response.ok) throw new Error(await response.text());
            await fetchTemas(); // Refrescar lista (ACTUALIZA)
        } catch (error) {
            console.error('Error al eliminar tema:', error);
            alert(`Error al eliminar tema: ${error.message}`);
        }
    }

    // --- 4. SECCIÓN: USUARIOS ---

    // GET /users
    async function fetchUsuarios() {
        try {
            const response = await fetch(USUARIO_API);
            if (!response.ok) throw new Error('Error al cargar usuarios');
            const usuarios = await response.json();

            usuarioList.innerHTML = ''; // Limpiar lista
            if (usuarios && usuarios.length > 0) {
                usuarios.forEach(user => {
                    const li = document.createElement('li');
                    // Mostramos solo datos no sensibles
                    li.textContent = `(ID: ${user.id_usuario}) - ${user.nombre_usuario} (${user.email}) `;
        
                    const deleteButton = document.createElement('button');
                    deleteButton.textContent = 'Eliminar';
                    deleteButton.onclick = () => deleteUsuario(user.id_usuario);
                
                    li.appendChild(deleteButton);
                    usuarioList.appendChild(li);
                });
            } else {
                usuarioList.innerHTML = '<li>No hay usuarios creados.</li>';
            }
        } catch (error) {
            console.error('Error en fetchUsuarios:', error);
            usuarioList.innerHTML = '<li>Error al cargar la lista.</li>';
        }
    }

    // POST /users
    async function handleUsuarioSubmit(event) {
        event.preventDefault();
        const data = {
            nombre_usuario: usuarioNombreInput.value.trim(),
            email: usuarioEmailInput.value.trim(),
            contrasena: usuarioContrasenaInput.value.trim()
        };

        if (!data.nombre_usuario || !data.email || !data.contrasena) {
            alert('Todos los campos de usuario son obligatorios');
        return;
        }

        try {
            const response = await fetch(USUARIO_API, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
            });
        
            if (!response.ok) throw new Error(await response.text());

            usuarioForm.reset();
            await fetchUsuarios(); // Refrescar lista
        } catch (error) {
            console.error('Error al crear usuario:', error);
            alert(`Error al crear usuario: ${error.message}`);
        }
    }

    // DELETE /users/{id}
    async function deleteUsuario(id) {
        if (!confirm(`¿Eliminar usuario con ID ${id}?`)) return;
        try {
            const response = await fetch(`${USUARIO_API}/${id}`, { method: 'DELETE' });
            if (!response.ok) throw new Error(await response.text());
            await fetchUsuarios(); // Refrescar lista
        } catch (error) {
            console.error('Error al eliminar usuario:', error);
            alert(`Error al eliminar usuario: ${error.message}`);
        }
    }


    // --- 6. INICIALIZACIÓN ---

    // Añadir listeners a los formularios
    temaForm.addEventListener('submit', handleTemaSubmit);
    usuarioForm.addEventListener('submit', handleUsuarioSubmit);
    //tarjetaForm.addEventListener('submit', handleTarjetaSubmit);
    
    // Cargar los datos iniciales de todas las entidades
    fetchTemas();
    fetchUsuarios();
    //fetchTarjetas();
});