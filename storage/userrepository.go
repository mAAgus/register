package storage

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/smtp"
	"register/internal/app/model"
	"time"
)

type UserRepository struct {
	storage *Storage
}

var (
	tableUser string = "users"
)

// Create создает нового пользователя и генерирует токен подтверждения
func (ur *UserRepository) Create(u *model.Users) (*model.Users, error) {
	// Генерация токена подтверждения
	token := make([]byte, 16)
	if _, err := rand.Read(token); err != nil {
		return nil, err
	}
	verificationToken := hex.EncodeToString(token)

	// Вставка пользователя в базу данных
	query := fmt.Sprintf("INSERT INTO %s (name, nick_name, email, verification_token, token_created_at) VALUES ($1, $2, $3, $4, $5) RETURNING ID", tableUser)
	if err := ur.storage.db.QueryRow(query, u.Name, u.NickName, u.Email, verificationToken, time.Now()).Scan(&u.ID); err != nil {
		return nil, err
	}
	u.CreatedAt = time.Now()
	u.IsVerificate = false
	u.VerificationToken = verificationToken
	u.TokenCreatedAt = time.Now()

	// Отправка e-mail с подтверждением
	if err := ur.SendVerificationEmail(u.Email, verificationToken); err != nil {
		log.Printf("Failed to send verification email: %v", err)
	}

	return u, nil
}

// VerifyEmail подтверждает e-mail пользователя по токену
func (ur *UserRepository) VerifyEmail(token string) (bool, error) {
	var user model.Users
	query := fmt.Sprintf("SELECT id, email FROM %s WHERE verification_token = $1 AND is_verificate = false", tableUser)
	err := ur.storage.db.QueryRow(query, token).Scan(&user.ID, &user.Email)
	if err != nil {
		return false, err
	}

	// Обновляем статус пользователя
	updateQuery := fmt.Sprintf("UPDATE %s SET is_verificate = true, verification_token = NULL WHERE id = $1", tableUser)
	_, err = ur.storage.db.Exec(updateQuery, user.ID)
	if err != nil {
		return false, err
	}

	return true, nil
}

// sendVerificationEmail отправляет e-mail с подтверждением
func (ur *UserRepository) SendVerificationEmail(email, token string) error {
	// Настройки SMTP
	smtpHost := "smtp.your-email-provider.com"
	smtpPort := "587"
	smtpUser := "YOUR_SMTP_USERNAME"
	smtpPass := "YOUR_SMTP_PASSWORD"

	// Создание сообщения
	verificationLink := fmt.Sprintf("https://yourapp.com/verify?token=%s", token)
	subject := "Email Verification"
	body := fmt.Sprintf("Please verify your email by clicking on the following link: %s", verificationLink)
	msg := []byte("To: " + email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	// Отправка e-mail
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{email}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
func (ur *UserRepository) FindByEmail(email string) (*model.Users, bool, error) {
	users, err := ur.SelectAll()
	var founded bool
	if err != nil {
		return nil, founded, err
	}
	var userFinded *model.Users
	for _, u := range users {
		if u.Email == email {
			userFinded = u
			founded = true
			break
		}
	}
	return userFinded, founded, nil
}
func (ur *UserRepository) SelectAll() ([]*model.Users, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tableUser)
	rows, err := ur.storage.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*model.Users, 0)
	for rows.Next() {
		u := model.Users{}
		err := rows.Scan(&u.ID, &u.Name, &u.NickName, u.Email)
		if err != nil {
			log.Println(err)
			continue
		}
		users = append(users, &u)
	}
	return users, nil
}
