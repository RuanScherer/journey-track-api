package model

type RegisterUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type RegisterUserResponse struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	IsVerified bool   `json:"is_verified"`
}

type VerifyUserRequest struct {
	UserID            string `json:"user_id"`
	VerificationToken string `json:"verification_token"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	AccessToken string `json:"access_token"`
}

type EditUserRequest struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

type EditUserResponse struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	IsVerified bool   `json:"is_verified"`
}

type RequestPasswordResetRequest struct {
	Email string `json:"email"`
}

type PasswordResetRequest struct {
	UserID             string `json:"user_id"`
	PasswordResetToken string `json:"password_reset_token"`
	Password           string `json:"password"`
}

type ShowUserResponse struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	IsVerified bool   `json:"is_verified"`
}

type SearchUserResponse = []*UserSearchResult

type UserSearchResult struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	IsVerified bool   `json:"is_verified"`
}
