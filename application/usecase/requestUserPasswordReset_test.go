package usecase

import (
	"errors"
	"testing"

	"github.com/RuanScherer/journey-track-api/application/factory"
	"github.com/RuanScherer/journey-track-api/application/kafka"
	"github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/repository"
	domainmodel "github.com/RuanScherer/journey-track-api/domain/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRequestUserPasswordResetUseCase_Execute(t *testing.T) {
	queuePasswordResetEmail = func(producerFactory kafka.ProducerFactory, user *domainmodel.User) {}

	ctrl := gomock.NewController(t)
	userRepositoryMock := repository.NewMockUserRepository(ctrl)
	producerFactoryMock := kafka.NewMockProducerFactory(ctrl)
	useCase := NewRequestUserPasswordResetUseCase(userRepositoryMock, producerFactoryMock)

	req := &model.RequestPasswordResetRequest{
		Email: "john.doe@gmail.com",
	}

	userRepositoryMock.
		EXPECT().
		FindByEmail(req.Email).
		Return(nil, errors.New("user not found"))

	err := useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [user_not_found] user not found")

	user, _ := factory.NewVerifiedUser("john.doe@gmail.com", "John Doe", "fake-password")
	userRepositoryMock.
		EXPECT().
		FindByEmail(req.Email).
		AnyTimes().
		Return(user, nil)
	userRepositoryMock.
		EXPECT().
		Save(user).
		Return(errors.New("unable to save"))

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(database) [unable_to_complete] unable to complete the request")

	userRepositoryMock.
		EXPECT().
		Save(user).
		AnyTimes().
		Return(nil)

	err = useCase.Execute(req)
	assert.Nil(t, err)
}

func TestRequestUserPasswordResetUseCase_sendPasswordResetEmail(t *testing.T) {
	queuePasswordResetEmail = doQueuePasswordResetEmail

	ctrl := gomock.NewController(t)
	producerFactoryMock := kafka.NewMockProducerFactory(ctrl)

	user, _ := factory.NewVerifiedUser("john.doe@gmail.com", "John Doe", "fake-password")
	user.RequestPasswordReset()

	producerMock := kafka.NewMockProducer(ctrl)
	producerFactoryMock.
		EXPECT().
		NewProducer(gomock.Any()).
		Return(producerMock, nil)
	producerMock.
		EXPECT().
		Produce(gomock.Any(), gomock.Any()).
		Return(nil)

	queuePasswordResetEmail(producerFactoryMock, user)
}
