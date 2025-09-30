package user

import (
	"fmt"
	"task/internal/history"
	"task/pkg/db"

	"gorm.io/gorm"
)

type UserRepository struct {
	Database    *db.Db
	HistoryRepo *history.HistoryRepository
}

func NewUserRepository(database *db.Db, historyRepo *history.HistoryRepository) *UserRepository {
	return &UserRepository{
		Database:    database,
		HistoryRepo: historyRepo,
	}
}

func (repo *UserRepository) AddBalance(userId int, amount float64) (*User, error) {
	var user User

	err := repo.Database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userId).First(&user, userId).Error; err != nil {
			return fmt.Errorf("пользователь не найден %w", err)
		}
		user.Balance += amount

		if err := tx.Save(&user).Error; err != nil {
			return fmt.Errorf("ошибка обновления баланса:%w", err)
		}
		transaction := &history.Transaction{
			UserID:      userId,
			Type:        history.TransactionTypeDeposit,
			Amount:      amount,
			Description: fmt.Sprintf("Пополнение баланса на сумму %.2f", amount),
		}

		if err := repo.HistoryRepo.CreateTransaction(transaction); err != nil {
			fmt.Printf("Ошибка записи транзакции: %v\n", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) CreateUser(name string) (*User, error) {
	user := User{
		Name:    name,
		Balance: 0,
	}
	result := repo.Database.DB.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo *UserRepository) GetUserBalance(userId int) (*User, error) {
	var user User

	result := repo.Database.DB.Where("user_id = ?", userId).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo *UserRepository) TransferFromUserToUser(fromUserId, toUserId int, amount float64) (*User, *User, error) {

	var fromUser, toUser User

	err := repo.Database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", fromUserId).First(&fromUser).Error; err != nil {
			return fmt.Errorf("отправитель не найден %w", err)
		}
		if err := tx.Where("user_id = ?", toUserId).First(&toUser).Error; err != nil {
			return fmt.Errorf("получатель не найден %w", err)
		}
		if fromUser == toUser {
			return fmt.Errorf("нельзя переводить самому себе")
		}
		if fromUser.Balance < amount {
			return fmt.Errorf("на счету недостаточно средств, текущий баланс %.2f", fromUser.Balance)
		}

		if err := tx.Model(&User{}).Where("user_id = ?", fromUserId).Update("balance", fromUser.Balance-amount).Error; err != nil {
			return fmt.Errorf("ошибка списания средств с баланса %w", err)
		}
		if err := tx.Model(&toUser).Where("user_id = ?", toUserId).Update("balance", toUser.Balance+amount).Error; err != nil {
			return fmt.Errorf("ошибка зачисления средств на ваш баланс %w", err)
		}
		if err := tx.Where("user_id = ?", fromUserId).First(&fromUser).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", toUserId).First(&toUser).Error; err != nil {
			return err
		}
		fromTransaction := &history.Transaction{
			UserID:      fromUserId,
			Type:        history.TransactionTypeTransfer,
			Amount:      -amount,
			Description: fmt.Sprintf("Перевод пользователю %d", toUserId),
			ToUserID:    &toUserId,
		}
		toTransaction := &history.Transaction{
			UserID:      toUserId,
			Type:        history.TransactionTypeTransfer,
			Amount:      amount,
			Description: fmt.Sprintf("Перевод от пользователя %d", fromUserId),
			FromUserID:  &fromUserId,
		}
		if err := repo.HistoryRepo.CreateTransaction(fromTransaction); err != nil {
			fmt.Printf("Ошибка записи транзакции отправителя: %v\n", err)
		}
		if err := repo.HistoryRepo.CreateTransaction(toTransaction); err != nil {
			fmt.Printf("Ошибка записи транзакции получателя: %v\n", err)
		}
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return &fromUser, &toUser, nil
}
