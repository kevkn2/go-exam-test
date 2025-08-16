package schemas

type LoginRequestSchema struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequestSchema struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type StudentSchema struct {
	Name   string
	School string
	UserID string
}

type RegisterStudentRequestSchema struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
	School   string `json:"school" binding:"required"`
}

type AuthResponseSchema struct {
	Token string `json:"token"`
	Type  string `json:"type"`
}

type UserInfoSchema struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	Authority string `json:"authority"`
}

type StudentInfoSchema struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	Authority string `json:"authority"`
	Name      string `json:"name"`
	School    string `json:"school"`
}
