#!/bin/bash

# make sure to change the main.go file if you change this
DB_FILE="database.db"

if [ -f "$DB_FILE" ]; then
	echo "error: file '$DB_FILE' already exists."
	exit 1
fi

sqlite3 "$DB_FILE" "CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    token TEXT
);" || { echo "Failed to create table"; exit 1; }

echo "Database was initiated successfully"
