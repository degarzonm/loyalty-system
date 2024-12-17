package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/degarzonm/customer_leal_service/internal/application"
	"github.com/degarzonm/customer_leal_service/internal/config"
	"github.com/degarzonm/customer_leal_service/internal/infrastructure/db"
	"github.com/degarzonm/customer_leal_service/internal/infrastructure/http"
	"github.com/degarzonm/customer_leal_service/internal/infrastructure/msg_broker"
	_ "github.com/lib/pq"
)

func main() {
	// Load initial configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// database connection and ping test
	dbConn, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Database connection established")
	// Create repositories
	customerRepo := db.NewPostgresCustomerRepo(dbConn)
	pointRepo := db.NewPostgresPointsRepo(dbConn)
	coinRepo := db.NewPostgresCoinsRepo(dbConn)
	purchaseRepo := db.NewPostgresPurchasesRepo(dbConn)
	redeemedRepo := db.NewPostgresRedeemedRepo(dbConn)

	// Create app services
	customerService := application.NewCustomerService(customerRepo)
	pointService := application.NewPointsService(pointRepo)
	coinService := application.NewCoinService(coinRepo)

	redeemService := application.NewRedeemService(redeemedRepo, pointRepo, customerRepo)

	// Kafka KafkaProducer initialization
	eventProducer, err := msgBroker.NewKafkaProducer()
	if err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %v", err)
	}
	defer eventProducer.(*msgBroker.KafkaProducer).Producer.Close()

	appService := application.NewAppService(pointRepo, customerRepo, coinRepo, eventProducer)
	purchaseService := application.NewPurchaseService(purchaseRepo, customerRepo, coinRepo, appService)

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
			log.Fatalf("Error executing Kafka listener: %v", err)
		}
	}()

	// Configure signal handling
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		log.Println("Signal received, shutting down...")
		cancel()
	}()

	// Create http handlers
	httpHandler := http.NewHandler(customerService, pointService, coinService, purchaseService, redeemService)

	// Create hhtp router
	router := http.NewRouter(httpHandler)

	// Init http server
	log.Printf("Starting HTTP server on port %s", cfg.HTTPServerPort)
	if err := router.Run(":" + cfg.HTTPServerPort); err != nil {
		log.Fatalf("Failed to run HTTP server: %v", err)
	}
}
