package helper

import (
	"context"
	"errors"
	"kg/procurement/cmd/utils"
	"time"
)

type RetryableFunc[T any] func() (T, error)

type SleepFunc func(time.Duration)

// Retry attempts to execute the provided function up to 'tries' times.
// It retries only if the error returned matches 'retryOn'.
// The delay between retries starts at 'delay' and is multiplied by 'backoff' after each attempt.
// A custom sleep function can be provided for testing purposes.
func Retry[T any](
	ctx context.Context,
	fn RetryableFunc[T],
	retryOn error,
	tries int,
	delay time.Duration,
	backoff int,
	sleep SleepFunc,
) (T, error) {
	var zeroValue T
	currentDelay := delay
	var lastError error

	for attempt := 1; attempt <= tries; attempt++ {
		select {
		case <-ctx.Done():
			return zeroValue, ctx.Err()
		default:
		}

		result, err := fn()
		if err == nil || !errors.Is(err, retryOn) {
			return result, err
		}

		lastError = err
		utils.Logger.Infof("Retry #%d for error: %v. Retrying after %v...\n", attempt, err, currentDelay)

		sleep(currentDelay)

		currentDelay *= time.Duration(backoff)
	}
	utils.Logger.Infof("Failed after %d attempts with error: %v\n", tries, lastError)

	return zeroValue, lastError
}
