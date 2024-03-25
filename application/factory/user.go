package factory

import "github.com/RuanScherer/journey-track-api/domain/model"

func NewVerifiedUser(email string, name string, password string) (*model.User, error) {
	user, err := model.NewUser(email, name, password)
	if err != nil {
		return nil, err
	}

	err = user.Verify(*user.VerificationToken)
	if err != nil {
		return nil, err
	}
	return user, nil
}
