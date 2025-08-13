package schemas

type LoginRequestSchema struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequestSchema struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterStudentRequestSchema struct {
	RegisterRequestSchema
	School string `json:"school" binding:"required"`
}

type AuthResponseSchema struct {
	Token string `json:"token"`
	Type  string `json:"type"`
}
