createpgcontainer:
	sudo docker container create --network simple-bank-network --name pgsql -p 5432:5432 -e POSTGRES_USER=tech_school -e POSTGRES_PASSWORD=21204444 postgres:15.4-alpine3.18

createdb:
	sudo docker container exec -it pgsql createdb --username=tech_school --owner=tech_school simple_bank

dropdb:
	sudo docker container exec -it pgsql dropdb --username=tech_school simple_bank

newmigrate:
	migrate create -ext sql -dir db/migration -seq $(name)

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
	sudo docker run --rm -v $(makeFileDir):/src -w /src sqlc/sqlc:1.19.1 init

sqlcgenerate:
	sudo docker run --rm -v $(makeFileDir):/src -w /src sqlc/sqlc:1.19.1 generate

test:
	go test -v -cover -short ./...

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
		  --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative --experimental_allow_proto3_optional \
          proto/*.proto


evans:
	evans --host localhost --port 9090 -r repl


createrediscontainer:
	sudo docker create --name rds -p 6379:6379 redis:7.2.0-alpine3.18


#  delete from verify_emails ;delete from sessions; delete from transfers; delete from entries ; delete from accounts ; delete from users ;
# atau jalankan saja make migdown, karna itu juga otomatis menghapus database dan isinya. lalu jalankan migup lagi

.PHONY: createpgcontainer createdb dropdb newmigrate migrateup migratedown migrateup1 migratedown1 sqlcinit sqlcgenerate test server mock proto evans createrediscontainer