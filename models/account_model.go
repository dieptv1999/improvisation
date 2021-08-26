package models

type Account struct {
	Id              int    `json:"id"`
	PhotoURL        string `json:"photoURL"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	UserEmail       string `json:"user_email"`
	Username        string `json:"username"`
	UserDisplayName string `json:"user_display_name"`
	TypeLogin       int    `json:"type_login"`
	CreatedAt       int64  `json:"created_at"`
}
