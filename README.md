# 📦 API REST con Go y SQLite

Este proyecto es una API REST escrita en Go utilizando el enrutador [`gorilla/mux`](https://github.com/gorilla/mux), con SQLite como base de datos. Se encuentra totalmente contenerizada con Docker.

## 🛠️ Tecnologías

- [Go 1.24](https://go.dev/)
- [gorilla/mux](https://github.com/gorilla/mux)
- [SQLite](https://www.sqlite.org/index.html)
- [Docker](https://www.docker.com/)
- [Swagger (swaggo)](https://github.com/swaggo/swag)

## 📁 Estructura del Proyecto

├── main.go
├── handlers.go
├── models.go
├── database.go
├── series.db
├── go.mod
├── go.sum
├── Dockerfile

## 🐳 Docker

### 🧱 Construcción de la imagen

```bash
docker build -t api .
```

### ▶️ Ejecución del contenedor

```bash
docker run -d -p 8080:8080 api
```

### 🛑 Detener contenedor

```bash
docker ps         # Ver ID o nombre del contenedor
docker stop <id o nombre>
```

### 🗃️ Acceso a la base de datos SQLite
Si deseas ingresar al contenedor e inspeccionar o modificar la base de datos SQLite:
```bash
docker ps
docker exec -it <id o nombre> sh
sqlite3 series.db
```

## 📄 Documentación con Swagger

La API incluye documentación interactiva generada con Swagger.

Una vez que el contenedor esté corriendo, puedes acceder a la documentación en tu navegador visitando:

👉 [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## 📬 Colección de Postman

Puedes probar los endpoints de la API utilizando la siguiente colección de Postman:

👉 [Abrir colección en Postman](https://angelesquit.postman.co/workspace/Angel-Esquit's-Workspace~6b3c6df4-f2c9-4a66-8550-c81b540553aa/collection/44134349-b9ee865e-c662-4a07-89b7-e0de63644f1e?action=share&creator=44134349)

Recuerda tener el contenedor de la API corriendo antes de probar los endpoints.
