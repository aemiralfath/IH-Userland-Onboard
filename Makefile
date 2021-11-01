userlandsh:
	docker exec -it ih-userland-onboard_userland_1 /bin/sh

redissh:
	docker exec -it ih-userland-onboard_redis_1 /bin/sh

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

mockstore:
	mockgen -destination datastore/mock/user_store.go github.com/aemiralfath/IH-Userland-Onboard/datastore UserStore
	mockgen -destination datastore/mock/session_store.go github.com/aemiralfath/IH-Userland-Onboard/datastore SessionStore
	mockgen -destination datastore/mock/client_store.go github.com/aemiralfath/IH-Userland-Onboard/datastore ClientStore
	mockgen -destination datastore/mock/profile_store.go github.com/aemiralfath/IH-Userland-Onboard/datastore ProfileStore
	mockgen -destination datastore/mock/password_store.go github.com/aemiralfath/IH-Userland-Onboard/datastore PasswordStore
	mockgen -destination datastore/mock/otp_store.go github.com/aemiralfath/IH-Userland-Onboard/datastore OTPStore
	mockgen -destination datastore/mock/crypto.go github.com/aemiralfath/IH-Userland-Onboard/datastore Crypto

.PHONY: postgres createdb dropdb migrateup migratedown checkdb migratedirty