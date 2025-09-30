package history

import (
	"net/http"
	"strconv"
	"task/pkg/res"
)

type HistoryHandler struct {
	HistoryRepository *HistoryRepository
}

type HistoryHandlerDeps struct {
	HistoryRepository *HistoryRepository
}

func NewHistoryHandler(router *http.ServeMux, deps HistoryHandlerDeps) {
	handler := &HistoryHandler{
		HistoryRepository: deps.HistoryRepository,
	}

	router.HandleFunc("GET /history", handler.GetTransactionHistory())
}

// GetTransactionHistory возвращает историю транзакций
// @Summary История транзакций
// @Description Возвращает историю транзакций пользователя с пагинацией и сортировкой
// @Tags history
// @Accept json
// @Produce json
// @Param user_id query int true "ID пользователя"
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Param page_size query int false "Размер страницы (по умолчанию 20)"
// @Param sort_by query string false "Поле для сортировки (amount/date)"
// @Param sort_order query string false "Порядок сортировки (asc/desc)"
// @Success 200 {object} TransactionHistoryResponse "История транзакций"
// @Failure 400 {string} string "Некорректные параметры"
// @Router /history [get]
func (handler *HistoryHandler) GetTransactionHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userIDStr := r.URL.Query().Get("user_id")
		pageStr := r.URL.Query().Get("page")
		pageSizeStr := r.URL.Query().Get("page_size")
		sortBy := r.URL.Query().Get("sort_by")
		sortOrder := r.URL.Query().Get("sort_order")

		if userIDStr == "" {
			res.Json(w, ErrorResponse{Error: "Параметр user_id обязателен"}, http.StatusBadRequest)
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil || userID <= 0 {
			res.Json(w, ErrorResponse{Error: "Некорректный user_id"}, http.StatusBadRequest)
			return
		}

		page := 1
		if pageStr != "" {
			if page, err = strconv.Atoi(pageStr); err != nil || page < 1 {
				res.Json(w, ErrorResponse{Error: "Некорректный номер страницы"}, http.StatusBadRequest)
				return
			}
		}

		pageSize := 20
		if pageSizeStr != "" {
			if pageSize, err = strconv.Atoi(pageSizeStr); err != nil || pageSize < 1 || pageSize > 100 {
				res.Json(w, ErrorResponse{Error: "Некорректный размер страницы. Допустимый диапазон: 1-100"}, http.StatusBadRequest)
				return
			}
		}

		if sortBy == "" {
			sortBy = "date"
		}
		if sortBy != "amount" && sortBy != "date" {
			res.Json(w, ErrorResponse{Error: "Некорректный параметр сортировки. Допустимые значения: amount, date"}, http.StatusBadRequest)
			return
		}

		if sortOrder == "" {
			sortOrder = "desc"
		}
		if sortOrder != "asc" && sortOrder != "desc" {
			res.Json(w, ErrorResponse{Error: "Некорректный порядок сортировки. Допустимые значения: asc, desc"}, http.StatusBadRequest)
			return
		}

		history, err := handler.HistoryRepository.GetUserTransactionHistory(
			userID, page, pageSize, sortBy, sortOrder,
		)

		if err != nil {
			res.Json(w, ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
			return
		}

		res.Json(w, history, http.StatusOK)
	}
}
