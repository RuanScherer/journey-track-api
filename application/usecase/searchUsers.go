package usecase

import (
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/repository"
)

type SearchUsersUseCase struct {
	userRepository repository.UserRepository
}

func NewSearchUsersUseCase(userRepository repository.UserRepository) *SearchUsersUseCase {
	return &SearchUsersUseCase{userRepository}
}

func (useCase *SearchUsersUseCase) Execute(req *appmodel.SearchUsersRequest) (*appmodel.SearchUsersResponse, error) {
	users, err := useCase.userRepository.Search(repository.UserSearchOptions{
		ExcludedProjectIDs: req.ExcludedProjectIDs,
		Email:              req.Email,
		Page:               req.Page,
		PageSize:           req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	var usersResponse appmodel.SearchUsersResponse
	for _, user := range users {
		if user.ID == req.ActorID {
			continue
		}
		usersResponse = append(usersResponse, &appmodel.UserSearchResult{
			ID:    user.ID,
			Email: *user.Email,
			Name:  user.Name,
		})
	}

	return &usersResponse, nil
}
