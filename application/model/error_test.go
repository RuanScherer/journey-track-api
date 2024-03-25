package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAppError(t *testing.T) {
	err := NewAppError("fake-code", "fake message", ErrorTypeValidation)
	assert.NotNil(t, err)
	assert.Equal(t, "fake-code", err.Code)
	assert.Equal(t, "fake message", err.Message)
	assert.Equal(t, ErrorTypeValidation, err.Type)
	assert.NotEmpty(t, err.Timestamp)
}

func TestAppError_Error(t *testing.T) {
	err := NewAppError("fake-code", "fake message", ErrorTypeValidation)
	assert.Equal(t, "(validation) [fake-code]: fake message", err.Error())
}
