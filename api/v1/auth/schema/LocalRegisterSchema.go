package schema

// LocalRegisterModel BodySchema 데이터
type LocalRegisterSchema struct {
	Username    string `json:"username" form:"username" validate:"required,max=16,min=1"`
	Email       string `json:"email" form:"email" validate:"required,email"`
	Password    string `json:"password" form:"password" validate:"required,min=6"`
}
