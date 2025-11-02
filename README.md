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

  

