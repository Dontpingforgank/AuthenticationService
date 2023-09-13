package Models

import "time"

type UserRegisterModel struct {
	Name     string `json:"name"`
	City     string `json:"city"`
	Country  string `json:"country"`
	Email    string `json:"email"`
	Status   string `json:"status"`
	Password string `json:"password"`
}

type UserLoginInfo struct {
	Platform     string    `json:"platform"`
	DateLoggedIn time.Time `json:"date_logged_in"`
	UserId       int       `json:"user_id"`
}
