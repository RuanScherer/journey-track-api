package usecase

import (
	"errors"
	"github.com/RuanScherer/journey-track-api/application/factory"
	"github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/repository"
	domainmodel "github.com/RuanScherer/journey-track-api/domain/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestSearchUsersUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepositoryMock := repository.NewMockUserRepository(ctrl)
	useCase := NewSearchUsersUseCase(userRepositoryMock)

	req := &model.SearchUsersRequest{
		ActorID:            "fake-actor-id",
		Email:              "doe",
		ExcludedProjectIDs: []string{"fake-project-id"},
		Page:               0,
		PageSize:           3,
	}

	userRepositoryMock.
		EXPECT().
		Search(gomock.Any()).
		Return(nil, errors.New("error searching users"))

	res, err := useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "error searching users")

	user1, _ := factory.NewVerifiedUser("john.doe@gmail.com", "John Doe", "fake-password")
	user2, _ := factory.NewVerifiedUser("jane.doe@gmail.com", "Jane Doe", "fake-password")
	user3, _ := factory.NewVerifiedUser("actor.doe@gmail.com", "Actor", "fake-password")
	req.ActorID = user3.ID
	userRepositoryMock.
		EXPECT().
		Search(gomock.Any()).
		Return([]*domainmodel.User{user1, user2, user3}, nil)

	res, err = useCase.Execute(req)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	users := *res
	assert.Equal(t, 2, len(users))
	assert.Equal(t, user1.ID, users[0].ID)
	assert.Equal(t, *user1.Email, users[0].Email)
	assert.Equal(t, user1.Name, users[0].Name)
	assert.Equal(t, user2.ID, users[1].ID)
	assert.Equal(t, *user2.Email, users[1].Email)
	assert.Equal(t, user2.Name, users[1].Name)
}
