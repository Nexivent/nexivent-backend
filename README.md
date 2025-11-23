# nexivent-Backend
Backend en Go para la tiketera Nexivent

# Requisitos previos
- Go
- Docker Desktop
  
# Configuración del entorno

1. Crear un archivo .env
2. Levantar la base de datos con Docker
        docker compose up -d
3. Verificar si se levantó la base de datos
        docker ps

# Ver los logs de la base de datos
Usar el siguiente comando
    docker logs -f nexivent-db

# Comando para ejecutar la creación de tablas
Get-Content -Raw .\migrations\001_create_tables.sql |   
>>   docker exec -i nexivent-db psql -U postgres -d nexivent

# Comando para verificar las tablas
1. docker exec -it nexivent-db psql -U postgres -d nexivent 
2. \dt

# Despliegue en Railway (backend + frontend)

1. Instala la CLI de Railway y autentícate:
   ```bash
   npm i -g @railway/cli
   railway login
   ```
2. Backend (este repo):
   - Ejecuta `railway init --service backend` para crear o vincular el servicio.
   - Añade Postgres con `railway add postgres`. El backend ya lee `DATABASE_URL` y las variables `PG*` que expone Railway y usa `PORT` automáticamente.
   - Variables recomendadas: `ENABLE_SWAGGER=false`, `CORS_ALLOWED_ORIGINS=https://tu-frontend.railway.app` (puedes añadir varias separadas por comas), `AWS_*` si usas S3, `MAIL_*`, `FACTILIZA_TOKEN`.
   - Despliega con `railway up --service backend`. Railway usará el `Dockerfile` y `railway.json` (healthcheck en `/health-check/`).
3. Frontend React:
   - En el repo del frontend crea otro servicio en el mismo proyecto Railway.
   - Define `VITE_API_URL` (o la variable equivalente) apuntando al dominio del backend.
   - Ejemplo mínimo de `railway.json` para un frontend Vite:
     ```json
     {
       "$schema": "https://railway.com/railway.schema.json",
       "build": { "builder": "NIXPACKS" },
       "deploy": { "startCommand": "npm run preview -- --host 0.0.0.0 --port $PORT" }
     }
     ```
   - Despliega desde el repo del frontend con `railway up --service frontend`.
4. Dominios/CORS: añade el dominio público del frontend en `CORS_ALLOWED_ORIGINS` para que las peticiones en producción funcionen.
