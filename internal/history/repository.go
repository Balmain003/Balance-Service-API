package history

import (
	"fmt"
	"task/pkg/db"
)

type HistoryRepository struct {
	Database *db.Db
}

func NewHistoryRepository(database *db.Db) *HistoryRepository {
	return &HistoryRepository{
		Database: database,
	}
}

func (repo *HistoryRepository) CreateTransaction(transaction *Transaction) error {
	return repo.Database.DB.Create(transaction).Error
}

func (repo *HistoryRepository) GetUserTransactionHistory(userID, page, pageSize int, sortBy, sortOrder string) (*TransactionHistoryResponse, error) {
	var transactions []Transaction
	var totalCount int64

	query := repo.Database.DB.Where("user_id = ? OR from_user_id = ? OR to_user_id = ?",
		userID, userID, userID)

	if err := query.Model(&Transaction{}).Count(&totalCount).Error; err != nil {
		return nil, fmt.Errorf("ошибка получения общего количества транзакций: %w", err)
	}

	offset := (page - 1) * pageSize

	orderBy := getOrderBy(sortBy, sortOrder)
	query = query.Order(orderBy)

	if err := query.Offset(offset).Limit(pageSize).Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("ошибка получения транзакций: %w", err)
	}

	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize != 0 {
		totalPages++
	}

	return &TransactionHistoryResponse{
		Transactions: transactions,
		TotalCount:   int(totalCount),
		TotalPages:   totalPages,
		CurrentPage:  page,
		PageSize:     pageSize,
	}, nil
}

func getOrderBy(sortBy, sortOrder string) string {
	var orderBy string

	switch sortBy {
	case "amount":
		orderBy = "amount"
	case "date":
		fallthrough
	default:
		orderBy = "created_at"
	}

	if sortOrder == "asc" {
		orderBy += " ASC"
	} else {
		orderBy += " DESC"
	}

	return orderBy
}

func (repo *HistoryRepository) GetTransactionByID(transactionID uint) (*Transaction, error) {
	var transaction Transaction
	if err := repo.Database.DB.First(&transaction, transactionID).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}
