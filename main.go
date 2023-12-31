package main

import (
	"context"
	"net"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // subpackage untuk database postgres dari modul migrate
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	// _ "github.com/golang-migrate/migrate/v4/source/github"  // ubah /github menjadi /file karena sumber migrate kita berada di local file system
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tech_school/simple_bank/api"
	db "github.com/tech_school/simple_bank/db/sqlc"
	"github.com/tech_school/simple_bank/gapi"
	"github.com/tech_school/simple_bank/mail"
	"github.com/tech_school/simple_bank/pb"
	"github.com/tech_school/simple_bank/utils/conf"
	"github.com/tech_school/simple_bank/worker"
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
		log.Fatal().Err(err).Msg("tidak bisa membaca configuration")
	}

	// ZEROLOG pretty-logging mode
	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	/* OLD LIB/PQ
	conn, err := sql.Open(config.DBDriver, config.DBSource) // create new connection to db
	if err != nil {
		log.Fatal().Err(err).Msg("tidak bisa terkoneksi dengan database")
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)
	*/

	// NEW PGX

	connPool, err := pgxpool.New(context.Background(), config.DBSource) // create new connection pool
	if err != nil {
		log.Fatal().Err(err).Msg("tidak bisa terkoneksi dengan database")
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(connPool)

	// this asynq.RedisClientOpt object allows us to set up many different parameters to communicate with the Redis server.
	// asynq.RedisClientOpt use to know how to connect to reddis
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
		// Username:     "", // username to authenticate the connection // for this lecture we use default value
		// Password:     "", // password to authenticate the connection // for this lecture we use default value
		// DB:           0, // redis DB number // for this lecture we use default value
		// TLSConfig:    &tls.Config{}, // if connect to production redis server that use TLS to secure the connections, you must set appropriate value for this field
	}

	// create a new task distributor
	taskDist := worker.NewRedisTaskDistributor(&redisOpt)

	// runGinServer(config,store) // gin server

	go runTaskProcessor(config, redisOpt, store) // run task processor in seprate go routine because when asynq server start it will block and keep plling redis for new tasks like http web server so it block like http server block while waiting request from client.
	go runGatewayServer(config, store, taskDist) // go membuat run gateway berjalan di routine yang berbeda. runGatewayerver dan runGRPCServer berjalan pada routine yang sama maka server yang pertama akan memblock server yang kedua, sehingga kita harus memisahkan nya dari main go  routine dengan go routine yang lain
	runGRPCServer(config, store, taskDist)

}

func runTaskProcessor(conf conf.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	gmailSender := mail.NewGmailSender(conf.EmailSenderName, conf.EmailSenderAddress, conf.EmailSenderPassword)
	// create task processor
	taskProcessor := worker.NewRedisTaskProcessor(&redisOpt, store, gmailSender)
	log.Info().Msg("memulai redis taks processor")

	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("gagal memulai redis task processor")
	}
}

func runGRPCServer(config conf.Config, store db.Store, td worker.TaskDistributor) {
	// implementasi server bank mandiri kita  // sama seperti newserver yang ada pad gin
	server, err := gapi.NewServer(config, store, td)
	if err != nil {
		log.Fatal().Err(err).Msg("tidak bisa membuat server")
	}

	// karna semua rpc pada kasus kali ini adalah unary maka yang ini
	// function ini akan melakukan set interceptor pada server
	grpcLoggerServerOption := grpc.UnaryInterceptor(gapi.GRPCLogger)

	grpcServer := grpc.NewServer(grpcLoggerServerOption) // create new grpc server object. this funtction can take multiple server option
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
		log.Fatal().Err(err).Msg("tidak bisa membuat listener")
	}

	log.Info().Msgf("memulai server GRPC pada %s", listener.Addr().String())

	// start server
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("tidak bisa memulai GRPC server")
	}

}

// setup HTTP gateway with in-process translation method
func runGatewayServer(config conf.Config, store db.Store, td worker.TaskDistributor) {
	// implementasi server bank mandiri kita  // sama seperti newserver yang ada pad gin
	server, err := gapi.NewServer(config, store, td)
	if err != nil {
		log.Fatal().Err(err).Msg("tidak bisa membuat server")
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
		log.Fatal().Err(err).Msg("tidak bisa register handler server")
	}

	// create a new http serve mux
	mux := http.NewServeMux() // this mux will actualy receive http req from client

	// convert mux(client req) to grpc format. dengan cara reroutes ke grpcmux
	// single slash (/) is use to cover all routes
	mux.Handle("/", grpcMux)

	// menggunakan protocol tcp
	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("tidak bisa membuat listener")
	}

	log.Info().Msgf("memulai server HTTP GATEWAY pada %s", listener.Addr().String())

	handler := gapi.HTTPLogger(mux) // return new http handler with logger middleware // added http logger to the gateway server

	// start server
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Err(err).Msg("tidak bisa memulai HTTP GATEWAY server")
	}

}

func runGinServer(config conf.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("tidak bisa membuat server")
	}

	// start server
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("Tidak bisa memulai server")
	}
}

func runDBMigration(migrationUrl string, dbSource string) {
	migration, err := migrate.New(migrationUrl, dbSource) // return migration object
	if err != nil {
		log.Fatal().Err(err).Msg("tidak bisa membuat instansiasi migrate")
	}

	if err := migration.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Info().Msg("tidak ada perubahan pada schema migration")
		} else {
			log.Fatal().Err(err).Msg("gagal menjalankan migrate up")
		}
	} // run all the migration up files

	// jika tidak ada error yang muncul
	log.Info().Msg("db migration berhasil dijalankan")
}
