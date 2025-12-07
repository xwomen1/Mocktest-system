package errors

import (
	"context"
	"fmt"
	"math"
	"time"
)

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxAttempts  int           `json:"max_attempts" yaml:"max_attempts"`
	InitialDelay time.Duration `json:"initial_delay" yaml:"initial_delay"`
	MaxDelay     time.Duration `json:"max_delay" yaml:"max_delay"`
	Multiplier   float64       `json:"multiplier" yaml:"multiplier"`
	Jitter       bool          `json:"jitter" yaml:"jitter"`

	// Which errors to retry (nil means all)
	RetryableErrors []ErrorCode `json:"retryable_errors" yaml:"retryable_errors"`
}

// DefaultRetryConfig returns sensible default retry config
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     10 * time.Second,
		Multiplier:   2.0,
		Jitter:       true,
		RetryableErrors: []ErrorCode{
			CodeNetworkError,
			CodeTimeout,
			CodeConnectionLost,
			CodeServiceUnavailable,
		},
	}
}

// Retry executes a function with retry logic
func Retry(ctx context.Context, config RetryConfig, fn func() error) error {
	var lastErr error

	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		// Execute function
		err := fn()
		if err == nil {
			return nil // Success
		}

		lastErr = err

		// Check if error is retryable
		if !isRetryableError(err, config.RetryableErrors) {
			return Wrapf(err, CodeInternalError,
				"non-retryable error after %d attempts", attempt+1)
		}

		// Check if we should stop (last attempt or context cancelled)
		if attempt == config.MaxAttempts-1 {
			break
		}

		if ctx.Err() != nil {
			return Wrapf(err, CodeTimeout, "context cancelled during retry")
		}

		// Calculate delay with exponential backoff
		delay := calculateDelay(config, attempt)

		// Wait before retry
		select {
		case <-time.After(delay):
			// Continue to next attempt
		case <-ctx.Done():
			return Wrapf(err, CodeTimeout, "context cancelled during retry delay")
		}
	}

	return Wrapf(lastErr, CodeInternalError,
		"failed after %d attempts", config.MaxAttempts)
}

// CircuitBreaker implements circuit breaker pattern
type CircuitBreaker struct {
	name                string
	failureThreshold    int
	resetTimeout        time.Duration
	halfOpenMaxAttempts int

	state            circuitState
	failureCount     int
	lastFailure      time.Time
	halfOpenAttempts int
}

type circuitState int

const (
	stateClosed circuitState = iota
	stateOpen
	stateHalfOpen
)

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(name string, failureThreshold int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		name:                name,
		failureThreshold:    failureThreshold,
		resetTimeout:        resetTimeout,
		halfOpenMaxAttempts: 3,
		state:               stateClosed,
	}
}

// Execute runs a function with circuit breaker protection
func (cb *CircuitBreaker) Execute(fn func() error) error {
	// Check circuit state
	switch cb.state {
	case stateOpen:
		if time.Since(cb.lastFailure) > cb.resetTimeout {
			// Move to half-open state
			cb.state = stateHalfOpen
			cb.halfOpenAttempts = 0
		} else {
			return New(CodeServiceUnavailable,
				fmt.Sprintf("circuit breaker '%s' is open", cb.name))
		}
	case stateHalfOpen:
		if cb.halfOpenAttempts >= cb.halfOpenMaxAttempts {
			cb.state = stateOpen
			cb.lastFailure = time.Now()
			return New(CodeServiceUnavailable,
				fmt.Sprintf("circuit breaker '%s' half-open attempts exhausted", cb.name))
		}
	}

	// Execute function
	err := fn()

	// Update circuit state
	cb.updateState(err)

	return err
}

func (cb *CircuitBreaker) updateState(err error) {
	switch cb.state {
	case stateClosed:
		if err != nil {
			cb.failureCount++
			if cb.failureCount >= cb.failureThreshold {
				cb.state = stateOpen
				cb.lastFailure = time.Now()
			}
		} else {
			cb.failureCount = 0 // Reset on success
		}

	case stateHalfOpen:
		cb.halfOpenAttempts++
		if err != nil {
			// Failed again, go back to open
			cb.state = stateOpen
			cb.lastFailure = time.Now()
			cb.failureCount = cb.failureThreshold
		} else {
			// Success, close the circuit
			cb.state = stateClosed
			cb.failureCount = 0
			cb.halfOpenAttempts = 0
		}
	}
}

// Helper functions

func isRetryableError(err error, retryableCodes []ErrorCode) bool {
	if len(retryableCodes) == 0 {
		return true // All errors are retryable
	}

	// Check if error matches any retryable code
	var appErr *Error
	if e, ok := err.(*Error); ok {
		appErr = e
	} else {
		// Wrap regular error to check
		appErr = Wrap(err, CodeInternalError, err.Error())
	}

	for _, code := range retryableCodes {
		if appErr.Code == code {
			return true
		}
	}

	return false
}

func calculateDelay(config RetryConfig, attempt int) time.Duration {
	// Exponential backoff: delay = initial * multiplier^attempt
	delay := float64(config.InitialDelay) * math.Pow(config.Multiplier, float64(attempt))

	// Cap at max delay
	if delay > float64(config.MaxDelay) {
		delay = float64(config.MaxDelay)
	}

	// Add jitter if enabled
	if config.Jitter {
		// Add ±10% jitter
		jitter := 0.1 * delay
		delay = delay - jitter/2 + (float64(time.Now().UnixNano()) * jitter / 1e9)
	}

	return time.Duration(delay)
}
