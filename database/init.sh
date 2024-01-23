#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USERNAME" --dbname="$DB_NAME"<<-EOSQL
   CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
EOSQL