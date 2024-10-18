package helper

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/onsi/gomega"
)

// TestRetrySuccessOnFirstAttempt tests that the function succeeds on the first try without retries.
func TestRetrySuccessOnFirstAttempt(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	fn := func() (int, error) {
		return 42, nil
	}

	retryOn := errors.New("temporary error")
	tries := 3
	delay := 10 * time.Millisecond
	backoff := 2

	sleep := func(d time.Duration) {}

	ctx := context.Background()

	result, err := Retry(ctx, fn, retryOn, tries, delay, backoff, sleep)
	g.Expect(err).To(gomega.BeNil())
	g.Expect(result).To(gomega.Equal(42))
}

// TestRetryNonRetryableError tests that the function fails immediately with a non-retryable error.
func TestRetryNonRetryableError(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	fn := func() (int, error) {
		return 0, errors.New("fatal error")
	}
	retryOn := errors.New("temporary error")
	tries := 3
	delay := 10 * time.Millisecond
	backoff := 2

	sleep := func(d time.Duration) {}

	ctx := context.Background()

	result, err := Retry(ctx, fn, retryOn, tries, delay, backoff, sleep)
	g.Expect(err).To(gomega.MatchError("fatal error"))
	g.Expect(result).To(gomega.Equal(0))
}

// TestRetryRetryableErrorThenSuccess tests that the function fails a certain number of times with a retryable error and then succeeds.
func TestRetryRetryableErrorThenSuccess(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	retryOn := errors.New("temporary error")
	failuresBeforeSuccess := 2
	attempts := 0

	fn := func() (string, error) {
		attempts++
		if attempts <= failuresBeforeSuccess {
			return "", retryOn
		}
		return "success", nil
	}
	tries := 5
	delay := 10 * time.Millisecond
	backoff := 2

	sleep := func(d time.Duration) {}

	ctx := context.Background()

	result, err := Retry(ctx, fn, retryOn, tries, delay, backoff, sleep)
	g.Expect(err).To(gomega.BeNil())
	g.Expect(result).To(gomega.Equal("success"))
	g.Expect(attempts).To(gomega.Equal(failuresBeforeSuccess + 1))
}

// TestRetryRetryableErrorExhaustRetries tests that the function fails consistently with a retryable error, exhausting all retries.
func TestRetryRetryableErrorExhaustRetries(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	retryOn := errors.New("temporary error")
	fn := func() (bool, error) {
		return false, retryOn
	}
	tries := 3
	delay := 10 * time.Millisecond
	backoff := 2

	sleep := func(d time.Duration) {}

	ctx := context.Background()

	result, err := Retry(ctx, fn, retryOn, tries, delay, backoff, sleep)
	g.Expect(err).To(gomega.MatchError(retryOn))
	g.Expect(result).To(gomega.BeFalse())
}

// TestRetryBackoffFunctionality tests that the backoff multiplier is applied correctly.
func TestRetryBackoffFunctionality(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	retryOn := errors.New("temporary error")
	failuresBeforeSuccess := 2
	attempts := 0

	var delays []time.Duration

	fn := func() (string, error) {
		attempts++
		if attempts <= failuresBeforeSuccess {
			return "", retryOn
		}
		return "success", nil
	}

	tries := 5
	initialDelay := 10 * time.Millisecond
	backoff := 2

	sleep := func(d time.Duration) {
		delays = append(delays, d)
	}

	ctx := context.Background()
	result, err := Retry(ctx, fn, retryOn, tries, initialDelay, backoff, sleep)

	g.Expect(err).To(gomega.BeNil())
	g.Expect(result).To(gomega.Equal("success"))
	g.Expect(attempts).To(gomega.Equal(failuresBeforeSuccess + 1))

	expectedDelays := []time.Duration{
		10 * time.Millisecond, // After first failure
		20 * time.Millisecond, // After second failure
	}
	g.Expect(delays).To(gomega.Equal(expectedDelays))
}
