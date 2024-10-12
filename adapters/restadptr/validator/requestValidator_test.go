package validator

import (
	"net/http"
	"testing"

	"github.com/RuanScherer/journey-track-api/adapters/restadptr/model"
	"github.com/stretchr/testify/require"
)

type Test struct {
	Name string `json:"name" valid:"required~Name is required"`
}

func TestValidateRequestBody(t *testing.T) {
	test := &Test{}
	err := ValidateRequestBody(test)
	require.NotNil(t, err)

	restErr, ok := err.(*model.RestApiError)
	require.True(t, ok)
	require.Equal(t, "invalid_request_data", restErr.Code)
	require.Equal(t, "Name is required", restErr.Message)
	require.Equal(t, http.StatusBadRequest, restErr.StatusCode)

	test = &Test{
		Name: "Ruan Scherer",
	}
	err = ValidateRequestBody(test)
	require.Nil(t, err)
}
