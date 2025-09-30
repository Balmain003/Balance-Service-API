package reservations

import (
	"fmt"
	"net/http"
	"task/internal/user"
	"task/pkg/req"
	"task/pkg/res"
)

type ReserveHandler struct {
	ReserveRepository *ReserveRepository
	UserRepository    user.IUserRepository
}

type ReserveHandlerDeps struct {
	ReserveRepository *ReserveRepository
	UserRepository    user.IUserRepository
}

func NewReservationsHandler(router *http.ServeMux, deps ReserveHandlerDeps) {
	handler := &ReserveHandler{
		ReserveRepository: deps.ReserveRepository,
		UserRepository:    deps.UserRepository,
	}
	router.HandleFunc("POST /reserve", handler.Reservation())
	router.HandleFunc("PATCH /reserve/confirm/{reservation_id}", handler.ConfirmReserve())
	router.HandleFunc("PATCH /reserve/cancel/{reservation_id}", handler.CancelReserve())
}

// Reservation резервирует средства
// @Summary Резервирование средств
// @Description Резервирует средства на балансе пользователя для указанной услуги
// @Tags reservations
// @Accept json
// @Produce json
// @Param input body ReserveRequest true "Данные для резервирования"
// @Success 200 {object} ReserveResponese "Результат резервирования"
// @Failure 400 {string} string "Некорректные данные для резервирования"
// @Router /reserve [post]
func (handler *ReserveHandler) Reservation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[ReserveRequest](&w, r)
		if err != nil {
			return
		}
		if body.UserId <= 0 || body.OrderId <= 0 || body.ServiceId <= 0 || body.Amount <= 0 {
			http.Error(w, "некоректные данные для резервации", http.StatusBadRequest)
			return
		}
		reservation, err := handler.ReserveRepository.Reservation(
			body.UserId,
			body.ServiceId,
			body.OrderId,
			body.Amount,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := handler.UserRepository.GetUserBalance(body.UserId)
		if err != nil {
			http.Error(w, "Ошибка получения баланса пользователя", http.StatusBadRequest)
			return
		}
		res.Json(w, "Резрвирование средств произошло успешно! Данные по резервированию", http.StatusOK)
		res.Json(w, ReserveResponese{
			ID:         reservation.ID,
			UserId:     reservation.UserId,
			OrderId:    reservation.OrderId,
			ServiceId:  reservation.ServiceId,
			Amount:     reservation.Amount,
			Status:     reservation.Status,
			NewBalance: user.Balance,
		}, http.StatusOK)
	}
}

// ConfirmReserve подтверждает резервирование
// @Summary Подтверждение резервирования
// @Description Подтверждает ранее созданное резервирование средств
// @Tags reservations
// @Produce json
// @Param reservation_id path int true "ID резервирования"
// @Success 200 {string} string "Успешное подтверждение"
// @Failure 400 {string} string "Ошибка подтверждения"
// @Router /reserve/confirm/{reservation_id} [patch]
func (handler *ReserveHandler) ConfirmReserve() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reservationId int
		path := r.URL.Path

		_, err := fmt.Sscanf(path, "/reserve/confirm/%d", &reservationId)
		if err != nil || reservationId == 0 {
			http.Error(w, "Некорректный ID резервирования", http.StatusBadRequest)
			return
		}
		err = handler.ReserveRepository.ConfirmReserve(reservationId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		answer := "Оплата прошла успешно!"
		res.Json(w, answer, http.StatusOK)
	}
}

// CancelReserve отменяет резервирование
// @Summary Отмена резервирования
// @Description Отменяет ранее созданное резервирование средств
// @Tags reservations
// @Produce json
// @Param reservation_id path int true "ID резервирования"
// @Success 200 {string} string "Успешная отмена"
// @Failure 400 {string} string "Ошибка отмены"
// @Router /reserve/cancel/{reservation_id} [patch]
func (handler *ReserveHandler) CancelReserve() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var reservationId int
		path := r.URL.Path

		_, err := fmt.Sscanf(path, "/reserve/cancel/%d", &reservationId)
		if err != nil || reservationId == 0 {
			http.Error(w, "Некорректный ID резервирования", http.StatusBadRequest)
			return
		}
		err = handler.ReserveRepository.CancelReserve(reservationId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		answer := "Оплата отменена"
		res.Json(w, answer, http.StatusOK)
	}
}
