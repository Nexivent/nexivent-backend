#!/usr/bin/env bash
set -euo pipefail

if ! [[ -x "$(command -v psql)" ]]; then
	echo "Error: psql is not installed." >&2
	exit 1
fi

if ! [[ -x "$(command -v rg)" ]]; then
	echo "Error: ripgrep is not installed." >&2
	exit 1
fi

docker compose up -d

# Waiting for initialization
sleep 10

until (
	set -a 
	source <(rg -v -e '^#|DB_PASSWORD|DATABASE_URL|APP' -e '^[[:space:]]*$' .env)
	set +a
	psql -h $DB_HOST -U $DB_USER -p $DB_PORT -d $DB_NAME -w -c '\q'
	)
do
	echo "Postgres is still unavailable - sleeping" >&2
	sleep 1
done

(
	set -a  
	source <(rg '^DB_PORT' .env)
	set +a
	echo "Postgres is up and running on port $DB_PORT - running migrations now!" >&2
)

