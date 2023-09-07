package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	db "github.com/tech_school/simple_bank/db/sqlc"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error

	// function to process the task send verify email
	// it must follow asynq's task handler function  signature (ctx context.Context, task *asynq.Task) error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

// RedisTaskProcessor will implement TaskProcessor interface\
// just like how we did for the task distributor
type RedisTaskProcessor struct {
	server *asynq.Server // RedisTaskProcessor must contain an asynq.server object as one of its field.
	store  db.Store      // when processing the task it will need database
}

// redis client opt to connect to redis
func NewRedisTaskProcessor(redisOpt *asynq.RedisClientOpt, store db.Store) TaskProcessor {
	customLogger := NewCustomLogger()
	redis.SetLogger(customLogger) // apply custom logger to go-redis package internal logger

	server := asynq.NewServer(
		redisOpt,

		// asynq.config object allow us to control many different parameters of asynq server. for example  Maximum number of concurrent processing of tasks, retry delay for a failed task, a predicate function to menentukan apakah error yang dikembalikan dari handler adalah sebuah kegagalan atau bukan, a map of task queues together with their priority values, and many more.
		// for now keep it simple and leave it config empty. yang berarti kita akan menggunakan asynq's predefined default configurations.
		asynq.Config{
			// queues map tell asyncq about the queue names and their correspondeing priority values.
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
				// "low":      1,
			},

			// it is in fact a type conversion, but for function instead of normal variable
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().
					Err(err).
					Str("type", task.Type()).
					Bytes("payload", task.Payload()).
					Msg("task gagal di proses")

					// you can modif this error handler function to send notification to your email, slack or whatever channel you want
			}), // add error handler

			Logger: customLogger, // register our custom logger 


			// specify a custom logger for the asynq server
		},
	)

	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux() // create a new mux
	// asynq design is prety simmilar to that http server. we can use mux to register each task with its handler function. just like how we use http mux to register each route

	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Start(mux)
}
