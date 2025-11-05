package main

import (
	"os"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/joho/godotenv"

	"SlotSwapper/config"
	"SlotSwapper/models"
	"SlotSwapper/utils"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg("No .env file found, proceeding with environment variables")
	}

	config.ConnectDatabase()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Info().Msgf("Server starting on port %s", port)

	db := config.DB
	db.AutoMigrate(&models.User{})

	token, err := utils.GenerateJWTToken(1, "alice@test.com", "alice")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to generate JWT token")
	}
	fmt.Println("Token:", token)

	claims, err := utils.ValidateJWTToken(token)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to validate JWT token")
	}
	fmt.Printf("Claims: %+v\n", claims)
}