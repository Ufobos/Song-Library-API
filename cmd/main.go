package main

import (
	"fmt"
	"net/http"
	"song-library/configs"
	"song-library/internal/handlers"
	"song-library/internal/repository"
	"song-library/internal/services"
	"song-library/internal/utils"

	_ "song-library/docs"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// @title Song Library API
// @version 1.0
// @description API Server for Song Library
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	config := configs.LoadConfig()
	utils.InitLogger()
	defer utils.Logger.Sync()

	// Connect to the database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPass, config.DBName)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		utils.Logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Run migrations
	err = repository.Migrate(db)
	if err != nil {
		utils.Logger.Fatal("Migration failed", zap.Error(err))
	}

	// Initialize repository, service, and handler
	songRepo := repository.NewSongRepository(db)
	songService := services.NewSongService(songRepo, config)
	songHandler := handlers.NewSongHandler(songService)

	// Set up routes
	router := mux.NewRouter()

	router.HandleFunc("/songs", songHandler.AddSong).Methods("POST")
	router.HandleFunc("/songs", songHandler.GetSongs).Methods("GET")
	router.HandleFunc("/songs/{id}/text", songHandler.GetSongText).Methods("GET")
	router.HandleFunc("/songs/{id}", songHandler.UpdateSong).Methods("PUT")
	router.HandleFunc("/songs/{id}", songHandler.DeleteSong).Methods("DELETE")

	// Swagger
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	utils.Logger.Info("Server is running on port 8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		utils.Logger.Fatal("Failed to start server", zap.Error(err))
	}
}
