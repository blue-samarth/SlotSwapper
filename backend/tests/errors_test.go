package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"SlotSwapper/utils"
)

func TestErrorCodes_Defined(t *testing.T) {
	// Verify all error codes are properly defined
	assert.Equal(t, "VALIDATION_ERROR", utils.ErrCodeValidation)
	assert.Equal(t, "AUTHENTICATION_ERROR", utils.ErrCodeAuthentication)
	assert.Equal(t, "AUTHORIZATION_ERROR", utils.ErrCodeAuthorization)
	assert.Equal(t, "NOT_FOUND", utils.ErrCodeNotFound)
	assert.Equal(t, "CONFLICT", utils.ErrCodeConflict)
	assert.Equal(t, "INTERNAL_SERVER_ERROR", utils.ErrCodeInternalServer)
	assert.Equal(t, "BAD_REQUEST", utils.ErrCodeBadRequest)
	assert.Equal(t, "RATE_LIMIT_EXCEEDED", utils.ErrCodeRateLimit)
	assert.Equal(t, "SWAP_LOGIC_ERROR", utils.ErrCodeSwapLogic)
}

func TestErrorResponse_Structure(t *testing.T) {
	errorResp := utils.ErrorResponse{
		Error:      utils.ErrCodeValidation,
		Message:    "Invalid input",
		Details:    map[string]interface{}{"field": "email"},
		StatusCode: 400,
		Timestamp:  "2024-01-01T00:00:00Z",
	}

	assert.Equal(t, utils.ErrCodeValidation, errorResp.Error)
	assert.Equal(t, "Invalid input", errorResp.Message)
	assert.Equal(t, 400, errorResp.StatusCode)
	assert.NotNil(t, errorResp.Details)
}

func TestSuccessResponse_Structure(t *testing.T) {
	successResp := utils.SuccessResponse{
		Message: "Operation successful",
		Data:    map[string]interface{}{"id": 1},
	}

	assert.Equal(t, "Operation successful", successResp.Message)
	assert.NotNil(t, successResp.Data)
}

func TestErrorResponse_WithDetails(t *testing.T) {
	details := map[string]interface{}{
		"field":  "email",
		"reason": "invalid format",
	}

	errorResp := utils.ErrorResponse{
		Error:      utils.ErrCodeValidation,
		Message:    "Validation failed",
		Details:    details,
		StatusCode: 400,
	}

	assert.Equal(t, "email", errorResp.Details["field"])
	assert.Equal(t, "invalid format", errorResp.Details["reason"])
}

func TestErrorResponse_WithoutDetails(t *testing.T) {
	errorResp := utils.ErrorResponse{
		Error:      utils.ErrCodeNotFound,
		Message:    "Resource not found",
		Details:    nil,
		StatusCode: 404,
	}

	assert.Nil(t, errorResp.Details)
}

func TestSuccessResponse_WithData(t *testing.T) {
	data := map[string]interface{}{
		"id":   1,
		"name": "Test User",
	}

	successResp := utils.SuccessResponse{
		Message: "User created",
		Data:    data,
	}

	dataMap := successResp.Data.(map[string]interface{})
	assert.Equal(t, 1, dataMap["id"])
	assert.Equal(t, "Test User", dataMap["name"])
}

func TestSuccessResponse_WithoutMessage(t *testing.T) {
	successResp := utils.SuccessResponse{
		Message: "",
		Data:    map[string]interface{}{"count": 5},
	}

	assert.Empty(t, successResp.Message)
	assert.NotNil(t, successResp.Data)
}

func TestErrorCodes_Uniqueness(t *testing.T) {
	// Verify all error codes are unique
	codes := []string{
		utils.ErrCodeValidation,
		utils.ErrCodeAuthentication,
		utils.ErrCodeAuthorization,
		utils.ErrCodeNotFound,
		utils.ErrCodeConflict,
		utils.ErrCodeInternalServer,
		utils.ErrCodeBadRequest,
		utils.ErrCodeRateLimit,
		utils.ErrCodeSwapLogic,
	}

	codeSet := make(map[string]bool)
	for _, code := range codes {
		assert.False(t, codeSet[code], "Error code %s should be unique", code)
		codeSet[code] = true
	}

	assert.Equal(t, 9, len(codeSet), "Should have 9 unique error codes")
}
