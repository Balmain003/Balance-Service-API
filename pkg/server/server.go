package server

import (
	"net/http"
	"task/config"
	"task/internal/history"
	"task/internal/reservations"
	"task/internal/stat"
	"task/internal/user"
	"task/pkg/db"

	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
)

// @title Balance Service API
// @version 1.0
// @description Микросервис для работы с балансом пользователей
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@balance-service.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8081
// @BasePath /
func AppWithDB(dbg *gorm.DB) http.Handler {
	router := http.NewServeMux()

	dbWrapper := &db.Db{DB: dbg}

	historyRepository := history.NewHistoryRepository(dbWrapper)
	userRepository := user.NewUserRepository(dbWrapper, historyRepository)
	reserveRepository := reservations.NewReserveRepository(dbWrapper, historyRepository)
	statrepository := stat.NewStatRepository(dbWrapper)

	user.NewUserHandler(router, user.UserHandlerDeps{
		UserRepository: userRepository,
	})
	reservations.NewReservationsHandler(router, reservations.ReserveHandlerDeps{
		ReserveRepository: reserveRepository,
		UserRepository:    userRepository,
	})
	stat.NewStatHandler(router, stat.StatHandlerDeps{
		StatRepository: statrepository,
	})
	history.NewHistoryHandler(router, history.HistoryHandlerDeps{
		HistoryRepository: historyRepository,
	})

	router.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
		httpSwagger.UIConfig(map[string]string{
			"defaultModelsExpandDepth": "3",
			"displayRequestDuration":   "true",
		}),
		httpSwagger.URL("/swagger/doc.json"),
	))

	router.HandleFunc("GET /swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	return router
}

func App() http.Handler {
	conf := config.LoadConfig()
	db := db.NewDB(conf)
	return AppWithDB(db.DB)
}
