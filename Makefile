createpgcontainer:
	sudo docker container create --network simple-bank-network --name pgsql -p 5432:5432 -e POSTGRES_USER=tech_school -e POSTGRES_PASSWORD=21204444 postgres:15.4-alpine3.18

createdb:
	sudo docker container exec -it pgsql createdb --username=tech_school --owner=tech_school simple_bank

dropdb:
	sudo docker container exec -it pgsql dropdb --username=tech_school simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://tech_school:21204444@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://tech_school:21204444@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://tech_school:21204444@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://tech_school:21204444@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

makeFileDir := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

sqlcinit:
	sudo docker run --rm -v $(makeFileDir):/src -w /src sqlc/sqlc:1.8.0 init

sqlcgenerate:
	sudo docker run --rm -v $(makeFileDir):/src -w /src sqlc/sqlc:1.8.0 generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/tech_school/simple_bank/db/sqlc Store

# https://grpc.io/docs/languages/go/quickstart/
# rm -f pb/*.go = hapus semua go code yang ada di pb sebelum generate yang baru
# proto:
# 	rm -f pb/*.go
# 	rm -f doc/swagger/*.swagger.json
# 	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
#     --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
# 	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
# 	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=bank_mandiri \
#     proto/*.proto 
# *.proto karena kita mau generate dari semua proto file yang ada difile proto
proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
          --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		  --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
          proto/*.proto


evans:
	evans --host localhost --port 9090 -r repl


#  delete from transfers; delete from entries ; delete from accounts ; delete from users ;

.PHONY: createpgcontainer createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlcinit sqlcgenerate test server mock proto evans