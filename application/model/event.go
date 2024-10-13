package model

type TrackEventRequest struct {
	ProjectToken string `json:"project_token" valid:"required"`
	Name         string `json:"name" valid:"required"`
}
