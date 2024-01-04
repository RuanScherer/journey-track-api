package model

import (
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
)

type RestApiError struct {
	appmodel.AppError
	StatusCode int `json:"-"`
}

func NewRestApiError(statusCode int, appError *appmodel.AppError) *RestApiError {
	return &RestApiError{
		StatusCode: statusCode,
		AppError:   *appError,
	}
}
