package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

// type name
// this constant is very important because it's a way for asynq to recognize what kind of task it is distributing or processing
const TaskSendVerifyEmail = "task:send_verify_email"

// this struct will contain all data of the task that we want to store in Redis;
// and later the worker will be able to retrieve it from the queue
type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opt ...asynq.Option, // list of option
) error {
	jsonPayload, err := json.Marshal(payload) // karna payload berupa object json, kita harus serialize untuk menjadi
	if err != nil {
		return fmt.Errorf("failed to marshal task payload : %w", err)
	}

	// pass option argument will allow us to control how the task is distributed, run, or retried.
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opt...) // create a new task

	// now we're ready to send task to redis queue
	// taskInfor, err := distributor.client.EnqueueContext(ctx, task, opt...) // no need to add option because it already did in newtask above
	taskInfo, err := distributor.client.EnqueueContext(ctx, task) // add or enqueue task to the redis queue
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queue", taskInfo.Queue).        // show name of the queue
		Int("max_retry", taskInfo.MaxRetry). // show maximum number of retries in case of failure
		Msg("task berhasil di enqueue")

	return nil
}

// asynq has already taken care of the core part, which is pulling the task form Redis. and feed it to the background worker to process it via the task parameter of this handler funciton.
// this is a task handler function
func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail

	// parse the task to get the payload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		// if task payload is unmarshallable there's no point to retrying it. we can tell asynq about that by wrapping the asynq.SkipRetry error here
		return fmt.Errorf("gagal untuk unmarshal PayloadSendVerifyEmail dari task : %w", asynq.SkipRetry)
	}

	// retrieve user record from the database
	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		// if err == sql.ErrNoRows {
		// 	// jika gagal karna user tidak ditemukan maka tidak perlu retry
		// 	return fmt.Errorf("user tidak ditemukan : %w", asynq.SkipRetry)
		// }

		// akan lebih baik jika kita biarkan lakukan retry, karena bisa jadi error disebabkan transaksi yang lama diproses sehingga user yang dicari belum tersedia.
		// jikapun transaksi batal dilakukan dan user tidak pernah tersedia. maka retry tidak akan selamanya dilakukan karena retry memililki batas maksimal percobaan

		// in other case there's some internal error with the db, so it's retryable. therefore we simply wrap the original error
		return fmt.Errorf("gagal mendapatkan user : %w", err)

	}

	// if no error occurs, then we can send email to the user here
	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", user.Email).
		Msg("sedang memproses task") // for now just write log here TEMP

	return nil // tell asynq that task has ben processed successfully
}
