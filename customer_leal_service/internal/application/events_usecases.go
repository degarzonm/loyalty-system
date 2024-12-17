package application

import (
	"log"

	"github.com/degarzonm/customer_leal_service/internal/config"
	"github.com/degarzonm/customer_leal_service/internal/domain"
)

type AppService struct {
	pointRepo     domain.PointsRepository
	customerRepo  domain.CustomerRepository
	coinRepo      domain.CoinsRepository
	eventProducer domain.EventProducer
}

// NewAppService creates a new application service
func NewAppService(pointRepo domain.PointsRepository, customerRepo domain.CustomerRepository, coinRepo domain.CoinsRepository, producer domain.EventProducer) *AppService {
	return &AppService{
		pointRepo:     pointRepo,
		customerRepo:  customerRepo,
		coinRepo:      coinRepo,
		eventProducer: producer,
	}
}

// ProcessApplyPointsEvent processes a new points application event, by recording the transaction,
// updating the customer points and updating the customer coins.
func (s *AppService) ProcessApplyPointsEvent(pointsEvent domain.LealPointsApply) error {
	log.Println("Service: ProcessApplyPointsEvent, with points: ", pointsEvent)
	//Record transaction
	err := s.pointRepo.RecordPointsTransaction(&domain.LealPointsTransaction{CustomerID: pointsEvent.CustomerID,
		BrandID: pointsEvent.BrandID, Change: pointsEvent.Points, Reason: pointsEvent.Reason})
	if err != nil {
		return err
	}

	//Update points
	err = s.pointRepo.UpdatePoints(pointsEvent.CustomerID, pointsEvent.BrandID, pointsEvent.Points)
	if err != nil {
		return err
	}

	//update coins
	err = s.coinRepo.UpdateCustomerCoins(pointsEvent.CustomerID, pointsEvent.Coins)
	if err != nil {
		return err
	}
	return nil
}

// SendPurchaseEvent sends a purchase message to a configured Kafka topic.
// The method retrieves the topic name from the global configuration and sends
// the raw purchase object without marshaling it. It returns an error if the
// message could not be sent.

func (s *AppService) SendPurchaseEvent(purchase domain.Purchase) error {

	cfg := config.GetConfig()
	return s.eventProducer.SendMessage(cfg.MsgPurchaseTopic, purchase)
}
