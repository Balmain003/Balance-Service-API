package stat

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"task/pkg/res"
)

type StatHandler struct {
	StatRepository *StatRepository
}

type StatHandlerDeps struct {
	StatRepository *StatRepository
}

func NewStatHandler(router *http.ServeMux, deps StatHandlerDeps) {
	handler := &StatHandler{
		StatRepository: deps.StatRepository,
	}
	router.HandleFunc("GET /report", handler.MonthlyReport())
	router.HandleFunc("GET /download/{filename}", handler.DownloadReport())
}

// MonthlyReport генерирует месячный отчет
// @Summary Генерация отчета
// @Description Генерирует CSV отчет по выручке за указанный месяц
// @Tags reports
// @Accept json
// @Produce json
// @Param year query int true "Год (например: 2024)"
// @Param month query int true "Месяц (1-12)"
// @Success 200 {object} ReportResponse "Ссылка на отчет"
// @Failure 400 {string} string "Некорректные параметры"
// @Router /report [get]
func (handler *StatHandler) MonthlyReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		yearStr := r.URL.Query().Get("year")
		monthStr := r.URL.Query().Get("month")

		if yearStr == "" || monthStr == "" {
			res.Json(w, ErrorResponse{Error: "Необходимо указать год и месяц"}, http.StatusBadRequest)
			return
		}
		year, err := strconv.Atoi(yearStr)
		if err != nil || year < 2020 || year > 2030 {
			res.Json(w, ErrorResponse{Error: "Некоректный год. Допустимый диапазон 2020-2030"}, http.StatusBadRequest)
		}
		month, err := strconv.Atoi(monthStr)
		if err != nil || month < 1 || month > 12 {
			res.Json(w, ErrorResponse{Error: "Некоректный месяц. Допустимый диапазон 1-12"}, http.StatusBadRequest)
			return
		}
		reportPath, err := handler.StatRepository.GenerateCSV(year, month)
		if err != nil {
			res.Json(w, ErrorResponse{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		baseURL := "http://" + r.Host
		fileName := filepath.Base(reportPath)
		reportURL := fmt.Sprintf("%s/download/%s", baseURL, fileName)

		res.Json(w, ReportResponse{
			ReportURL: reportURL,
			Message:   fmt.Sprintf("Отчет за %d-%02d успешно сгенерирован", year, month),
			FileName:  fileName,
		}, http.StatusOK)
	}
}

func (handler *StatHandler) DownloadReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		path := r.URL.Path

		var filename string

		_, err := fmt.Sscanf(path, "/download/%s", &filename)
		if err != nil || filename == "" {
			res.Json(w, ErrorResponse{Error: "Некоректный запрос"}, http.StatusBadRequest)
			return
		}

		filePath := filepath.Join("reports", filepath.Base(filename))

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			res.Json(w, ErrorResponse{Error: "Файл не найден"}, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
		w.Header().Set("Content-Description", "File Transfer")

		http.ServeFile(w, r, filePath)
	}
}
