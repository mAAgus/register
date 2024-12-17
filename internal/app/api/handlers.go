package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"register/internal/app/model"
)

type Message struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	IsError    bool   `json:"is_error"`
}

func initHeaders(writer http.ResponseWriter) {
	writer.Header().Set("Content_Type", "application/json")
}
func (api *API) RegisterUser(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Post user register POST /api/v1/user/register")
	var user model.Users
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		api.logger.Info("Invalid json received from client")
		msg := Message{
			StatusCode: 400,
			Message:    "Provided json is invalid",
			IsError:    true,
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	// Проверка на существование пользователя
	_, ok, err := api.storage.User().FindByEmail(user.Email)
	if err != nil {
		api.logger.Info("Troubles while accessing database table(users). Error:", err)
		msg := Message{
			StatusCode: 500,
			Message:    "We have some troubles accessing the database. Try again",
			IsError:    true,
		}
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	if ok {
		api.logger.Info("User  with that email already exists")
		msg := Message{
			StatusCode: 400,
			Message:    "User  with that email already exists in the database",
			IsError:    true,
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	// Создание пользователя
	userAdded, err := api.storage.User().Create(&user)
	if err != nil {
		api.logger.Info("Troubles while accessing database table(users). Error:", err)
		msg := Message{
			StatusCode: 500,
			Message:    "We have some troubles accessing the database. Try again",
			IsError:    true,
		}
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	// Отправка e-mail с подтверждением
	if err := api.storage.User().SendVerificationEmail(userAdded.Email, userAdded.VerificationToken); err != nil {
		api.logger.Info("Failed to send verification email: ", err)
		msg := Message{
			StatusCode: 500,
			Message:    "User  registered, but failed to send verification email. Please check your email.",
			IsError:    true,
		}
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	// Успешный ответ
	msg := Message{
		StatusCode: 201,
		Message:    fmt.Sprintf("User  {Email:%s} successfully registered! Please check your email to verify your account.", userAdded.Email),
		IsError:    false,
	}
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(msg)
}

func (api *API) GetAllUsers(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Get All Users GET /api/v1/users")
	articles, err := api.storage.User().SelectAll()
	if err != nil {
		api.logger.Info("Error while Articles.SelectAll : ", err)
		msg := Message{
			StatusCode: 501,
			Message:    "We have some trobles to accessing database. Try again later",
			IsError:    true,
		}
		writer.WriteHeader(501)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode(articles)
}
