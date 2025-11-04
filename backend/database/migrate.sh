#!/bin/bash

# Check for correct number of arguments
if [ $# -ne 2 ]; then
    echo "Usage: $0 <up|down> <migration_path>"
    exit 1
fi

# Assign arguments to variables
MIGRATION_DIRECTION=$1
MIGRATION_PATH=$2

# Validate migration direction
if [ "$MIGRATION_DIRECTION" != "up" ] && [ "$MIGRATION_DIRECTION" != "down" ]; then
    echo "Invalid migration direction: $MIGRATION_DIRECTION. Please specify 'up' or 'down'."
    exit 1
fi

# Validate migration path
if [ ! -d "$MIGRATION_PATH" ]; then
    echo "Invalid migration path: $MIGRATION_PATH does not exist or is not a directory."
    exit 1
fi

# Check if DB_CONNECTION_STRING is set
if [ -z "${DB_CONNECTION_STRING}" ]; then
    echo "DB_CONNECTION_STRING is not set. Attempting to construct from components..."

    # Check if required components are set
    if [ -z "${DB_HOST}" ] || [ -z "${DB_NAME}" ] || [ -z "${DB_PASSWORD}" ] || [ -z "${DB_PORT}" ] || [ -z "${DB_USER}" ]; then
        echo "One or more required environment variables (DB_HOST, DB_NAME, DB_PASSWORD, DB_PORT, DB_USER) are missing."
        export DB_CONNECTION_STRING="postgres://root:root@localhost:5432/app?sslmode=disable"
        echo "Using default DB_CONNECTION_STRING: ${DB_CONNECTION_STRING}"
    else
        # Construct DB_CONNECTION_STRING
        export DB_CONNECTION_STRING="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}"
        echo "Constructed DB_CONNECTION_STRING: ${DB_CONNECTION_STRING}"
    fi
else
    echo "Using existing DB_CONNECTION_STRING: ${DB_CONNECTION_STRING}"
fi

# Adiciona parâmetros para resolver problemas de prepared statement no PostgreSQL da nuvem
if [[ "${DB_CONNECTION_STRING}" == *"?"* ]]; then
    # Se já tem parâmetros, adiciona os novos com &
    ENHANCED_DB_CONNECTION_STRING="${DB_CONNECTION_STRING}&x-migrations-table=schema_migrations&x-migrations-table-quoted=false&statement_timeout=0&lock_timeout=0"
else
    # Se não tem parâmetros, adiciona com ?
    ENHANCED_DB_CONNECTION_STRING="${DB_CONNECTION_STRING}?x-migrations-table=schema_migrations&x-migrations-table-quoted=false&statement_timeout=0&lock_timeout=0"
fi

echo "Enhanced connection string (masked): ${ENHANCED_DB_CONNECTION_STRING//:*@/:***@}"

# Configura retry para resolver problemas temporários de conexão
RETRY_ATTEMPTS=3
RETRY_COUNT=0

while [ $RETRY_COUNT -lt $RETRY_ATTEMPTS ]; do
    echo "Migration attempt $((RETRY_COUNT + 1)) of $RETRY_ATTEMPTS..."
    
    # Run migrations using go-migrate in the specified direction and path
    if migrate -path "$MIGRATION_PATH" -database "${ENHANCED_DB_CONNECTION_STRING}" "$MIGRATION_DIRECTION"; then
        echo "Migrations applied successfully"
        exit 0
    else
        RETRY_COUNT=$((RETRY_COUNT + 1))
        if [ $RETRY_COUNT -lt $RETRY_ATTEMPTS ]; then
            echo "Migration failed. Retrying in 5 seconds... (attempt $RETRY_COUNT/$RETRY_ATTEMPTS)"
            sleep 5
        fi
    fi
done

echo "Failed to apply migrations after $RETRY_ATTEMPTS attempts."
exit 1
