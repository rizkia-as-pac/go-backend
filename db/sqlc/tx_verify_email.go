package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type VerifyEmailTxParams struct {
	EmailId    int64
	SecretCode string
}

type VerifyEmailTxResult struct {
	UpdatedUser        User
	UpdatedVerifyEmail VerifyEmail
}

func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.UpdatedVerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return fmt.Errorf("failed to update ver email : %s", err)
		}

		result.UpdatedUser, err = q.UpdateUser(ctx, UpdateUserParams{
			/* OLD LIB/PQ
				IsEmailVerified: sql.NullBool{
					Bool:  true,
					Valid: true,
				},
			*/
			IsEmailVerified: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},

			Username: result.UpdatedVerifyEmail.Username,
		})
		if err != nil {
			return fmt.Errorf("failed to update user : %s", err)

		}

		return err
	})

	return result, err

}
