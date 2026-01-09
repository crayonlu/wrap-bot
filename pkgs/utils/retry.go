package utils

import (
	"fmt"
	"strings"
	"time"
)

type RetryConfig struct {
	MaxRetries    int           // 最大重试次数
	InitialDelay  time.Duration // 初始延迟
	MaxDelay      time.Duration // 最大延迟
	BackoffFactor float64       // 退避倍数
}

func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:    3,
		InitialDelay:  2 * time.Second,
		MaxDelay:      10 * time.Second,
		BackoffFactor: 2.0,
	}
}

type ShouldRetryFunc func(error) bool

func IsRateLimitError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "Too Many Requests") ||
		strings.Contains(errStr, "429") ||
		strings.Contains(errStr, "rate limit")
}

func RetryWithBackoff(operation func() error, config RetryConfig, shouldRetry ShouldRetryFunc) error {
	var lastErr error

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		err := operation()
		if err == nil {
			return nil
		}

		lastErr = err

		if shouldRetry != nil && !shouldRetry(err) {
			return fmt.Errorf("operation failed (not retryable): %w", err)
		}

		if attempt == config.MaxRetries {
			break
		}

		delay := time.Duration(float64(config.InitialDelay) * pow(config.BackoffFactor, float64(attempt)))
		if delay > config.MaxDelay {
			delay = config.MaxDelay
		}

		fmt.Printf("Retry attempt %d/%d after %v, error: %v\n", attempt+1, config.MaxRetries, delay, err)
		time.Sleep(delay)
	}

	return fmt.Errorf("operation failed after %d attempts: %w", config.MaxRetries+1, lastErr)
}

func Retry(operation func() error) error {
	return RetryWithBackoff(operation, DefaultRetryConfig(), IsRateLimitError)
}

func pow(base, exp float64) float64 {
	result := 1.0
	for i := 0; i < int(exp); i++ {
		result *= base
	}
	return result
}
