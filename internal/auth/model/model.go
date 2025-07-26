package authmodel

type SignUpRequestBody struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
