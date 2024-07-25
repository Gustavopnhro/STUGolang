#!/bin/bash

LOG_FILE="migration.log"

function log_message {
    local message="$1"
    echo "$(date +'%Y-%m-%d %H:%M:%S') - $message"
    echo "$(date +'%Y-%m-%d %H:%M:%S') - $message" >> "$LOG_FILE"
}

function check_variable {
    if [[ -z "${!1}" ]]; then
        log_message "Error: $1 is not set or is empty. Check your .env file."
        exit 1
    fi
}

set -o allexport
source .env
set +o allexport

check_variable "CONTAINER_NAME"
check_variable "DB_USERNAME"
check_variable "DB_PASSWORD"
check_variable "DB_DATABASE"

# Assign variables
CONTAINER_NAME=$CONTAINER_NAME
DB_USERNAME=$DB_USERNAME
DB_PASSWORD=$DB_PASSWORD
DB_DATABASE=$DB_DATABASE

function check_error {
    if [ $? -ne 0 ]; then
        log_message "Error occurred. Exiting..."
        exit 1
    fi
}

function execute_sql_file {
    local sql_file="$1"
    log_message "Executing SQL file $sql_file inside $CONTAINER_NAME..."
    docker exec -i $CONTAINER_NAME mysql -u $DB_USERNAME -p$DB_PASSWORD -D $DB_DATABASE < "$sql_file"
    check_error
    log_message "SQL file $sql_file execution completed successfully."
}

MIGRATIONS_DIR="./migrations"

if [ ! -d "$MIGRATIONS_DIR" ]; then
    log_message "Error: Migrations directory $MIGRATIONS_DIR not found."
    exit 1
fi

for sql_file in "$MIGRATIONS_DIR"/*.sql; do
    if [ -f "$sql_file" ]; then
        execute_sql_file "$sql_file"
    fi
done

log_message "All SQL files in $MIGRATIONS_DIR executed."
exit 0
