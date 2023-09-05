package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opt ...asynq.Option, // list of option 
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client // we will use this client later to send the task to redis queue
}

// the reason we return TaskDistributor interface is we're forcing the RedisTaskDistributor to implement the TaskDistributor interface. if it doesn't implement all required functions of the interface, the compiler will complain
// we want to use automatic type checking feature of go compiler
func NewRedisTaskDistributor(redisOpt *asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt) // create a new client

	return  &RedisTaskDistributor{
		client: client,
	}
}

