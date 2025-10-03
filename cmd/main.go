// main.go
package main

import (
	"fmt"
	"net/http"
	"task/config"
	"task/internal/history"
	"task/internal/reservations"
	"task/internal/stat"
	"task/internal/user"
	"task/pkg/db"

	_ "task/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	Server()
}

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

func Server() {
	app := App()

	server := http.Server{
		Addr:    ":8081",
		Handler: app,
	}

	fmt.Println("Server is listening on port 8081")
	fmt.Println("Swagger UI: http://localhost:8081/swagger/index.html")
	fmt.Println("Swagger JSON: http://localhost:8081/swagger/doc.json")

	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
func App() http.Handler {
	conf := config.LoadConfig()
	db := db.NewDB(conf)
	router := http.NewServeMux()
	historyRepository := history.NewHistoryRepository(db)
	userRepository := user.NewUserRepository(db, historyRepository)
	reserveRepository := reservations.NewReserveRepository(db, historyRepository)
	statrepository := stat.NewStatRepository(db)

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
		httpSwagger.URL("/swagger/doc.json"), // URL для swagger.json
	))

	// Serve swagger.json directly
	router.HandleFunc("GET /swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	return router
}
