package domain

type User struct {
	Username     string `json:"username"`
	PasswordHash []byte `json:"-"`
}

type UserCreate struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
