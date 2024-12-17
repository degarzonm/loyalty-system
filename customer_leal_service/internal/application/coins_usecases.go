package application

import (
	"github.com/degarzonm/customer_leal_service/internal/domain"
)

type coinService struct {
	coinsRepo domain.CoinsRepository
}

func NewCoinService(p domain.CoinsRepository) domain.CoinService {
	return &coinService{coinsRepo: p}
}

// GetCustomerCoins retrieves the amount of coins a customer has
func (s *coinService) GetCustomerCoins(customerID int) (int, error) {
	return s.coinsRepo.GetCoinsByCustomerID(customerID)
}

// UpdateCustomerCoins updates the amount of coins a customer has.
//
// If the amount of coins given is negative, it will subtract from the customer's current amount.
// If the amount of coins given is positive, it will add to the customer's current amount.
func (s *coinService) UpdateCustomerCoins(customerID int, coins int) error {
	return s.coinsRepo.UpdateCustomerCoins(customerID, coins)
}
