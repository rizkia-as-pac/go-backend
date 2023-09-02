package gapi

import (
	db "github.com/tech_school/simple_bank/db/sqlc"
	"github.com/tech_school/simple_bank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt), // timestamp type in protobuf is not same as the golang's time type
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
