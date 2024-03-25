package validator

import (
	"net/http"

	"github.com/RuanScherer/journey-track-api/adapters/rest/model"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/asaskevich/govalidator"
)

func ValidateRequestBody(body interface{}) error {
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		appErr := appmodel.ErrInvalidReqData
		appErr.Message = err.Error()
		return model.NewRestApiError(http.StatusBadRequest, appErr)
	}
	return nil
}
