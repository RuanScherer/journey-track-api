package usecase

import (
	"errors"
	"github.com/RuanScherer/journey-track-api/application/jwt"
	"github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/repository"
	domainmodel "github.com/RuanScherer/journey-track-api/domain/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestSignInUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepositoryMock := repository.NewMockUserRepository(ctrl)
	jwtManagerMock := jwt.NewMockManager(ctrl)
	useCase := NewSignInUseCase(userRepositoryMock, jwtManagerMock)

	email := "john.doe@gmail.com"
	password := "123456"
	req := &model.SignInRequest{
		Email:    email,
		Password: password,
	}

	userRepositoryMock.
		EXPECT().
		FindByEmail(email).
		Return(nil, errors.New("user not found"))

	res, err := useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [invalid_auth_credentials]: Invalid authentication credentials")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("a@bh8i32#1"), bcrypt.DefaultCost)
	user := &domainmodel.User{
		ID:         "1",
		Name:       "John Doe",
		Email:      &email,
		Password:   string(passwordHash),
		IsVerified: false,
	}
	userRepositoryMock.
		EXPECT().
		FindByEmail(email).
		AnyTimes().
		Return(user, nil)

	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [user_not_verified]: User is not verified")

	user.IsVerified = true
	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [invalid_auth_credentials]: Invalid authentication credentials")

	req.Password = "a@bh8i32#1"
	jwtManagerMock.
		EXPECT().
		CreateJwtFromUser(user).
		Return("", errors.New("unexpected error"))

	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(server) [unexpected_error]: unexpected error")

	jwtManagerMock.
		EXPECT().
		CreateJwtFromUser(user).
		Return("fake-token", nil)

	res, err = useCase.Execute(req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.User.ID, user.ID)
	assert.Equal(t, res.User.Email, *user.Email)
	assert.Equal(t, res.User.Name, user.Name)
	assert.Equal(t, res.AccessToken, "fake-token")
}
