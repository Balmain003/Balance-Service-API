package reservations

import (
	"fmt"
	"task/internal/history"
	"task/internal/user"
	"task/pkg/db"

	"gorm.io/gorm"
)

type ReserveRepository struct {
	Database    *db.Db
	HistoryRepo *history.HistoryRepository
}

func NewReserveRepository(database *db.Db, historyRepo *history.HistoryRepository) *ReserveRepository {
	return &ReserveRepository{
		Database:    database,
		HistoryRepo: historyRepo,
	}
}

func (repo *ReserveRepository) Reservation(userId, serviceId, orderId int, amount float64) (*Reservations, error) {
	var reserve Reservations
	err := repo.Database.DB.Transaction(func(tx *gorm.DB) error {

		var user user.User
		if err := tx.Where("user_id = ?", userId).First(&user).Error; err != nil {
			return fmt.Errorf("пользователь не найден %w", err)
		}
		if user.Balance < amount {
			return fmt.Errorf("на счету недостаточно средств, текущий баланс %.2f", user.Balance)
		}
		var existingReservation Reservations
		if err := tx.Where("user_id = ? AND service_id = ? AND order_id = ?",
			userId, serviceId, orderId).First(&existingReservation).Error; err == nil {
			return fmt.Errorf("резервация для этого товара или услуги уже существует")
		}
		if err := tx.Model(&user).Where("user_id = ?", userId).Update("balance", user.Balance-amount).Error; err != nil {
			return fmt.Errorf("ошибка списания средств %w", err)
		}
		reserve = Reservations{
			UserId:    userId,
			OrderId:   orderId,
			ServiceId: serviceId,
			Amount:    amount,
			Status:    "reserved",
		}
		if err := tx.Create(&reserve).Error; err != nil {
			return fmt.Errorf("ошибка резервирования средств %w", err)
		}
		transaction := &history.Transaction{
			UserID:      userId,
			Type:        history.TransactionTypeReservation,
			Amount:      -amount,
			Description: fmt.Sprintf("Резервирование для услуги %d, заказ %d", serviceId, orderId),
			ServiceID:   &serviceId,
			OrderID:     &orderId,
			RelatedID:   &reserve.ID,
		}

		if err := repo.HistoryRepo.CreateTransaction(transaction); err != nil {
			fmt.Printf("Ошибка записи транзакции резервирования: %v\n", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return &reserve, nil
}

func (repo *ReserveRepository) GetReservation(reservationId uint) (*Reservations, error) {
	var reservation Reservations
	if err := repo.Database.DB.First(&reservation, reservationId).Error; err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (repo *ReserveRepository) ConfirmReserve(resevationId int) error {
	return repo.Database.DB.Transaction(func(tx *gorm.DB) error {
		var reservation Reservations

		if err := tx.First(&reservation, resevationId).Error; err != nil {
			return fmt.Errorf("резервирование не найдено: %w", err)
		}
		if reservation.Status != "reserved" {
			return fmt.Errorf("невозможно подтвердить оплату со статусом %s", reservation.Status)
		}
		if err := tx.Model(&reservation).Update("status", "confirmed").Error; err != nil {
			return err
		}
		transaction := &history.Transaction{
			UserID:      reservation.UserId,
			Type:        history.TransactionTypeConfirmation,
			Amount:      0,
			Description: fmt.Sprintf("Подтверждение оплаты услуги %d, заказ %d", reservation.ServiceId, reservation.OrderId),
			ServiceID:   &reservation.ServiceId,
			OrderID:     &reservation.OrderId,
			RelatedID:   &reservation.ID,
		}
		if err := repo.HistoryRepo.CreateTransaction(transaction); err != nil {
			fmt.Printf("Ошибка записи транзакции подтверждения: %v\n", err)
		}
		return nil
	})
}

func (repo *ReserveRepository) CancelReserve(reservationId int) error {
	return repo.Database.DB.Transaction(func(tx *gorm.DB) error {
		var reservation Reservations

		if err := tx.First(&reservation, reservationId).Error; err != nil {
			return fmt.Errorf("резервирование не найдено: %w", err)
		}
		if reservation.Status == "confirmed" {
			return fmt.Errorf("невозможно отменить оплату, так как она уже подтверждена")
		}
		if reservation.Status == "canceled" {
			return fmt.Errorf("невозможно отменить оплату, так как она уже отменена")
		}
		if err := tx.Model(&reservation).Update("status", "canceled").Error; err != nil {
			return err
		}
		if err := tx.Model(&user.User{}).Where("user_id = ?", reservation.UserId).Update("balance", gorm.Expr("balance + ?", reservation.Amount)).Error; err != nil {
			return err
		}
		transaction := &history.Transaction{
			UserID:      reservation.UserId,
			Type:        history.TransactionTypeCancellation,
			Amount:      reservation.Amount,
			Description: fmt.Sprintf("Отмена резервирования для услуги %d, заказ %d", reservation.ServiceId, reservation.OrderId),
			ServiceID:   &reservation.ServiceId,
			OrderID:     &reservation.OrderId,
			RelatedID:   &reservation.ID,
		}

		if err := repo.HistoryRepo.CreateTransaction(transaction); err != nil {
			fmt.Printf("Ошибка записи транзакции отмены: %v\n", err)
		}
		return nil
	})
}
