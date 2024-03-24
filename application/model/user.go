package model

type RegisterUserRequest struct {
	Email    string `json:"smtpemail" valid:"smtpemail,required"`
	Name     string `json:"name" valid:"required"`
	Password string `json:"password" valid:"required,minstringlength(8)"`
}

type RegisterUserResponse struct {
	ID         string `json:"id"`
	Email      string `json:"smtpemail"`
	Name       string `json:"name"`
	IsVerified bool   `json:"is_verified"`
}

type VerifyUserRequest struct {
	UserID            string `json:"user_id" valid:"required~user id is required"`
	VerificationToken string `json:"verification_token" valid:"required~verification token is required"`
}

type SignInRequest struct {
	Email    string `json:"smtpemail" valid:"smtpemail,required"`
	Password string `json:"password" valid:"required"`
}

type SignInResponse struct {
	AccessToken string     `json:"access_token"`
	User        SignInUser `json:"user"`
}

type SignInUser struct {
	ID    string `json:"id"`
	Email string `json:"smtpemail"`
	Name  string `json:"name"`
}

type EditUserRequest struct {
	UserID string `json:"user_id" valid:"required~user id is required"`
	Name   string `json:"name" valid:"required"`
}

type EditUserResponse struct {
	ID         string `json:"id"`
	Email      string `json:"smtpemail"`
	Name       string `json:"name"`
	IsVerified bool   `json:"is_verified"`
}

type RequestPasswordResetRequest struct {
	Email string `json:"smtpemail" valid:"smtpemail,required"`
}

type PasswordResetRequest struct {
	UserID             string `json:"user_id" valid:"required~user id is required"`
	PasswordResetToken string `json:"password_reset_token" valid:"required~password reset token is required"`
	Password           string `json:"password" valid:"required"`
}

type ShowUserResponse struct {
	ID         string `json:"id"`
	Email      string `json:"smtpemail"`
	Name       string `json:"name"`
	IsVerified bool   `json:"is_verified"`
}

type SearchUsersRequest struct {
	ActorID            string   `json:"actor_id"`
	ExcludedProjectIDs []string `json:"excluded_project_ids"`
	Email              string   `json:"smtpemail"`
	Page               int      `json:"page"`
	PageSize           int      `json:"page_size"`
}

type SearchUsersResponse = []*UserSearchResult

type UserSearchResult struct {
	ID    string `json:"id"`
	Email string `json:"smtpemail"`
	Name  string `json:"name"`
}
