# Extensión UUID postgres
para ejecutar postgres necesita tener la extensión *uuid-ossp*, existen varias formar de instalar con docker-compose. Pero la que me funcionó fue:
```bash
docker exec -it srvpostgres bash
```
```bash
psql -U ecotrueque -d ecotrueque -f /docker-entrypoint-initdb.d/init.sql
```



[Readme Sintaxis](https://docs.github.com/es/get-started/writing-on-github/getting-started-with-writing-and-formatting-on-github/basic-writing-and-formatting-syntax)