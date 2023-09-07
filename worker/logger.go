package worker

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type CustomLogger struct{}

// create a new CustomLogger
// return asynq.Logger so we can use golang interface check feature
// func NewCustomLogger() asynq.Logger {

// kita mengembalikan *CustomLogger daripada asynq.Logger interface karena fungsi redis.SetLogger(customLogger) mendeteksi  asynq.Logger yang tidak menerapkan PrintF yang merupakan method dari interface internal.Logging. hal ini mengakibatkan error meskipun CustomLogger sudah menerapkan Print
// karna itulah kita mengembalikan *CustomLogger yang didalamnya sudah menerapkan kedua interface diatas 
func NewCustomLogger() *CustomLogger {
	return &CustomLogger{}
}

// args... = rest argument
func (logger *CustomLogger) Print(level zerolog.Level, args ...interface{}) {
	mergedString := fmt.Sprint(args...) // merge all input string together into one string format
	log.WithLevel(level).
		Msg(mergedString)
}

func (logger *CustomLogger) Printf(ctx context.Context, format string, v ...interface{}) {
	log.WithLevel(zerolog.DebugLevel).Msgf(format, v...)
}

// Debug logs a message at Debug level.
func (logger *CustomLogger) Debug(args ...interface{}) {
	logger.Print(zerolog.DebugLevel, args...)
}

// Info logs a message at Info level.
func (logger *CustomLogger) Info(args ...interface{}) {
	logger.Print(zerolog.InfoLevel, args...)
}

// Warn logs a message at Warning level.
func (logger *CustomLogger) Warn(args ...interface{}) {
	logger.Print(zerolog.WarnLevel, args...)
}

// Error logs a message at Error level.
func (logger *CustomLogger) Error(args ...interface{}) {
	logger.Print(zerolog.WarnLevel, args...)
}

// Fatal logs a message at Fatal level
// and process will exit with status set to 1.
func (logger *CustomLogger) Fatal(args ...interface{}) {
	logger.Print(zerolog.FatalLevel, args...)

}
