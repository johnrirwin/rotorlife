package models

import (
	"encoding/json"
	"testing"
)

// TestELRSSanitizationStripsSensitiveData verifies that bind phrase, model match,
// and other sensitive fields are NEVER present in sanitized ELRS settings.
// This is a CRITICAL security test.
func TestELRSSanitizationStripsSensitiveData(t *testing.T) {
	// Create ELRS settings with sensitive data
	fullSettings := &ELRSSettingsData{
		// SENSITIVE fields that must be stripped
		BindPhrase: "super-secret-bind-phrase",
		ModelMatch: boolPtr(true),
		ModelID:    intPtr(42),
		UID:        "rx-unique-id-12345",

		// SAFE fields that should be preserved
		ReceiverModel:    "EP1",
		PacketRate:       "500Hz",
		TelemetryRatio:   "1:64",
		SwitchMode:       "Hybrid",
		OutputPower:      "250mW",
		RegulatoryDomain: "FCC",
		FirmwareVersion:  "3.4.0",
		RXProtocol:       "CRSF",
		RFProfile:        "D500",
	}

	// Sanitize the settings
	sanitized := fullSettings.Sanitize()

	if sanitized == nil {
		t.Fatal("Sanitize() returned nil")
	}

	// CRITICAL: Verify sensitive fields are NOT present in sanitized output
	// We need to check the JSON output to ensure nothing leaks
	jsonBytes, err := json.Marshal(sanitized)
	if err != nil {
		t.Fatalf("Failed to marshal sanitized settings: %v", err)
	}
	jsonStr := string(jsonBytes)

	// These strings should NEVER appear in sanitized output
	sensitiveValues := []string{
		"super-secret-bind-phrase",
		"bindPhrase",
		"modelMatch",
		"modelId",
		"rx-unique-id-12345",
		"uid",
	}

	for _, sensitive := range sensitiveValues {
		if containsIgnoreCase(jsonStr, sensitive) {
			t.Errorf("SECURITY VIOLATION: Sanitized output contains sensitive value '%s': %s", sensitive, jsonStr)
		}
	}

	// Verify safe fields ARE present
	if sanitized.ReceiverModel != "EP1" {
		t.Errorf("Expected ReceiverModel 'EP1', got '%s'", sanitized.ReceiverModel)
	}
	if sanitized.PacketRate != "500Hz" {
		t.Errorf("Expected PacketRate '500Hz', got '%s'", sanitized.PacketRate)
	}
	if sanitized.TelemetryRatio != "1:64" {
		t.Errorf("Expected TelemetryRatio '1:64', got '%s'", sanitized.TelemetryRatio)
	}
	if sanitized.SwitchMode != "Hybrid" {
		t.Errorf("Expected SwitchMode 'Hybrid', got '%s'", sanitized.SwitchMode)
	}
	if sanitized.OutputPower != "250mW" {
		t.Errorf("Expected OutputPower '250mW', got '%s'", sanitized.OutputPower)
	}
	if sanitized.RegulatoryDomain != "FCC" {
		t.Errorf("Expected RegulatoryDomain 'FCC', got '%s'", sanitized.RegulatoryDomain)
	}
	if sanitized.FirmwareVersion != "3.4.0" {
		t.Errorf("Expected FirmwareVersion '3.4.0', got '%s'", sanitized.FirmwareVersion)
	}
}

// TestSanitizeELRSSettingsFromRawJSON tests the SanitizeELRSSettings function
// which parses raw JSON and returns sanitized settings.
func TestSanitizeELRSSettingsFromRawJSON(t *testing.T) {
	// Create a raw JSON ELRS settings object
	rawJSON := `{
		"bindPhrase": "my-secret-phrase",
		"modelMatch": true,
		"modelId": 7,
		"uid": "receiver-uid-xyz",
		"receiverModel": "RP3",
		"packetRate": "250Hz",
		"telemetryRatio": "1:128",
		"switchMode": "Wide",
		"outputPower": "100mW",
		"regulatoryDomain": "LBT",
		"firmwareVersion": "3.3.2",
		"rxProtocol": "CRSF"
	}`

	aircraftSettings := &AircraftELRSSettings{
		ID:         "test-id",
		AircraftID: "aircraft-id",
		Settings:   json.RawMessage(rawJSON),
	}

	sanitized := SanitizeELRSSettings(aircraftSettings)

	if sanitized == nil {
		t.Fatal("SanitizeELRSSettings returned nil")
	}

	// Verify JSON output doesn't contain sensitive data
	jsonBytes, err := json.Marshal(sanitized)
	if err != nil {
		t.Fatalf("Failed to marshal sanitized settings: %v", err)
	}
	jsonStr := string(jsonBytes)

	sensitiveValues := []string{
		"my-secret-phrase",
		"receiver-uid-xyz",
		"bindPhrase",
		"modelMatch",
		"modelId",
		"uid",
	}

	for _, sensitive := range sensitiveValues {
		if containsIgnoreCase(jsonStr, sensitive) {
			t.Errorf("SECURITY VIOLATION: Sanitized output contains sensitive value '%s': %s", sensitive, jsonStr)
		}
	}

	// Verify safe fields are preserved
	if sanitized.ReceiverModel != "RP3" {
		t.Errorf("Expected ReceiverModel 'RP3', got '%s'", sanitized.ReceiverModel)
	}
	if sanitized.PacketRate != "250Hz" {
		t.Errorf("Expected PacketRate '250Hz', got '%s'", sanitized.PacketRate)
	}
}

// TestSanitizeNilSettings verifies that sanitizing nil settings returns nil
func TestSanitizeNilSettings(t *testing.T) {
	var settings *ELRSSettingsData = nil
	sanitized := settings.Sanitize()
	if sanitized != nil {
		t.Error("Expected nil result for nil settings")
	}

	result := SanitizeELRSSettings(nil)
	if result != nil {
		t.Error("Expected nil result for nil AircraftELRSSettings")
	}
}

// TestSanitizeEmptySettings verifies that sanitizing empty settings works
func TestSanitizeEmptySettings(t *testing.T) {
	settings := &ELRSSettingsData{}
	sanitized := settings.Sanitize()

	if sanitized == nil {
		t.Fatal("Expected non-nil result for empty settings")
	}

	// All fields should be empty
	if sanitized.ReceiverModel != "" || sanitized.PacketRate != "" {
		t.Error("Expected empty fields for empty settings")
	}
}

// TestELRSSanitizedSettingsHasNoSensitiveFields verifies that the
// ELRSSanitizedSettings struct itself has no sensitive fields defined.
// This is a compile-time-like check through reflection.
func TestELRSSanitizedSettingsHasNoSensitiveFields(t *testing.T) {
	sanitized := &ELRSSanitizedSettings{}
	
	// Marshal to JSON and check structure
	jsonBytes, _ := json.Marshal(sanitized)
	
	// Unmarshal into a map to check field names
	var fieldMap map[string]interface{}
	json.Unmarshal(jsonBytes, &fieldMap)

	// List of fields that should NEVER exist in sanitized settings
	forbiddenFields := []string{
		"bindPhrase",
		"modelMatch",
		"modelId",
		"uid",
		"extra",
	}

	// Note: The check here is structural - the type itself doesn't have these fields
	// This test documents the expectation
	for _, field := range forbiddenFields {
		if _, exists := fieldMap[field]; exists {
			t.Errorf("SECURITY VIOLATION: ELRSSanitizedSettings should not have field '%s'", field)
		}
	}
}

// Helper functions
func boolPtr(b bool) *bool {
	return &b
}

func intPtr(i int) *int {
	return &i
}

func containsIgnoreCase(s, substr string) bool {
	return contains(toLower(s), toLower(substr))
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c = c + ('a' - 'A')
		}
		result[i] = c
	}
	return string(result)
}
