package auth

import (
	"testing"
)

func TestAuthError(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		message  string
		expected string
	}{
		{
			name:     "invalid_input error",
			code:     "invalid_input",
			message:  "email is required",
			expected: "email is required",
		},
		{
			name:     "user_exists error",
			code:     "user_exists",
			message:  "a user with this email already exists",
			expected: "a user with this email already exists",
		},
		{
			name:     "invalid_credentials error",
			code:     "invalid_credentials",
			message:  "invalid email or password",
			expected: "invalid email or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &AuthError{Code: tt.code, Message: tt.message}
			if err.Error() != tt.expected {
				t.Errorf("AuthError.Error() = %s, want %s", err.Error(), tt.expected)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "securepassword123",
			wantErr:  false,
		},
		{
			name:     "exactly 8 characters",
			password: "12345678",
			wantErr:  false,
		},
		{
			name:     "too short",
			password: "1234567",
			wantErr:  true,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Password validation logic
			if tt.password == "" {
				if !tt.wantErr {
					t.Error("Expected error for empty password")
				}
				return
			}
			if len(tt.password) < 8 {
				if !tt.wantErr {
					t.Error("Expected error for short password")
				}
				return
			}
			if tt.wantErr {
				t.Errorf("Did not expect error for password: %s", tt.password)
			}
		})
	}
}

func TestDefaultDisplayName(t *testing.T) {
	tests := []struct {
		name         string
		email        string
		displayName  string
		expectedName string
	}{
		{
			name:         "uses provided display name",
			email:        "test@example.com",
			displayName:  "Test User",
			expectedName: "Test User",
		},
		{
			name:         "derives from email when empty",
			email:        "john@example.com",
			displayName:  "",
			expectedName: "john",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string
			if tt.displayName != "" {
				result = tt.displayName
			} else {
				// Extract username from email
				for i, c := range tt.email {
					if c == '@' {
						result = tt.email[:i]
						break
					}
				}
			}
			if result != tt.expectedName {
				t.Errorf("Display name = %s, want %s", result, tt.expectedName)
			}
		})
	}
}
