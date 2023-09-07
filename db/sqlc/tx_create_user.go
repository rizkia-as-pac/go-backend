package db

import "context"

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error // we will run some function after user created inside the same transaction. in other word we're gonna use it as callback. error output will be used to decided wether the to commit or rollback the transaction
	// then from outside, we will use tha callback function to send the async task to redisk
}

type CreateUserTxResult struct {
	User User
}

func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.User) // execute callback function 
	})

	return result, err

}
