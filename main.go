package main

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // subpackage untuk database postgres dari modul migrate

	// _ "github.com/golang-migrate/migrate/v4/source/github"  // ubah /github menjadi /file karena sumber migrate kita berada di local file system
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/tech_school/simple_bank/api"
	db "github.com/tech_school/simple_bank/db/sqlc"
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

	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("tidak bisa membuat server : ", err)
	}

	err = server.Start(config.ServerAddress)
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
