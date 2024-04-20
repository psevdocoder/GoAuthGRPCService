protoc-gen:
	protoc -I protos/proto protos/proto/sso/sso.proto --go_out=./protos/gen/go/ --go_opt=paths=source_relative --go-grpc_out=./protos/gen/go/ --go-grpc_opt=paths=source_relative

run:
	go run cmd/sso/main.go

run-migrator:
	go run cmd/migrator/main.go --storage-path=./storage/database.db --migrations-path=./migrations

run-migrator-tests:
	go run .\cmd\migrator\main.go --storage-path=./storage/database.db --migrations-path=./tests/migrations --migrations-table=migrations_test

migrate:
	migrate -path=./migrations -database "sqlite3://storage/database.db" up

migrate-down:
	migrate -path=./migrations -database "sqlite3://storage/database.db" down 1

migrate-install-sqlite:
	go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-install-postgres:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
