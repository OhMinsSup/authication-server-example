package schema

type EmailVerificationSchema struct {
	Email string `json:"email" form:"email" query:"email" validate:"required,email"`
}