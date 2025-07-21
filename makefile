run:
	templ generate --watch --proxy="http://localhost:8080" --cmd="go run ."
start:
	templ generate
	air \
		--build.cmd "go build -o bin/app" \
		--build.bin "./bin/app" \
		--build.exclude_dir "vendor" \
		--build.exclude_dir "static" \
		--build.include_ext "go" \
		--build.kill_delay "0.5s" \
		--build.poll "2s"

c_m: # create-migration: create migration of name=<migration_name>
	migrate create -ext sql -dir db/migrations -seq $(name)

count ?= 1
version ?= 1
db_username ?= postgres
db_password ?= admin
db_host ?= localhost
db_port ?= 5432
db_name ?= goldsavings_db
ssl_mode ?= disable

# Run database migrations up to apply pending changes
m_up: # migrate-up
	migrate -path db/migrations -database "postgres://${db_username}:${db_password}@${db_host}:${db_port}/${db_name}?sslmode=${ssl_mode}" up $(count)

# Fix dirty database state by forcing to previous clean version
m_fix: # migrate-fix: fix dirty database state
	migrate -path db/migrations -database "postgres://${db_username}:${db_password}@${db_host}:${db_port}/${db_name}?sslmode=${ssl_mode}" force $(version)

# Check current migration version number
m_version: # migrate-version
	migrate -path db/migrations -database "postgres://${db_username}:${db_password}@${db_host}:${db_port}/${db_name}?sslmode=${ssl_mode}" version

# Force database migration version without running migrations
m_fup: # migreate-force up
	migrate -path db/migrations -database "postgres://${db_username}:${db_password}@${db_host}:${db_port}/${db_name}?sslmode=${ssl_mode}" force $(count)

# Roll back database migrations
m_down: # migrate-down
	migrate -path db/migrations -database "postgres://${db_username}:${db_password}@${db_host}:${db_port}/${db_name}?sslmode=${ssl_mode}" down $(count)

# Start services - PostgreSQL | Bitgo | Redis containers in detached mode
s_up: #
	docker compose -f docker-compose.services.yml up -d

# Stop and remove services - PostgreSQL | Bitgo | Redis container
s_down: #
	docker compose -f docker-compose.services.yml down

container_name ?= goldsavings_postgres

# Create a new PostgreSQL database
db_up: # database-up: create a new database
	docker exec -it ${container_name} createdb --username=${db_username} --owner=${db_username} ${db_name}

# Create a full backup of the database
db_backup:
	docker exec -it ${container_name} pg_dump --username=${db_username} ${db_name} > db_backup.sql

# Restore database from a full backup
db_restore:
	docker exec -i ${container_name} psql --username=${db_username} ${db_name} < db_backup.sql

# Backup specific tables from the database
# Usage: make db_backup_specific tables="table1 table2 table3"
db_backup_specific:
	docker exec -it ${container_name} pg_dump --username=${db_username} ${db_name} --table=$(subst $(space),$(,),$(tables)) > db_backup_specific.sql

# Restore specific tables from a backup
# Usage: make db_restore_specific tables="table1 table2 table3"
db_restore_specific:
	docker exec -i ${container_name} psql --username=${db_username} ${db_name} < db_backup_specific.sql

# Drop/delete the database
db_down: # database-down: drop a database
	docker exec -it ${container_name} dropdb --username=${db_username} ${db_name}

# Generate Go code from SQL using sqlc
sqlc: # sqlc-generate
	sqlc generate