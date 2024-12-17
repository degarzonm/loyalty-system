package application

import (
	"errors"

	"github.com/degarzonm/customer_leal_service/internal/domain"
)

type redeemService struct {
	redeemRepo   domain.RedeemedRepository
	pointRepo    domain.PointsRepository
	customerRepo domain.CustomerRepository
}

func NewRedeemService(p domain.RedeemedRepository, pr domain.PointsRepository, cr domain.CustomerRepository) domain.RedeemService {
	return &redeemService{redeemRepo: p, pointRepo: pr, customerRepo: cr}
}

// RedeemReward executes a redeem operation for a given customer and brand
//
// It will first verify if the customer has enough points to redeem the reward,
// then it will update the customer points, create a transaction for the points
// modification, and finally record the redeem in the database.
//
// If any of the steps fail, it will return an error.
func (s *redeemService) RedeemReward(redeem *domain.Redeemed) (*domain.Redeemed, error) {
	//get customer points given the customer id and brand id
	points, err := s.pointRepo.GetPoinysByCustomerIDAndBrandID(redeem.CustomerID, redeem.BrandID)
	if err != nil {
		return nil, err
	}

	if points.Points < redeem.PointsSpend {
		return nil, errors.New("not enough points")
	}

	//update customer points
	err = s.pointRepo.UpdatePoints(redeem.CustomerID, redeem.BrandID, -redeem.PointsSpend)
	if err != nil {
		return nil, err
	}

	//update points transactions

	var transaction = &domain.LealPointsTransaction{
		CustomerID: redeem.CustomerID,
		BrandID:    redeem.BrandID,
		Change:     -redeem.PointsSpend,
		Reason:     "Redeem reward",
	}

	err = s.pointRepo.RecordPointsTransaction(transaction)
	if err != nil {
		return nil, err
	}

	// record redeem

	redeem, err = s.redeemRepo.RedeemReward(redeem)
	if err != nil {
		return nil, err
	}

	return redeem, nil

}
