package models

import "testing"

func TestUserStatus_Values(t *testing.T) {
	// Verify the constants have expected values
	if UserStatusActive != "active" {
		t.Errorf("UserStatusActive = %q, want %q", UserStatusActive, "active")
	}
	if UserStatusDisabled != "disabled" {
		t.Errorf("UserStatusDisabled = %q, want %q", UserStatusDisabled, "disabled")
	}
	if UserStatusPending != "pending" {
		t.Errorf("UserStatusPending = %q, want %q", UserStatusPending, "pending")
	}
}

func TestAuthProvider_Values(t *testing.T) {
	if AuthProviderGoogle != "google" {
		t.Errorf("AuthProviderGoogle = %q, want %q", AuthProviderGoogle, "google")
	}
}

func TestIsValidUserStatus(t *testing.T) {
	tests := []struct {
		name   string
		status UserStatus
		want   bool
	}{
		{name: "active", status: UserStatusActive, want: true},
		{name: "disabled", status: UserStatusDisabled, want: true},
		{name: "pending", status: UserStatusPending, want: true},
		{name: "invalid", status: UserStatus("suspended"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidUserStatus(tt.status)
			if got != tt.want {
				t.Fatalf("IsValidUserStatus(%q) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestUser_Creation(t *testing.T) {
	// Verify a user can be created with basic fields
	user := User{
		ID:          "123",
		Email:       "test@example.com",
		DisplayName: "Test User",
		Status:      UserStatusActive,
	}

	if user.ID != "123" {
		t.Errorf("User.ID = %q, want %q", user.ID, "123")
	}
	if user.Email != "test@example.com" {
		t.Errorf("User.Email = %q, want %q", user.Email, "test@example.com")
	}
}
