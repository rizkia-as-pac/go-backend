package gapi

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	db "github.com/tech_school/simple_bank/db/sqlc"
	"github.com/tech_school/simple_bank/pb"
	"github.com/tech_school/simple_bank/utils/pass"
	"github.com/tech_school/simple_bank/validator"
	"github.com/tech_school/simple_bank/worker"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// jika pada http server kita perlu melakukan bind pada request yang dikirimkan client. pada grpc framework itu semua sudah otomatis
	// grpc has already bound all the input data into req object for us

	// if violations not nil then it means that there's atleast one invalid parameter
	if violations := validateCreateUserRequest(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

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

	// TODO : create user and send task to Redis in 1 single DB transaction. if we fail to send task, the tx will rolled back, and the client can retry later.
	//
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, "username sudah ada : %s", err)
		}
		return nil, status.Errorf(codes.Internal, "gagal untuk membuat user : %s", err)
	}

	taskPayload := &worker.PayloadSendVerifyEmail{
		Username: user.Username,
	}
	// send verify email to user (sistributing task)
	opts := []asynq.Option{
		asynq.MaxRetry(10), // we only allows task to retried at most 10 times if it fails
		asynq.ProcessIn(10 * time.Second), // procces task after delay 10 second
		asynq.Queue(worker.QueueCritical), // sending task to different level of queue based on it importances

	}
	err = server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to distribute task to send verify email : %s", err)
	}

	// we should not mix up the DB layer struct with the API struct. because sometimes we don't want to return every field in te DB to the client. that's way we have CreateUserResponse struct and also created ConverUser function
	res := &pb.CreateUserResponse{
		User: ConvertUser(user),
	}

	return res, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, filedViolation("username", err))
	}

	if err := validator.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, filedViolation("password", err))
	}

	if err := validator.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, filedViolation("full_name", err))
	}

	if err := validator.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, filedViolation("email", err))
	}

	return violations
}
