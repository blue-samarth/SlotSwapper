package main

import (
	"os"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/joho/godotenv"

	"SlotSwapper/config"
	"SlotSwapper/models"
	"SlotSwapper/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg("No .env file found, proceeding with environment variables")
	}

	db := config.ConnectDatabase()
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	} 
	if err := db.AutoMigrate(&models.Event{}); err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	}
	log.Info().Msg("Database migrated successfully")

	route := routes.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Info().Msgf("Server starting on port %s", port)
	if err := route.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}