package gapi

import (
	"fmt"

	db "github.com/tech_school/simple_bank/db/sqlc"
	"github.com/tech_school/simple_bank/pb"
	"github.com/tech_school/simple_bank/token"
	"github.com/tech_school/simple_bank/utils/conf"
)

// Server serves GRPC requests for our banking service.
type Server struct {
	// enable forward compability. membuat rpc seperti CreateUser dan LoginUser meskipun belum diimplementasikan namun sudah bisa diterima pada server dari client sebagai testing. hal ini membuat kita bisa mudah bekerja dalam tim membuat banyak rpc sekaligus secara paralel tanpa blok atau konflik dengan satu sama lain
	// hal ini juga membuat kita bisa mengakses service apa saja melalui evans tool cli
	pb.UnimplementedSimpleBankServer
	config     conf.Config // we will use duration later
	store      db.Store
	tokenMaker token.Maker
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

	// server.setupRouter() // tidak seperti HTTP, tidak ada route di grpc
	// client akan memanggil server dengan mengeksekusi sebuah rpc sama seperti memanggil function pada umumnya
	return server, nil
}
