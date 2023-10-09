#!/bin/sh
set -e

# Restore the database
litestream restore -v -if-replica-exists -o ./database.sqlite ${DB_REPLICA_REMOTE_PATH}

if [ -f ./database.sqlite ]; then
    echo "---- Restored from Cloud Storage ----"
else
    echo "---- Failed to restore from Cloud Storage ----"
    ./init_db.sh
fi

# Run litestream with your app as the subprocess.
exec litestream replicate -exec "air" ./database.sqlite ${DB_REPLICA_REMOTE_PATH}
