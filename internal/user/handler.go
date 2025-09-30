package user

import (
	"fmt"
	"net/http"
	"task/pkg/req"
	"task/pkg/res"

	"gorm.io/gorm"
)

type UserHandler struct {
	UserRepository *UserRepository
}

type UserHandlerDeps struct {
	UserRepository *UserRepository
}

func NewUserHandler(router *http.ServeMux, deps UserHandlerDeps) {
	handler := &UserHandler{
		UserRepository: deps.UserRepository,
	}
	router.HandleFunc("POST /user", handler.CreateUser())
	router.HandleFunc("PATCH /balance", handler.AddBalance())
	router.HandleFunc("GET /balance/{user_id}", handler.GetUserBalance())
	router.HandleFunc("POST /transfer", handler.Transfer())
}

// AddBalance пополняет баланс пользователя
// @Summary Пополнение баланса
// @Description Зачисление средств на баланс пользователя
// @Tags balance
// @Accept json
// @Produce json
// @Param input body AddBalanceRequest true "Данные для пополнения баланса"
// @Success 200 {object} User "Данные пользователя с обновленным балансом"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /balance [patch]
func (handler *UserHandler) AddBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[AddBalanceRequest](&w, r)
		if err != nil {
			return
		}
		if body.Amount <= 0 || body.UserId <= 0 {
			http.Error(w, "Некорректный запрос", http.StatusBadRequest)
			return
		}
		user, err := handler.UserRepository.AddBalance(body.UserId, body.Amount)
		if err != nil {
			http.Error(w, "Ошибка пополнения баланса", http.StatusBadRequest)
			return
		}
		res.Json(w, user, 200)
	}
}

func (handler *UserHandler) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[UserCreateRequest](&w, r)
		if err != nil {
			return
		}
		user := User{
			Name:    body.Name,
			Balance: 0,
		}
		if err := handler.UserRepository.Database.DB.Create(&user).Error; err != nil {
			http.Error(w, "Ошибка создания пользователя", http.StatusBadRequest)
			return
		}
		res.Json(w, user, 201)
	}
}

// GetUserBalance возвращает баланс пользователя
// @Summary Получение баланса
// @Description Возвращает информацию о балансе пользователя
// @Tags balance
// @Accept json
// @Produce json
// @Param user_id path int true "ID пользователя"
// @Success 200 {object} User "Данные пользователя"
// @Failure 400 {string} string "Некорректный ID пользователя"
// @Failure 404 {string} string "Пользователь не найден"
// @Router /balance/{user_id} [get]
func (handler *UserHandler) GetUserBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		var userId int
		_, err := fmt.Sscanf(path, "/balance/%d", &userId)
		if err != nil || userId <= 0 {
			http.Error(w, "Некорректный ID пользователя", http.StatusBadRequest)
			return
		}
		user, err := handler.UserRepository.GetUserBalance(userId)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, "Пользователь не найден", http.StatusNotFound)
			} else {
				http.Error(w, "Ошибка получения баланса", http.StatusInternalServerError)
			}
			return
		}

		res.Json(w, user, http.StatusOK)
	}
}

// Transfer осуществляет перевод между пользователями
// @Summary Перевод средств
// @Description Переводит указанную сумму от одного пользователя другому
// @Tags transfer
// @Accept json
// @Produce json
// @Param input body TransferRequset true "Данные для перевода"
// @Success 200 {object} map[string]interface{} "Результат перевода"
// @Failure 400 {string} string "Некорректные данные для перевода"
// @Router /transfer [post]
func (handler *UserHandler) Transfer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[TransferRequest](&w, r)
		if err != nil {
			return
		}

		if body.FromUserId <= 0 || body.ToUserId <= 0 {
			http.Error(w, "Некоректные данные для перевода", http.StatusBadRequest)
			return
		}

		fromUser, toUser, err := handler.UserRepository.TransferFromUserToUser(body.FromUserId, body.ToUserId, body.Amount)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, map[string]interface{}{
			"message": "Перевод выполнен успешно",
			"from_user": map[string]interface{}{
				"user_id":     fromUser.UserId,
				"name":        fromUser.Name,
				"new_balance": fromUser.Balance,
			},
			"to_user": map[string]interface{}{
				"user_id":     toUser.UserId,
				"name":        toUser.Name,
				"new_balance": toUser.Balance,
			},
			"amount": body.Amount,
		}, http.StatusOK)
	}
}
