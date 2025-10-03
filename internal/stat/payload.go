// internal/stat/payload.go
package stat

// ReportRequest представляет запрос на получении статистики
type ReportRequest struct {
	Year  int `json:"year" validate:"required, min=2020,max=2030"` // Год
	Month int `json:"month" validate:"required, min=1,max=12"`     // Месяц
}

// ReportResponse представляет ответ со статистикой
type ReportResponse struct {
	ReportURL string `json:"report_url"` // Путь
	Message   string `json:"message"`    // Сообщение
	FileName  string `json:"file_name"`  // Название файла
}

// ErrorResponse представляет ответ в виде ошибки
type ErrorResponse struct {
	Error string `json:"error"` // Ошибка
}
