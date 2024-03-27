package usecase

import (
	"errors"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/repository"
	"gorm.io/gorm"
)

type ListProjectsByMemberUseCase struct {
	projectRepository repository.ProjectRepository
}

func NewListProjectsByMemberUseCase(
	projectRepository repository.ProjectRepository,
) *ListProjectsByMemberUseCase {
	return &ListProjectsByMemberUseCase{projectRepository}
}

func (useCase *ListProjectsByMemberUseCase) Execute(memberId string) (*appmodel.ListProjectByMemberResponse, error) {
	projects, err := useCase.projectRepository.FindByMemberId(memberId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &appmodel.ListProjectByMemberResponse{}, nil
		}
		return nil, appmodel.NewAppError("unable_to_find_projects", err.Error(), appmodel.ErrorTypeDatabase)
	}

	var projectsResponse appmodel.ListProjectByMemberResponse
	for _, p := range projects {
		project := &appmodel.ProjectByMember{
			ID:      p.ID,
			Name:    p.Name,
			OwnerID: p.OwnerID,
		}
		projectsResponse = append(projectsResponse, project)
	}

	return &projectsResponse, nil
}
