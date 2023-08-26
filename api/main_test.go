package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	db "github.com/tech_school/simple_bank/db/sqlc"
	"github.com/tech_school/simple_bank/utils/conf"
	"github.com/tech_school/simple_bank/utils/random"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	conf := conf.Config{
		TokenSymmetricKey:   random.RandomString(32, "abcdefghijklmnopqrstuvwxyz"),
		AccessTokenDuration: time.Minute,
	}

	Server, err := NewServer(conf, store)
	require.NoError(t, err)

	return Server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode) // set gin to test mode
	os.Exit(m.Run())          // to start unit test, mengembalikan pass atau fail
}
