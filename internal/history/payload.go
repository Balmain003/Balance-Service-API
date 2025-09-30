package history

// TransactionHistoryRequest представляет запрос на получении истории пользователя
type TransactionHistoryRequest struct {
	UserID    int    `json:"user_id" validate:"required,min=1"`  // ID пользователя
	Page      int    `json:"page" validate:"min=1"`              // Страница поиска
	PageSize  int    `json:"page_size" validate:"min=1,max=100"` // Размер страницы
	SortBy    string `json:"sort_by"`                            // Критерий сортировки
	SortOrder string `json:"sort_order"`                         // Критерий порядка сортировки
}

// TransactionHistoryResponese представляет ответ с историей пользователя
type TransactionHistoryResponse struct {
	Transactions []Transaction `json:"transactions"` // Слайс представляющий информацию
	TotalCount   int           `json:"total_count"`  // Общее количесво
	TotalPages   int           `json:"total_pages"`  // Общее количество страниц
	CurrentPage  int           `json:"current_page"` // Текущая страница
	PageSize     int           `json:"page_size"`
}

// ErrorResponse представляет ответ в виде ошибки
type ErrorResponse struct {
	Error string `json:"error"` // Ошибка
}
