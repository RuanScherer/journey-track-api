package model

type UserVerificationEmailConfig struct {
	UserName         string
	VerificationLink string
}

type UserPasswordResetEmailConfig struct {
	UserName          string
	PasswordResetLink string
}

type ProjectInviteEmailConfig struct {
	UserName         string
	IssuerName       string
	ProjectName      string
	AnswerInviteLink string
}
