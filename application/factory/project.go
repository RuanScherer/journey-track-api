package factory

import "github.com/RuanScherer/journey-track-api/domain/model"

func NewProjectWithDefaultOwner(projectName string) (*model.Project, error) {
	owner, err := NewVerifiedUser("john.doe@gmail.com", "John Doe", "fake-password")
	if err != nil {
		return nil, err
	}

	project, err := model.NewProject(projectName, owner)
	if err != nil {
		return nil, err
	}
	return project, nil
}
