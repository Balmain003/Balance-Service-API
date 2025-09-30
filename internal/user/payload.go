package user

// AddBalanceRequest представляет запрос на пополнение баланса
type UserCreateRequest struct {
	Name string `json:"name"` // Имя пользователя
}

// UserCreateResponese предоставляет ответ на создание пользователя
type UserCreateResponese struct {
	UserId  int     `json:"user_id"` //ID пользователя
	Name    string  `json:"name"`    // Имя пользователя
	Balance float64 `json:"balance"` // Текущий баланс
}

// AddBalanceRequest представляет запрос на пополнение баланса
type AddBalanceRequest struct {
	UserId int     `json:"user_id"` //ID пользователя
	Amount float64 `json:"amount"`  // Сумма пополнения
}

// AddBalanceResponese представляет ответ с данными пользователя
type AddBalanceResponese struct {
	UserId  int     `json:"user_id"` //ID пользователя
	Name    string  `json:"name"`    // Имя пользователя
	Balance float64 `json:"balance"` // Текущий баланс
}

// TransferRequest предоставляет запрос на перевод пользователя
type TransferRequest struct {
	FromUserId int     `json:"from_user_id"` // Имя пользователя отправителя
	ToUserId   int     `json:"to_user_id"`   // Имя пользователя получателя
	Amount     float64 `json:"amount"`       // Сумма пополнения
}

// TransferResponese предоставляет ответ о переводе пользователя
type TransferResponese struct {
	FromUserId int     `json:"from_user_id"` // Имя пользователя отправителя
	ToUserId   int     `json:"to_user_id"`   // Имя пользователя получателя
	Amount     float64 `json:"amount"`       // Сумма пополнения
	NewBalance float64 `json:"new_balance"`  // Новый баланс
}
