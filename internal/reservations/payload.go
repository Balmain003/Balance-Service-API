package reservations

// ReserveRequest представляет запрос для резервирования средств перед оплатой
type ReserveRequest struct {
	UserId    int     `json:"user_id"`    // ID пользователя
	OrderId   int     `json:"order_id"`   // ID заказа
	ServiceId int     `json:"service_id"` // ID услуги
	Amount    float64 `json:"amount"`     // Сумма которая резервируется
}

// ReserveRequest представляет ответ с данными для резервирования средств перед оплатой
type ReserveResponese struct {
	ID         int     `json:"id"`          // ID резервации
	UserId     int     `json:"user_id"`     // ID пользователя
	OrderId    int     `json:"order_id"`    // ID заказа
	ServiceId  int     `json:"service_id"`  // ID услуги
	Amount     float64 `json:"amount"`      // Сумма которая резервируется
	Status     string  `json:"status"`      // Статус
	NewBalance float64 `json:"new_balance"` // Обновленый баланс пользователя после списания средств
}
