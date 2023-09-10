package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tech_school/simple_bank/utils/conf"
)

// untuk sementara kita gunakan konstanta, dalam real case kita akan menarik data dari environment variable
// const (
// 	dbDriver = "postgres"
// 	dbSource = "postgresql://tech_school:21204444@localhost:5432/simple_bank?sslmode=disable" // copy saja dari migrate command
// )

/* OLD LIB/PQ
  	var testStore *Queries // didefinisikan secara global karena akan menggunakannya secara intensif di semua unit test kita
	var testDB *sql.DB

	// the convention the TestMain function is main entry point to all unit test inside one spesific golang package
	func TestMain(m *testing.M) {
		config, err := conf.LoadConfig("../..")
		if err != nil {
			log.Fatal("tidak bisa membaca configuration : ", err)
		}

		testDB, err = sql.Open(config.DBDriver, config.DBSource) // create new connection to db
		if err != nil {
			log.Fatal("tidak bisa terkoneksi dengan database : ", err)
		}

		testStore = New(testDB) // function new dari file  yang digen sqlc
		os.Exit(m.Run())          // to start unit test, mengembalikan pass atau fail
	}
*/

// NEW PGX

var testStore Store // didefinisikan secara global karena akan menggunakannya secara intensif di semua unit test kita

// the convention the TestMain function is main entry point to all unit test inside one spesific golang package
func TestMain(m *testing.M) {
	config, err := conf.LoadConfig("../..")
	if err != nil {
		log.Fatal("tidak bisa membaca configuration : ", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource) // create new connection pool
	if err != nil {
		log.Fatal("tidak bisa terkoneksi dengan database : ", err)
	}

	testStore = NewStore(connPool)
	os.Exit(m.Run()) // to start unit test, mengembalikan pass atau fail
}
