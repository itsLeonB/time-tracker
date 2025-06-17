package dto

type RegisterRequest struct {
	Email                string `json:"email" form:"email" binding:"required,email"`
	Password             string `json:"password" form:"password" binding:"required"`
	PasswordConfirmation string `json:"passwordConfirmation" form:"passwordConfirmation" binding:"required,eqfield=Password"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

type LoginRequest struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
}

type LoginResponse struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}
