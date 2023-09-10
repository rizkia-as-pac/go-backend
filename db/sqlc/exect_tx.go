package db

import (
	"context"
	"fmt"
)

/* OLD LIB/PQ
  	// function yang mengeksekusi generic database transaction
	// execTx executes a function within a database transaction
	func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
		// mulai transaction baru
		// &sql.TxOptions{} membolehkan kita untuk melakukan custom level dari isolation. jika tidak diterapkan maka level nya akan default sesuai dengan jenis database yg digukanan. untuk postgres levelnya write commited
		// tx,err := store.db.BeginTx(ctx, &sql.TxOptions{})
		tx, err := store.db.BeginTx(ctx, nil) // start transaction
		if err != nil {
			return err
		}

		q := New(tx) // sama seperti New pada Store struct, bedanya di store adalah sql.DB disini sql.Tx. ini bisa dilakukan karena New menerima DBTX yang merupakan singkatan DB dan TX
		err = fn(q)  // memanggil input function didalam tx queries

		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("tx err: %v, rb err : %v", err, rbErr) // kalau error pada tx terjadi dan terjadi juga error pada proses rollback
			}

			return err // jika rollback sukses, kita hanya perlu mengembalikan error pada tx saja
		}

		// jika tidak ada kendala dan seluruh operasi sukses maka lakukan commit tx
		// kembalikan error pada caller jika ada
		return tx.Commit()
	}
*/

// NEW PGX
// function yang mengeksekusi generic database transaction
// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	// mulai transaction baru
	// &sql.TxOptions{} membolehkan kita untuk melakukan custom level dari isolation. jika tidak diterapkan maka level nya akan default sesuai dengan jenis database yg digukanan. untuk postgres levelnya write commited
	// tx,err := store.db.BeginTx(ctx, &sql.TxOptions{})
	tx, err := store.connPool.Begin(ctx) // start transaction // behind the scene it called connPool.beginTx object but with empty transaction option
	if err != nil {
		return err
	}

	q := New(tx) // sama seperti New pada Store struct, bedanya di store adalah sql.DB disini sql.Tx. ini bisa dilakukan karena New menerima DBTX yang merupakan singkatan DB dan TX
	err = fn(q)  // memanggil input function didalam tx queries

	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err : %v", err, rbErr) // kalau error pada tx terjadi dan terjadi juga error pada proses rollback
		}

		return err // jika rollback sukses, kita hanya perlu mengembalikan error pada tx saja
	}

	// jika tidak ada kendala dan seluruh operasi sukses maka lakukan commit tx
	// kembalikan error pada caller jika ada
	return tx.Commit(ctx)
}
