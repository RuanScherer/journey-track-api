package model

type RegisterUserRequest struct {
	Email    string `json:"email" valid:"email,required"`
	Name     string `json:"name" valid:"required"`
	Password string `json:"password" valid:"required,minstringlength(8)"`
}

type RegisterUserResponse struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	IsVerified bool   `json:"is_verified"`
}

type VerifyUserRequest struct {
	UserID            string `json:"user_id" valid:"required~user id is required"`
	VerificationToken string `json:"verification_token" valid:"required~verification token is required"`
}

type SignInRequest struct {
	Email    string `json:"email" valid:"email,required"`
	Password string `json:"password" valid:"required"`
}

type SignInResponse struct {
	AccessToken string `json:"access_token"`
}

type EditUserRequest struct {
	UserID string `json:"user_id" valid:"required~user id is required"`
	Name   string `json:"name" valid:"required"`
}

type EditUserResponse struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	IsVerified bool   `json:"is_verified"`
}

type RequestPasswordResetRequest struct {
	Email string `json:"email" valid:"email,required"`
}

type PasswordResetRequest struct {
	UserID             string `json:"user_id" valid:"required~user id is required"`
	PasswordResetToken string `json:"password_reset_token" valid:"required~password reset token is required"`
	Password           string `json:"password" valid:"required"`
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
