package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

// untuk sementara kita gunakan konstanta, dalam real case kita akan menarik data dari environment variable
const (
	dbDriver = "postgres"
	dbSource = "postgresql://tech_school:21204444@localhost:5432/simple_bank?sslmode=disable" // copy saja dari migrate command
)

var testQueries *Queries // didefinisikan secara global karena akan menggunakannya secara intensif di semua unit test kita
var testDB *sql.DB

// the convention the TestMain function is main entry point to all unit test inside one spesific golang package
func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource) // create new connection to db
	if err != nil {
		log.Fatal("tidak bisa terkoneksi dengan database : ", err)
	}

	testQueries = New(testDB) // function new dari file  yang digen sqlc
	os.Exit(m.Run())          // to start unit test, mengembalikan pass atau fail
}
