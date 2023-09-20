package Models

import "time"

type UserLogInModel struct {
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Date     time.Time `json:"date"`
}
