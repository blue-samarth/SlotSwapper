package utils

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

type ErrorResponse struct {
    Error      string                 `json:"error"`
    Message    string                 `json:"message,omitempty"`
    Details    map[string]interface{} `json:"details,omitempty"`
    StatusCode int                    `json:"status_code"`
    Timestamp  string                 `json:"timestamp"`
}

type SuccessResponse struct {
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
}

const (
    ErrCodeValidation      = "VALIDATION_ERROR"
    ErrCodeAuthentication  = "AUTHENTICATION_ERROR"
    ErrCodeAuthorization   = "AUTHORIZATION_ERROR"
    ErrCodeNotFound        = "NOT_FOUND"
    ErrCodeConflict        = "CONFLICT"
    ErrCodeInternalServer  = "INTERNAL_SERVER_ERROR"
    ErrCodeBadRequest      = "BAD_REQUEST"
    ErrCodeRateLimit       = "RATE_LIMIT_EXCEEDED"
    ErrCodeSwapLogic       = "SWAP_LOGIC_ERROR"
)

func RespondWithError(c *gin.Context, statusCode int, errorCode, message string, details map[string]interface{}) {
    c.JSON(statusCode, ErrorResponse{
        Error:      errorCode,
        Message:    message,
        Details:    details,
        StatusCode: statusCode,
        Timestamp:  getCurrentTimestamp(),
    })
}

func RespondWithSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
    response := SuccessResponse{
        Message: message,
        Data:    data,
    }
    c.JSON(statusCode, response)
}

func BadRequest(c *gin.Context, message string, details map[string]interface{}) {
    RespondWithError(c, http.StatusBadRequest, ErrCodeBadRequest, message, details)
}

func ValidationError(c *gin.Context, message string, details map[string]interface{}) {
    RespondWithError(c, http.StatusBadRequest, ErrCodeValidation, message, details)
}

func Unauthorized(c *gin.Context, message string) {
    RespondWithError(c, http.StatusUnauthorized, ErrCodeAuthentication, message, nil)
}

func Forbidden(c *gin.Context, message string) {
    RespondWithError(c, http.StatusForbidden, ErrCodeAuthorization, message, nil)
}

func NotFound(c *gin.Context, resource string) {
    RespondWithError(c, http.StatusNotFound, ErrCodeNotFound, resource+" not found", nil)
}

func Conflict(c *gin.Context, message string, details map[string]interface{}) {
    RespondWithError(c, http.StatusConflict, ErrCodeConflict, message, details)
}

func InternalServerError(c *gin.Context, message string) {
    RespondWithError(c, http.StatusInternalServerError, ErrCodeInternalServer, message, nil)
}

func SwapLogicError(c *gin.Context, message string, details map[string]interface{}) {
    RespondWithError(c, http.StatusBadRequest, ErrCodeSwapLogic, message, details)
}

func getCurrentTimestamp() string {
    return ""
}