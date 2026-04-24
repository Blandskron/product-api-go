# Product Management API (Go)


Una API REST robusta y de nivel profesional construida en Go para la gestión de un inventario de productos. La API implementa un CRUD completo para productos y endpoints transaccionales para registrar ventas y compras, asegurando la integridad de los datos y el stock.

Este proyecto sirve como un excelente ejemplo de cómo aplicar principios de **Arquitectura Limpia (Clean Architecture)**, inyección de dependencias, manejo de errores centralizado y otras mejores prácticas en un entorno de Go.

## ✨ Características Principales

- **CRUD Completo para Productos**: Endpoints para crear, leer, actualizar y eliminar productos.
- **Gestión de Inventario Transaccional**:
  - **Endpoint de Ventas**: Descuenta el stock de un producto de forma atómica.
  - **Endpoint de Compras**: Aumenta el stock de un producto existente.
- **Arquitectura Limpia**: Código modular y desacoplado, separado en capas (`domain`, `repository`, `usecase`, `handler`).
- **Validación de Entrada Robusta**: Utiliza `go-playground/validator` para validar los datos de entrada a nivel de la lógica de negocio.
- **Manejo de Errores Centralizado**: Una única función `handleError` traduce los errores internos (de base de datos, validación, negocio) en respuestas HTTP claras y consistentes.
- **Inyección de Dependencias**: Las dependencias se inicializan en un único lugar y se inyectan en las capas que las necesitan, facilitando las pruebas y la mantenibilidad.
- **Documentación de API con Swagger**: Documentación generada automáticamente y una UI interactiva para probar los endpoints.
- **Configuración Externa**: Carga la configuración de la base de datos desde un archivo `.env` para una fácil adaptación entre entornos.
- **Apagado Controlado (Graceful Shutdown)**: El servidor finaliza las peticiones en curso antes de apagarse, previniendo la corrupción de datos.

## 🏛️ Arquitectura

El proyecto sigue los principios de la Arquitectura Limpia para separar las responsabilidades y asegurar un bajo acoplamiento entre componentes.

- **`domain`**: Contiene las entidades de negocio principales (`Product`, `Sale`, `Purchase`) y sus reglas de validación. Es el núcleo de la aplicación.
- **`repository`**: Define las interfaces para la capa de persistencia de datos. Su propósito es abstraer el acceso a la base de datos.
- **`usecase`**: Contiene la lógica de negocio específica de la aplicación. Orquesta el flujo de datos entre los repositorios y aplica las reglas de negocio (ej. validaciones, transacciones).
- **`handler`**: Es la capa de presentación (API). Se encarga de recibir las peticiones HTTP, decodificar los datos, llamar a los casos de uso correspondientes y devolver las respuestas HTTP.
- **`infrastructure`**: Contiene las implementaciones concretas de las interfaces, como el router de Gin, la conexión a la base de datos (GORM), y la implementación de los repositorios.

## 🚀 Cómo Empezar

Sigue estos pasos para levantar el proyecto en tu entorno local.

### Prerrequisitos

- **Go** (versión 1.21 o superior)
- **PostgreSQL**: Una instancia de base de datos en ejecución.
- **Git**

### Pasos de Instalación

1.  **Clona el repositorio:**
    ```bash
    git clone <URL_DEL_REPOSITORIO>
    cd product-api-go
    ```

2.  **Configura la Base de Datos:**
    Asegúrate de que tu instancia de PostgreSQL esté corriendo y crea una base de datos. Por ejemplo:
    ```sql
    CREATE DATABASE productos_db;
    ```

3.  **Crea el archivo de configuración `.env`:**
    En la raíz del proyecto, crea un archivo llamado `.env` y añade la cadena de conexión a tu base de datos.

    ```env
    # d:\SUSTANTIVA\ProductosApi\product-api-go\.env
    DATABASE_DSN="host=localhost user=postgres password=admin1234 dbname=productos_db port=5432 sslmode=disable TimeZone=UTC"
    ```
    *Ajusta los valores de `user`, `password`, `dbname`, etc., según tu configuración.*

4.  **Instala las dependencias:**
    Go Modules se encargará de descargar todas las librerías necesarias.
    ```bash
    go mod tidy
    ```

5.  **Ejecuta la aplicación:**
    ```bash
    go run cmd/api/main.go
    ```
    El servidor se iniciará en `http://localhost:8080`. Las tablas de la base de datos se crearán automáticamente la primera vez que se ejecute.

## 📚 Documentación de la API (Swagger)

Una vez que el servidor esté en funcionamiento, puedes acceder a la documentación interactiva de la API en la siguiente URL:

**http://localhost:8080/swagger/index.html**

Desde esta interfaz, puedes ver todos los endpoints, sus parámetros, y probarlos directamente.

#### Regenerar la Documentación

Si realizas cambios en los comentarios de Swagger (`// @Summary`, etc.), necesitas regenerar los archivos de documentación. Ejecuta el siguiente comando en la raíz del proyecto:
```bash
swag init -g cmd/api/main.go
```

## 🧪 Pruebas

El proyecto incluye pruebas unitarias para los casos de uso, utilizando mocks para aislar la lógica de negocio de la base de datos. Para ejecutar todas las pruebas:

```bash
go test ./...
```

## Endpoints de la API

| Método | Ruta                      | Descripción                               |
|--------|---------------------------|-------------------------------------------|
| `GET`  | `/api/products`           | Obtiene una lista de todos los productos. |
| `GET`  | `/api/products/{id}`      | Obtiene un producto por su ID.            |
| `POST` | `/api/products`           | Crea un nuevo producto.                   |
| `PUT`  | `/api/products/{id}`      | Actualiza un producto existente.          |
| `DELETE`| `/api/products/{id}`      | Elimina un producto por su ID.            |
| `POST` | `/api/sales`              | Procesa una venta y descuenta el stock.   |
| `POST` | `/api/purchases`          | Procesa una compra y aumenta el stock.    |
