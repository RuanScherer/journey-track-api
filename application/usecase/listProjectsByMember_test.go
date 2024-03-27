package usecase

import (
	"errors"
	"github.com/RuanScherer/journey-track-api/application/factory"
	"github.com/RuanScherer/journey-track-api/application/repository"
	"github.com/RuanScherer/journey-track-api/domain/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"testing"
)

func TestListProjectsByMemberUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	projectRepositoryMock := repository.NewMockProjectRepository(ctrl)
	useCase := NewListProjectsByMemberUseCase(projectRepositoryMock)
	memberId := "fake-member-id"

	projectRepositoryMock.
		EXPECT().
		FindByMemberId(memberId).
		Return(nil, gorm.ErrRecordNotFound)

	res, err := useCase.Execute(memberId)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Empty(t, *res)

	projectRepositoryMock.
		EXPECT().
		FindByMemberId(memberId).
		Return(nil, errors.New("unexpected error"))

	res, err = useCase.Execute(memberId)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [unable_to_find_projects] unexpected error")

	project1, _ := factory.NewProjectWithDefaultOwner("fake project 1")
	project2, _ := factory.NewProjectWithDefaultOwner("fake project 2")
	projectRepositoryMock.
		EXPECT().
		FindByMemberId(memberId).
		Return([]*model.Project{project1, project2}, nil)

	res, err = useCase.Execute(memberId)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Len(t, *res, 2)
	assert.Equal(t, (*res)[0].ID, project1.ID)
	assert.Equal(t, (*res)[0].Name, project1.Name)
	assert.Equal(t, (*res)[0].OwnerID, project1.OwnerID)
	assert.Equal(t, (*res)[1].ID, project2.ID)
	assert.Equal(t, (*res)[1].Name, project2.Name)
	assert.Equal(t, (*res)[1].OwnerID, project2.OwnerID)
}
