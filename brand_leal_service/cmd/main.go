package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/degarzonm/brand_leal_service/internal/application"
	"github.com/degarzonm/brand_leal_service/internal/config"
	"github.com/degarzonm/brand_leal_service/internal/infrastructure/db"
	"github.com/degarzonm/brand_leal_service/internal/infrastructure/http"
	"github.com/degarzonm/brand_leal_service/internal/infrastructure/msgBroker"
	_ "github.com/lib/pq"
)

func main() {
	// Load initial configuration
	_, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}
	cfg := config.GetConfig()

	// Database connection
	dbConn, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	))
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		log.Fatalf("Error  pinging database: %v", err)
	}
	log.Println("Conection to database established")

	// Create repositories
	brandRepo := db.NewPostgresBrandRepo(dbConn)
	branchRepo := db.NewPostgresBranchRepo(dbConn)
	campaignRepo := db.NewPostgresCampaignRepo(dbConn)
	rewardRepo := db.NewPostgresRewardRepo(dbConn)

	// Create services
	brandService := application.NewBrandService(brandRepo, campaignRepo)
	branchService := application.NewBranchService(branchRepo)
	campaignService := application.NewCampaignService(campaignRepo)

	// Initialize Kafka producer
	eventProducer, err := msgBroker.NewKafkaProducer()
	if err != nil {
		log.Fatalf("Error  initializing Kafka producer: %v", err)
	}
	defer eventProducer.(*msgBroker.KafkaProducer).Producer.Close()

	// Create app service
	appService := application.NewAppService(campaignRepo, brandRepo, eventProducer)

	// Initialize Kafka listener
	kafkaListener, err := msgBroker.NewKafkaListener(appService)
	if err != nil {
		log.Fatalf("Error initializing Kafka listener: %v", err)
	}

	// Context for graceful shutdown
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Execute Kafka listener
	go func() {
		log.Println("Initializing Kafka listener...")
		if err := kafkaListener.Listen(); err != nil {
			log.Fatalf("Error en el listener de Kafka: %v", err)
		}
	}()

	// Configure graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		log.Println("Signal received, shutting down...")
		cancel()
	}()

	// Create HTTP handlers
	handler := http.NewHandler(brandService, branchService, campaignService, rewardRepo)

	// Create HTTP router
	router := http.NewRouter(handler)

	// Initialize HTTP server
	log.Printf("Initializing HTTP server on port %s", cfg.HTTPServerPort)
	if err := router.Run(":" + cfg.HTTPServerPort); err != nil {
		log.Fatalf("Error initializing HTTP server: %v", err)
	}
}
