FROM postgres:15-alpine

# Copy the database initialization script
COPY init.sql /docker-entrypoint-initdb.d/
# COPY init.sh /docker-entrypoint-initdb.d/