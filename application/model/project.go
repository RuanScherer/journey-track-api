package model

type CreateProjectRequest struct {
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
}

type CreateProjectResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
}

type EditProjectRequest struct {
	ActorID   string `json:"-"`
	ProjectID string `json:"project_id"`
	Name      string `json:"name"`
}

type EditProjectResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
}

type ShowProjectRequest struct {
	ActorID   string `json:"-"`
	ProjectID string `json:"project_id"`
}

type ShowProjectResponse struct {
	ID      string           `json:"id"`
	Name    string           `json:"name"`
	OwnerID string           `json:"owner_id"`
	Members []*ProjectMember `json:"members"`
}

type ProjectMember struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ListProjectByMemberResponse = []*ProjectByMember

type ProjectByMember struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
}

type DeleteProjectRequest struct {
	ActorID   string `json:"-"`
	ProjectID string `json:"project_id"`
}

type InviteProjectMemberRequest struct {
	ActorID   string `json:"-"`
	ProjectID string `json:"project_id"`
	UserID    string `json:"user_id"`
}

type InviteProjectMemberResponse struct {
	ID      string         `json:"id"`
	Project *InviteProject `json:"project"`
	User    *InviteUser    `json:"user"`
	Status  string         `json:"status"`
}

type InviteProject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type InviteUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type AnswerProjectInviteRequest struct {
	ProjectID   string `json:"project_id"`
	InviteToken string `json:"invite_token"`
}

type RevokeProjectInviteRequest struct {
	ActorID         string `json:"-"`
	ProjectInviteID string `json:"project_invite_id"`
}

type RegisterEventRequest struct {
	Name      string `json:"name"`
	ProjectID string `json:"project_id"`
}
