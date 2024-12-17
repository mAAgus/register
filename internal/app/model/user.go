package model

import "time"

type Users struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	NickName          string    `json:"nick_name"`
	Email             string    `json:"email"`
	IsVerificate      bool      `json:"is_verificate"`
	CreatedAt         time.Time `json:"created_at"`
	VerificationToken string    `json:"verification_token"`
	TokenCreatedAt    time.Time `json:"token_created_at"`
}
