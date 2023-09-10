package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	TransferTxV2(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
}

/* OLD LIB/PQ
  	// SQLStore provide all functions to execute db queries and transaction
	type SQLStore struct {
		// composition = extend struct functionality in golang instead inheritance
		*Queries // embed interface QUerier membuat store interface memiliki semua function yang dimiliki oleh interface querier
		db       *sql.DB
	}z

	// membuat new store
	func NewStore(db *sql.DB) Store {
		// membuat newstore object dan mengembalikannya
		return &SQLStore{
			Queries: New(db), // new membuat dan mengembalikan object Queries
			db:      db,
		}
	}
*/

// NEW PGX

// SQLStore provide all functions to execute db queries and transaction
type SQLStore struct {
	// composition = extend struct functionality in golang instead inheritance
	*Queries // embed interface QUerier membuat store interface memiliki semua function yang dimiliki oleh interface querier
	connPool *pgxpool.Pool
}

// membuat new store
func NewStore(connPool *pgxpool.Pool) Store {
	// # we always need a pool of connection (multiple connection) in order to handle multiple request in parallel. the pgx pool package will help us manage the connection pools

	// fungsi conPool sama dengan db yang kita gunakan pada lib/pq sebelumnya, yaitu menyediakan hal hal yang didapat dari koneksi dengan database postgres

	// membuat newstore object dan mengembalikannya
	return &SQLStore{
		Queries:  New(connPool),
		connPool: connPool,
	}
}
