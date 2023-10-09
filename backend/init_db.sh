#!/bin/bash

if [ ! -e ./database.sqlite ]; then
    cat _tools/sqlite3/pragma.sql _tools/sqlite3/schema.sql | sqlite3 database.sqlite
fi
