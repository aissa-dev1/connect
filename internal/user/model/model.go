package usermodel

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type MinimalUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

type MinimalUserWithPassword struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SearchUser struct {
	Username string `json:"username"`
}
