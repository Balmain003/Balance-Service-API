// internal/user/payload.go
package user

// AddBalanceRequest представляет запрос на пополнение баланса
// @Description Запрос для создания нового пользователя
type UserCreateRequest struct {
	Name string `json:"name"` // Имя пользователя
}

// UserCreateResponese предоставляет ответ на создание пользователя
// @Description Информация о балансе пользователя
type UserCreateResponse struct {
	UserId  int     `json:"user_id"` //ID пользователя
	Name    string  `json:"name"`    // Имя пользователя
	Balance float64 `json:"balance"` // Текущий баланс
}

// AddBalanceRequest представляет запрос на пополнение баланса
// @Description Запрос для пополнения баланса пользователя
type AddBalanceRequest struct {
	UserId int     `json:"user_id"` //ID пользователя
	Amount float64 `json:"amount"`  // Сумма пополнения
}

// AddBalanceResponese представляет ответ с данными пользователя
// @Description Информация о балансе пользователя
type AddBalanceResponse struct {
	UserId  int     `json:"user_id"` //ID пользователя
	Name    string  `json:"name"`    // Имя пользователя
	Balance float64 `json:"balance"` // Текущий баланс
}

// TransferRequest предоставляет запрос на перевод пользователя
// @Description Запрос для перевода средств между пользователями
type TransferRequest struct {
	FromUserId int     `json:"from_user_id"` // Имя пользователя отправителя
	ToUserId   int     `json:"to_user_id"`   // Имя пользователя получателя
	Amount     float64 `json:"amount"`       // Сумма пополнения
}

// TransferResponese предоставляет ответ о переводе пользователя
// @Description Ответ с результатом перевода средств
type TransferResponse struct {
	FromUserId int     `json:"from_user_id"` // Имя пользователя отправителя
	ToUserId   int     `json:"to_user_id"`   // Имя пользователя получателя
	Amount     float64 `json:"amount"`       // Сумма пополнения
	NewBalance float64 `json:"new_balance"`  // Новый баланс
}

// ErrorResponse представляет ответ об ошибке
// @Description Стандартный ответ при возникновении ошибки
type ErrorResponse struct {
	Error string `json:"error" example:"Описание ошибки"`
}
