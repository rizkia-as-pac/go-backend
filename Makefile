createpgcontainer:
	sudo docker container create --name pgsql -p 5432:5432 -e POSTGRES_USER=tech_school -e POSTGRES_PASSWORD=21204444 postgres:15.4-alpine3.18

createdb:
	sudo docker container exec -it pgsql createdb --username=tech_school --owner=tech_school simple_bank

dropdb:
	sudo docker container exec -it pgsql dropdb --username=tech_school simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://tech_school:21204444@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://tech_school:21204444@localhost:5432/simple_bank?sslmode=disable" -verbose down

makeFileDir := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

sqlcinit:
	sudo docker run --rm -v $(makeFileDir):/src -w /src sqlc/sqlc:1.8.0 init

sqlcgenerate:
	sudo docker run --rm -v $(makeFileDir):/src -w /src sqlc/sqlc:1.8.0 generate

.PHONY: createpgcontainer createdb dropdb migrateup migratedown sqlcinit sqlcgenerate