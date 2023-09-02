package gapi

import (
	"context"

	db "github.com/tech_school/simple_bank/db/sqlc"
	"github.com/tech_school/simple_bank/pb"
	"github.com/tech_school/simple_bank/utils/pass"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// jika pada http server kita perlu melakukan bind pada request yang dikirimkan client. pada grpc framework itu semua sudah otomatis

	// grpc has already bound all the input data into req object for us

	// hashedPassword, err := pass.HashedPassword(req.Password)
	hashedPassword, err := pass.HashedPassword(req.GetPassword()) // // bisa seperti atas, tapi menggunakan req.GetPassword() lebih baik karna sudah ada pengecekan terlebih dahulu passwordnya, just in case the request field is nil
	if err != nil {
		// status adalah sub package grpc
		// codes juga adalah sub package grpc
		return nil, status.Errorf(codes.Internal, "failed to hash password : %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, "username sudah ada : %s", err)
		}
		return nil, status.Errorf(codes.Internal, "gagal untuk membuat user : %s", err)
	}

	// we should not mix up the DB layer struct with the API struct. because sometimes we don't want to return every field in te DB to the client. that's way we have CreateUserResponse struct and also created ConverUser function
	res := &pb.CreateUserResponse{
		User: ConvertUser(user),
	}

	return res, nil
}
