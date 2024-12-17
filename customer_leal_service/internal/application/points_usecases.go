package application

import (
	"github.com/degarzonm/customer_leal_service/internal/domain"
)

type pointService struct {
	pointRepo domain.PointsRepository
}

func NewPointsService(p domain.PointsRepository) domain.PointService {
	return &pointService{pointRepo: p}
}

// GetCustomerPoints fetches all points for a given customer ID
// from the PointsRepository, returning an array of LealPoints
// and an error if any.
func (s *pointService) GetCustomerPoints(customerID int) ([]domain.LealPoints, error) {
	return s.pointRepo.GetPointsByCustomerID(customerID)
}

// UpdatePoints updates a customer's points balance for a given brand, updating the
// customer's points and recording the transaction. This method returns an error if
// any of the operations fail.
func (s *pointService) UpdatePoints(customerID int, brandID int, pointsDelta int, reason string) error {
	//update points transactions

	var transaction = &domain.LealPointsTransaction{
		CustomerID: customerID,
		BrandID:    brandID,
		Change:     pointsDelta,
		Reason:     reason,
	}

	err := s.pointRepo.RecordPointsTransaction(transaction)
	if err != nil {
		return err
	}
	err = s.pointRepo.UpdatePoints(customerID, brandID, pointsDelta)
	if err != nil {
		return err
	}

	//update coins

	err = s.pointRepo.RecordCoins(customerID, pointsDelta)
	if err != nil {
		return err
	}

	return nil
}
