package model

type CreateProjectRequest struct {
	Name    string `json:"name" valid:"required"`
	OwnerID string `json:"owner_id" valid:"required"`
}

type CreateProjectResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
}

type EditProjectRequest struct {
	ActorID   string `json:"-" valid:"required~actor id is required"`
	ProjectID string `json:"project_id" valid:"required"`
	Name      string `json:"name" valid:"required"`
}

type EditProjectResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
}

type ShowProjectRequest struct {
	ActorID   string `json:"-" valid:"required~actor id is required"`
	ProjectID string `json:"project_id" valid:"required"`
}

type ShowProjectResponse struct {
	ID      string           `json:"id"`
	Name    string           `json:"name"`
	OwnerID string           `json:"owner_id"`
	Members []*ProjectMember `json:"members"`
}

type GetProjectStatsRequest struct {
	ActorID   string `json:"-" valid:"required~actor id is required"`
	ProjectID string `json:"project_id" valid:"required"`
}

type GetProjectStatsResponse struct {
	MembersCount int64 `json:"members_count"`
	EventsCount  int64 `json:"events_count"`
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
	ActorID   string `json:"-" valid:"required~actor id is required"`
	ProjectID string `json:"project_id" valid:"required"`
}

type InviteProjectMembersRequest struct {
	ActorID   string   `json:"-" valid:"required~actor id is required"`
	ProjectID string   `json:"project_id" valid:"required"`
	UserIDs   []string `json:"users" valid:"required"`
}

type InviteProjectMembersResponse = []*ProjectInvite

type ProjectInvite struct {
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
	ProjectID   string `json:"project_id" valid:"required"`
	InviteToken string `json:"invite_token" valid:"required"`
}

type RevokeProjectInviteRequest struct {
	ActorID         string `json:"-" valid:"required~actor id is required"`
	ProjectInviteID string `json:"project_invite_id" valid:"required"`
}

type RegisterEventRequest struct {
	Name      string `json:"name" valid:"required"`
	ProjectID string `json:"project_id" valid:"required"`
}
