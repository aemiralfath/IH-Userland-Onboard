postgresh:
	docker exec -it ih-userland-onboard_postgres_1 /bin/sh

createdb:
	docker exec -it ih-userland-onboard_postgres_1 createdb --username=admin --owner=admin userland2

dropdb:
	docker exec -it ih-userland-onboard_postgres_1 dropdb userland2 -U admin

migrateup:
	migrate -path datastore/migrations -database "postgres://admin:admin@localhost:5438/userland?sslmode=disable" -verbose up

migratedown:
	migrate -path datastore/migrations -database "postgres://admin:admin@localhost:5438/userland?sslmode=disable" -verbose down

.PHONY: postgres createdb dropdb migrateup migratedown