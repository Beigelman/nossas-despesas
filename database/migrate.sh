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

# Run migrations using go-migrate in the specified direction and path
migrate -path "$MIGRATION_PATH" -database "${DB_CONNECTION_STRING}" "$MIGRATION_DIRECTION"

# Check if migrate succeeded
# shellcheck disable=SC2181
if [ $? -eq 0 ]; then
    echo "Migrations applied successfully"
else
    echo "Failed to apply migrations."
    exit 1
fi
