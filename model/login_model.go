package model

type Login struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}
