#!/bin/bash

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

database_file="$script_dir/../pkg/database/database.db"
schema_file="$script_dir/../pkg/database/schema.sql"
# seed_data_file="$script_dir/../pkg/database/seedData.sql"

if [ -e "$database_file" ]; then
    echo "Database file already exists. Doing nothing."
else
    sqlite3 "$database_file" < "$schema_file"
#    sqlite3 "$database_file" < "$seed_data_file"
    echo "Database created with schema."
fi