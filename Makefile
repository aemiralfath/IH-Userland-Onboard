postgresh:
	docker exec -it ih-userland-onboard_postgres_1 /bin/sh

createdb:
	docker exec -it ih-userland-onboard_postgres_1 createdb --username=admin --owner=admin userland2

dropdb:
	docker exec -it ih-userland-onboard_postgres_1 dropdb userland2 -U admin

checkdb:
	docker exec -it ih-userland-onboard_postgres_1 psql userland -U admin -c "\d users"

migrateup:
	migrate -path datastore/migrations -database "postgres://admin:admin@localhost:5431/userland?sslmode=disable" -verbose up

migratedown:
	migrate -path datastore/migrations -database "postgres://admin:admin@localhost:5431/userland?sslmode=disable" -verbose down

migratedirty:
	migrate -path datastore/migrations -database "postgres://admin:admin@localhost:5431/userland?sslmode=disable" force 000001

.PHONY: postgres createdb dropdb migrateup migratedown checkdb migratedirty