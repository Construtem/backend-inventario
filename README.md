# Backend - Inventario

Este microservicio gestiona el inventario y despacho en Construtem.

## 🛠️ Tecnologías
- Go (Lenguaje)
- Gin (Framework)
- GORM (Libreria de Go para interactuar con bases de datos)
- Postgres

## 🚀 Funcionalidades
- Control de stock por bodega y sucursal.
- Registro de movimientos (entrada/salida).
- Gestión de productos.

## Requisitos

- Docker Desktop instalado
- Git instalado 

  ## Instalación (entorno de desarrollo)

1. Clonar el repositorio en el directorio deseado:

*Desde la terminal debe situarse en el directorio que desee clonar repo (ej: "C:\Users\Admin\Desktop") y ejecutar siguiente comando*

<details>

<summary>**¿Cómo situarse en C:\Users\Admin\Desktop?**</summary>

1. Abrir terminal (Ya sea powershell, cmd, git bash, etc)
2. Te encontrarás situado en C:\Users\Admin o algo así
3. Debes ejecutar el comando
```bash
cd .\Desktop\
```
*Cualquier consulta escribirme a wsp +56979828311*
</details>

```bash
git clone https://github.com/Construtem/backend-inventario
cd backend-inventario
```
2. Correr aplicación desde directorio creado (ej "C:\Users\Admin\Desktop\backend-inventario"),
ejecutando el siguiente comando:
```bash
docker compose up
```
*Luego de ejecutar este comando, su app se encontrará corriendo en el puerto 8080 en "http://localhost:8080"*


## Contribución

1. Crea una rama para tu funcionalidad/tarea:

```bash
git switch -c feature/<nombre-funcionalidad>
```

2. Realiza cambios y haz commit:

```bash
git add <archivos-cambiados>
git commit -m "<descripcion pequeña del cambio>"
```

3. Pushea tus cambios de la rama:

```bash
git push origin feature/<nombre-funcionalidad> 
```

4. Crea un Pull Request (PR) a la rama ´develop´ desde GitHub para que sea revisado por otro desarrollador
