package application

import (
	"errors"
	"log"

	"github.com/degarzonm/customer_leal_service/internal/domain"
)

type purchaseService struct {
	purchaseRepo domain.PurchasesRepository
	customerRepo domain.CustomerRepository
	coinRepo     domain.CoinsRepository
	appService   *AppService // Referencia a AppService
}

func NewPurchaseService(p domain.PurchasesRepository, c domain.CustomerRepository, cr domain.CoinsRepository, app *AppService) domain.PurchaseService {
	return &purchaseService{purchaseRepo: p, customerRepo: c, coinRepo: cr, appService: app}
}

// ProcessPurchase process a purchase, verifing if the customer has enough coins to do the purchase,
// if yes, it updates the coins, records the purchase and send a purchase event to kafka
func (s *purchaseService) ProcessPurchase(attempPurchase *domain.Purchase) (*domain.Purchase, error) {
	// Verify if coins are enough
	log.Println("Processing purchase: ", attempPurchase)
	customer, err := s.customerRepo.GetCustomerByID(attempPurchase.CustomerID)
	if err != nil {
		return nil, err
	}
	if customer.LealCoins < attempPurchase.CoinsUsed {
		return nil, errors.New("not enough coins")
	}

	// Update coins
	log.Println("Update coins: ", attempPurchase.CoinsUsed, " for customer: ", attempPurchase.CustomerID)
	err = s.coinRepo.UpdateCustomerCoins(attempPurchase.CustomerID, -attempPurchase.CoinsUsed)
	if err != nil {
		return nil, err
	}

	// Record purchase
	attempPurchase, err = s.purchaseRepo.RecordPurchase(attempPurchase)
	if err != nil {
		return nil, err
	}

	// Call SendPurchaseEvent through AppService
	if err := s.appService.SendPurchaseEvent(*attempPurchase); err != nil {
		return nil, errors.New("failed to send apply points event")
	}

	return attempPurchase, nil
}
