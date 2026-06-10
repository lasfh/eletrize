package main

import (
	"os"
	"testing"
)

func TestLock(t *testing.T) {
	tests := []struct {
		name           string
		setupEnv       func()
		cleanupEnv     func()
		expectedError  error
		expectedLocked string
	}{
		{
			name: "should set ELETRIZE_LOCKED when not set",
			setupEnv: func() {
				os.Unsetenv(envLocked)
				os.Unsetenv(envSub)
			},
			cleanupEnv: func() {
				os.Unsetenv(envLocked)
				os.Unsetenv(envSub)
			},
			expectedError:  nil,
			expectedLocked: "1",
		},
		{
			name: "should return nil when ELETRIZE_SUB is 1",
			setupEnv: func() {
				os.Setenv(envLocked, "1")
				os.Setenv(envSub, "1")
			},
			cleanupEnv: func() {
				os.Unsetenv(envLocked)
				os.Unsetenv(envSub)
			},
			expectedError:  nil,
			expectedLocked: "1",
		},
		{
			name: "should return error when already locked and ELETRIZE_SUB is not 1",
			setupEnv: func() {
				os.Setenv(envLocked, "1")
				os.Unsetenv(envSub)
			},
			cleanupEnv: func() {
				os.Unsetenv(envLocked)
				os.Unsetenv(envSub)
			},
			expectedError:  ErrAlreadyLocked,
			expectedLocked: "1",
		},
		{
			name: "should return error when already locked and ELETRIZE_SUB is empty",
			setupEnv: func() {
				os.Setenv(envLocked, "1")
				os.Setenv(envSub, "")
			},
			cleanupEnv: func() {
				os.Unsetenv(envLocked)
				os.Unsetenv(envSub)
			},
			expectedError:  ErrAlreadyLocked,
			expectedLocked: "1",
		},
		{
			name: "should return error when already locked and ELETRIZE_SUB is different value",
			setupEnv: func() {
				os.Setenv(envLocked, "1")
				os.Setenv(envSub, "0")
			},
			cleanupEnv: func() {
				os.Unsetenv(envLocked)
				os.Unsetenv(envSub)
			},
			expectedError:  ErrAlreadyLocked,
			expectedLocked: "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment
			tt.setupEnv()
			defer tt.cleanupEnv()

			err := lock()
			if tt.expectedError == nil {
				if err != nil {
					t.Errorf("lock() returned error: %v, expected nil", err)
				}
			} else {
				if err != tt.expectedError {
					t.Errorf("lock() returned error: %v, expected: %v", err, tt.expectedError)
				}
			}

			// Check if ELETRIZE_LOCKED was set correctly (only for first test case)
			if tt.expectedLocked != "" {
				locked := os.Getenv(envLocked)
				if locked != tt.expectedLocked {
					t.Errorf("ELETRIZE_LOCKED = %q, expected %q", locked, tt.expectedLocked)
				}
			}
		})
	}
}
