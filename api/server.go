package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/tech_school/simple_bank/db/sqlc"
	"github.com/tech_school/simple_bank/token"
	"github.com/tech_school/simple_bank/utils/conf"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config     conf.Config // we will use duration later
	store      db.Store
	tokenMaker token.Maker
	routers    *gin.Engine // router mengirimkan setiap api request ke api handler yang tepat untuk diproses
}

// NewServer membuat instansi baru dari Server
// and setup all api route untuk semua service di Server tsb
func NewServer(conf conf.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(conf.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("tidak bisa membuat token maker : %w", err)
	}

	server := &Server{
		config:     conf,
		store:      store,
		tokenMaker: tokenMaker,
	}

	// mendaftarkan custom validator ke gin
	// get current validator engine yang gin gunakan, konversi outputnya menjadi (*validator.validate)
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//v.RegisterValidation("name validation tag", function)
		v.RegisterValidation("currency", validCurrency) // untuk mendaftarkan custom validator buatan kita
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	// add routes to router
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	// daripada router. kita gunkan authRoutes untuk memasukan kedalam group route. skrng semua ruoute di group tersebut akan berbagi middleware yang sama
	// smua request yang masuk rute ini akan melewati middleware terlebih dahulu
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// menambahkan route pada router
	// jika kita mengirimkan multiple function pada handlernya
	//  maka function yang urutannya terakhir akan menjadi handler yang asli
	// function sebelumnya hanya menjadi middle ware
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	// tidak perlu uri, kita akan mendapatkan req data dari query params
	authRoutes.GET("/accounts", server.listAccounts)
	authRoutes.POST("/transfers", server.transferMoneyTechSchool)

	server.routers = router
}

// start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.routers.Run(address)
	// LATER kita bisa mengimplementasikan gracefully shadow logic di function ini
}

// errorResponse function convert error into key value object so gin can serialize it to JSON before returning to the client
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
	// LATER bisa implementasikan untuk mengembalikan sesuai tipe errornya
}
