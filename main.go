package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // subpackage untuk database postgres dari modul migrate
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	// _ "github.com/golang-migrate/migrate/v4/source/github"  // ubah /github menjadi /file karena sumber migrate kita berada di local file system
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/tech_school/simple_bank/api"
	db "github.com/tech_school/simple_bank/db/sqlc"
	"github.com/tech_school/simple_bank/gapi"
	"github.com/tech_school/simple_bank/pb"
	"github.com/tech_school/simple_bank/utils/conf"

	_ "github.com/lib/pq"
)

// const (
// 	dbDriver      = "postgres"
// 	serverAddress = "0.0.0.0:8080"
// 	dbSource      = "postgresql://tech_school:21204444@localhost:5432/simple_bank?sslmode=disable" // copy saja dari migrate command
// )

func main() {
	// mengambil config yang sudah diberikan oleh viper
	config, err := conf.LoadConfig(".") // membaca file config dilokasi yang sama , lokasi cukup sampai pada foler yang nampung app.env saja, app.env tidak dituliska
	if err != nil {
		log.Fatal("tidak bisa membaca configuration : ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource) // create new connection to db
	if err != nil {
		log.Fatal("tidak bisa terkoneksi dengan database : ", err)
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)

	// runGinServer(config,store) // gin server

	// go runGatewayServer(config, store) // go membuat run gateway berjalan di routine yang berbeda. runGatewayerver dan runGRPCServer berjalan pada routine yang sama maka server yang pertama akan memblock server yang kedua, sehingga kita harus memisahkan nya dari main go  routine dengan go routine yang lain
	runGRPCServer(config, store)

}

func runGRPCServer(config conf.Config, store db.Store) {
	// implementasi server bank mandiri kita  // sama seperti newserver yang ada pad gin
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("tidak bisa membuat server : ", err)
	}

	grpcServer := grpc.NewServer() // create new grpc server object
	pb.RegisterSimpleBankServer(grpcServer, server)

	//OPTIONAL TAPI SANGAT DIREKOMENDASIKAN
	// register grpc reflection for our server
	// meskipun sederhana namun ini powerfull
	// ini allow grpc client to easily explore what rpc are avaible of the server and how to call them
	// ini seperti self documentation dari server
	// mengaktifkan reflection juga membuat kita dapat menggunakan tool seperti evans cli
	reflection.Register(grpcServer)

	// menggunakan protocol tcp
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("tidak bisa membuat listener : ", err)
	}

	log.Printf("memulai server GRPC pada %s", listener.Addr().String())

	// start server
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("tidak bisa memulai GRPC server : ", err)
	}

}

// setup HTTP gateway with in-process translation method
func runGatewayServer(config conf.Config, store db.Store) {
	// implementasi server bank mandiri kita  // sama seperti newserver yang ada pad gin
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("tidak bisa membuat server : ", err)
	}

	// this function comes from runtime pacakge which is a sub-package of grpc-gateway v2
	grpcMux := runtime.NewServeMux(
		// secara default JSON response yang akan dikirimkan ke client akan camelCase. jika kita ingin mengcustomisasinya menjadi snake_case tambahkan berikut https://grpc-ecosystem.github.io/grpc-gateway/docs/mapping/customizing_your_gateway/ :
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	ctx, cancel := context.WithCancel(context.Background()) // create context. return context object and function to cancel it
	defer cancel()                                          // we will defer the cancel call so that only will be executed before exiting this runGatewayServer function
	// canceling a context is just the way to prevent the system to doing unnecessary work

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("tidak bisa register handler server : ", err)
	}

	// create a new http serve mux
	mux := http.NewServeMux() // this mux will actualy receive http req from client

	// convert mux(client req) to grpc format. dengan cara reroutes ke grpcmux
	// single slash (/) is use to cover all routes
	mux.Handle("/", grpcMux)

	// menggunakan protocol tcp
	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("tidak bisa membuat listener : ", err)
	}

	log.Printf("memulai server HTTP GATEWAY pada %s", listener.Addr().String())

	// start server
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("tidak bisa memulai HTTP GATEWAY server : ", err)
	}

}

func runGinServer(config conf.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("tidak bisa membuat server : ", err)
	}

	// start server
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("Tidak bisa memulai server : ", err)
	}
}

func runDBMigration(migrationUrl string, dbSource string) {
	migration, err := migrate.New(migrationUrl, dbSource) // return migration object
	if err != nil {
		log.Fatal("tidak bisa membuat instansiasi migrate : ", err)
	}

	if err := migration.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("tidak ada perubahan pada schema migration")
		} else {
			log.Fatal("gagal menjalankan migrate up : ", err)
		}
	} // run all the migration up files

	// jika tidak ada error yang muncul
	log.Println("db migration berhasil dijalankan")
}
