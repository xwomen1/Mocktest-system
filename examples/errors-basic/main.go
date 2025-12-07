package main

import (
	"context"
	"fmt"
	"time"

	"upm-simple/pkg/errors"
)

func main() {
	fmt.Println("=== Error Handling Framework Example ===")

	// basic error creation
	fmt.Println("\n1. Basic Error Creation:")
	err1 := errors.New(errors.CodeInvalidArgument, "user ID cannot be empty")
	fmt.Printf("Error: %v\n", err1)
	fmt.Printf("Code: %s, HTTP Status: %d\n", err1.Code, err1.HTTPStatus())
	fmt.Printf("Is client error? %v\n", err1.IsClientError())

	// error with metadata
	fmt.Println("\n2. Error with Metadata:")
	err2 := errors.New(errors.CodeNotFound, "user not found")
	err2 = errors.AddMetadata(err2, "user_id", "12345")
	err2 = errors.AddMetadata(err2, "attempted_operation", "get_user")
	fmt.Printf("Error: %v\n", err2)
	fmt.Printf("Metadata: %v\n", err2.Metadata)

	// error wrapping
	fmt.Println("\n3. Error Wrapping:")
	dbErr := fmt.Errorf("connection refused")
	err3 := errors.Wrap(dbErr, errors.CodeNetworkError, "failed to connect to database")
	fmt.Printf("Wrapped error: %v\n", err3)
	fmt.Printf("Original cause: %v\n", err3.Cause)

	// utility functions
	fmt.Println("\n4. Utility Functions:")
	notFoundErr := errors.NotFoundError("user", "john@example.com")
	fmt.Printf("Not found error: %v\n", notFoundErr)

	validationErr := errors.ValidationError("email", "must be valid email address")
	fmt.Printf("Validation error: %v\n", validationErr)

	// retry mechanism
	fmt.Println("\n5. Retry Mechanism:")
	attempt := 0
	retryFunc := func() error {
		attempt++
		fmt.Printf("  Attempt %d... ", attempt)
		if attempt < 3 {
			fmt.Println("Failed")
			return fmt.Errorf("temporary failure")
		}
		fmt.Println("Success!")
		return nil
	}

	retryConfig := errors.DefaultRetryConfig()
	retryConfig.MaxAttempts = 5

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := errors.Retry(ctx, retryConfig, retryFunc)
	if err != nil {
		fmt.Printf("Retry failed: %v\n", err)
	} else {
		fmt.Println("Retry succeeded!")
	}

	fmt.Println("\n6. Error Checking:")
	testErr := errors.New(errors.CodeNotFound, "test")

	if errors.Is(testErr, errors.CodeNotFound) {
		fmt.Println("Error is a 'not found' error")
	}

	if errors.Is(testErr, errors.CodeNotFound, errors.CodeAlreadyExists) {
		fmt.Println("Error matches one of the codes")
	}

	fmt.Println("\n7. Circuit Breaker:")
	cb := errors.NewCircuitBreaker("database", 3, 10*time.Second)

	for i := 1; i <= 5; i++ {
		fmt.Printf("  Request %d: ", i)
		err := cb.Execute(func() error {
			if i <= 3 {
				return fmt.Errorf("database error")
			}
			return nil
		})

		if err != nil {
			fmt.Printf("Blocked: %v\n", err)
		} else {
			fmt.Println("Executed successfully")
		}

		time.Sleep(1 * time.Second)
	}

	fmt.Println("\n Error handling examples completed!")
}
