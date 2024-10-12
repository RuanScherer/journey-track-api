package usecase

import (
	"errors"
	"testing"

	"github.com/RuanScherer/journey-track-api/application/kafka"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/repository"
	"github.com/RuanScherer/journey-track-api/domain/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestRegisterUserUseCase_Execute(t *testing.T) {
	queueVerificationEmail = func(producerFactory kafka.ProducerFactory, user *model.User) {}

	ctrl := gomock.NewController(t)
	userRepositoryMock := repository.NewMockUserRepository(ctrl)
	producerFactoryMock := kafka.NewMockProducerFactory(ctrl)
	useCase := NewRegisterUserUseCase(userRepositoryMock, producerFactoryMock)

	req := &appmodel.RegisterUserRequest{
		Email:    "john.doe@gmail.com",
		Name:     "John Doe",
		Password: "",
	}

	res, err := useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [invalid_data_to_register_user] invalid data to register user")

	req.Password = "fake-password"
	userRepositoryMock.
		EXPECT().
		Register(gomock.Any()).
		Return(gorm.ErrDuplicatedKey)

	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [user_email_already_used] There's already an user using this email")

	userRepositoryMock.
		EXPECT().
		Register(gomock.Any()).
		Return(errors.New("unable to register user"))

	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [unable_to_register_user] unable to register user")

	userRepositoryMock.
		EXPECT().
		Register(gomock.Any()).
		AnyTimes().
		Return(nil)

	res, err = useCase.Execute(req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.ID)
	assert.Equal(t, req.Email, res.Email)
	assert.Equal(t, req.Name, res.Name)
	assert.False(t, res.IsVerified)
}

func TestRegisterUserUseCase_queueVerificationEmail(t *testing.T) {
	var queueVerificationEmail = doQueueVerificationEmail
	ctrl := gomock.NewController(t)
	producerFactoryMock := kafka.NewMockProducerFactory(ctrl)

	producerMock := kafka.NewMockProducer(ctrl)
	producerFactoryMock.
		EXPECT().
		NewProducer(gomock.Any()).
		Return(producerMock, nil)
	producerMock.
		EXPECT().
		Produce(gomock.Any(), gomock.Any()).
		Return(nil)

	user, _ := model.NewUser("john.doe@gmail.com", "John Doe", "fake-password")
	queueVerificationEmail(producerFactoryMock, user)
}
