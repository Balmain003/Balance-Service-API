package user

type IUserRepository interface {
	GetUserBalance(userId int) (*User, error)
}

type User struct {
	UserId  int     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance" gorm:"default:0"`
}

func NewUser(name string) *User {
	user := &User{
		Name:    name,
		Balance: 0,
	}
	return user
}
