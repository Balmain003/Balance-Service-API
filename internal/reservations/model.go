package reservations

type Reservations struct {
	ID        int     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId    int     `json:"user_id" gorm:"not null;index" `
	OrderId   int     `json:"order_id" gorm:"not null;index" `
	ServiceId int     `json:"service_id" gorm:"not null;index"`
	Amount    float64 `json:"amount" gorm:"not null"`
	Status    string  `json:"status" gorm:"type:varchar(20);default:'reserved'"`
}
