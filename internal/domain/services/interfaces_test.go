package services

import (
	"testing"

	"go.uber.org/mock/gomock"
)

// TestExecutionServiceInterface verifies the ExecutionService interface exists and has correct methods
func TestExecutionServiceInterface(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// This test ensures the ExecutionService interface compiles correctly
	var _ ExecutionService = NewMockExecutionService(ctrl)
}

// TestStatusServiceInterface verifies the StatusService interface exists and has correct methods
func TestStatusServiceInterface(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// This test ensures the StatusService interface compiles correctly
	var _ StatusService = NewMockStatusService(ctrl)
}

// TestConfigServiceInterface verifies the ConfigService interface exists and has correct methods
func TestConfigServiceInterface(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// This test ensures the ConfigService interface compiles correctly
	var _ ConfigService = NewMockConfigService(ctrl)
}

// TestValidationServiceInterface verifies the ValidationService interface exists and has correct methods
func TestValidationServiceInterface(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// This test ensures the ValidationService interface compiles correctly
	var _ ValidationService = NewMockValidationService(ctrl)
}

// TestLoggingServiceInterface verifies the LoggingService interface exists and has correct methods
func TestLoggingServiceInterface(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// This test ensures the LoggingService interface compiles correctly
	var _ LoggingService = NewMockLoggingService(ctrl)
}
