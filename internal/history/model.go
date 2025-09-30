package history

import "time"

type Transaction struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      int       `json:"user_id" gorm:"not null; index"`
	Type        string    `json:"type" gorm:"type:varchar(50);not null" `
	Amount      float64   `json:"amount" gorm:"type decimal(15,2);not null"`
	Description string    `gorm:"type:text" json:"description"`
	RelatedID   *int      `gorm:"index" json:"related_id"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	FromUserID  *int      `gorm:"index" json:"from_user_id"`
	ToUserID    *int      `gorm:"index" json:"to_user_id"`
	ServiceID   *int      `gorm:"index" json:"service_id"`
	OrderID     *int      `gorm:"index" json:"order_id"`
}

const (
	TransactionTypeDeposit      = "deposit"      // Пополнение
	TransactionTypeWithdrawal   = "withdrawal"   // Списание
	TransactionTypeTransfer     = "transfer"     // Перевод
	TransactionTypeReservation  = "reservation"  // Резервирование
	TransactionTypeConfirmation = "confirmation" // Подтверждение
	TransactionTypeCancellation = "cancellation" // Отмена
)
