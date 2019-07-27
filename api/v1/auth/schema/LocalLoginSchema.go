package schema

type LocalLoginSchema struct {
	Email       string `json:"email" form:"email" validate:"required,email"`
	Password    string `json:"password" form:"password" validate:"required,min=6"`
}