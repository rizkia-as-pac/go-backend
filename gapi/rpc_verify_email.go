package gapi

import (
	"context"

	db "github.com/tech_school/simple_bank/db/sqlc"
	"github.com/tech_school/simple_bank/pb"
	"github.com/tech_school/simple_bank/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	// jika pada http server kita perlu melakukan bind pada request yang dikirimkan client. pada grpc framework itu semua sudah otomatis
	// grpc has already bound all the input data into req object for us

	if violations := validateVerifyEmailRequest(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	arg := db.VerifyEmailTxParams{
		EmailId:    req.GetEmailId(),
		SecretCode: req.GetSecretCode(),
	}

	txResult, err := server.store.VerifyEmailTx(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "gagal untuk memverifikasi email : %s", err)
	}

	rsp := &pb.VerifyEmailResponse{
		IsVerified: txResult.UpdatedUser.IsEmailVerified,
	}

	return rsp, nil

}

func validateVerifyEmailRequest(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateEmailId(req.GetEmailId()); err != nil {
		violations = append(violations, filedViolation("email_id", err))
	}

	if err := validator.ValidateSecretCode(req.GetSecretCode()); err != nil {
		violations = append(violations, filedViolation("secret_code", err))
	}
	return violations
}
