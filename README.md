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