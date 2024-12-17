package application

import (
	"errors"
	"log"

	"github.com/degarzonm/brand_leal_service/internal/config"
	"github.com/degarzonm/brand_leal_service/internal/domain"
)

type AppService struct {
	campaignRepo  domain.CampaignRepository
	brandRepo     domain.BrandsRepository
	eventProducer domain.EventProducer
}

// NewAppService creates a new application service
func NewAppService(campaignRepo domain.CampaignRepository, brandRepo domain.BrandsRepository, producer domain.EventProducer) *AppService {
	return &AppService{
		campaignRepo:  campaignRepo,
		brandRepo:     brandRepo,
		eventProducer: producer,
	}
}

// ProcessPurchase processes a purchase transaction by calculating and applying points
// and coins based on base and active campaigns. It retrieves the base campaign for
// the brand and active campaigns for the branch, applying factors to the purchase
// amount to compute total points and coins. If the purchase meets campaign criteria,
// additional campaign bonuses are applied. The computed points and coins are then
// logged and sent as a message to a Kafka topic. Returns an error if any operation
// within the process fails.

func (s *AppService) ProcessPurchase(purchase domain.Purchase) error {
	// Get the base campaign
	baseCampaign, err := s.campaignRepo.GetBaseCampaignForBrand(purchase.BrandID)
	if err != nil {
		return errors.New("failed to retrieve base campaign")
	}
	if baseCampaign == nil {
		return errors.New("no base campaign found for the brand")
	}

	// Calculate base points and coins
	basePoints := purchase.Amount * baseCampaign.PointFactor
	baseCoins := purchase.Amount * baseCampaign.CoinFactor

	totalPoints := basePoints
	totalCoins := baseCoins

	// Fetch active campaigns for the branch
	campaigns, err := s.campaignRepo.GetCampaignsForBranch(purchase.BranchID)
	if err != nil {
		return errors.New("failed to retrieve campaigns for branch")
	}

	// Apply additional campaigns
	for _, campaign := range campaigns {
		if campaign.Status != "active" || purchase.Amount < campaign.MinValue || purchase.Amount > campaign.MaxValue {
			continue
		}
		if purchase.PurchaseDate.Before(campaign.StartDate) || purchase.PurchaseDate.After(campaign.EndDate) {
			continue
		}
		//update counter on campaign
		err = s.campaignRepo.UpdateCustomerCountCampaign(&campaign)
		if err != nil {
			return err
		}
		//add points and coins
		totalPoints += basePoints * campaign.PointFactor
		totalCoins += baseCoins * campaign.CoinFactor
	}
	log.Printf("Processed purchase for CustomerID=%d, Points=%f, Coins=%f\n", purchase.CustomerID, totalPoints, totalCoins)
	// Send calculated pointsInfo to Kafka
	pointsInfo := domain.LealPointsApply{
		CustomerID: purchase.CustomerID,
		BrandID:    purchase.BrandID,
		Points:     int(totalPoints),
		Coins:      int(totalCoins),
		Reason:     "purchase",
	}
	log.Println("message to sent to pointsInfo: ", pointsInfo)
	if err := s.SendApplyPointsEvent(pointsInfo); err != nil {
		return errors.New("failed to send apply points event")
	}

	log.Printf("Processed purchase for CustomerID=%d successfully", purchase.CustomerID)
	return nil
}

// SendApplyPointsEvent sends a message to the apply-points topic with the given points info.
// The message is sent to the configured topic name from the global configuration.
// The method returns an error if the message could not be sent.
func (s *AppService) SendApplyPointsEvent(points domain.LealPointsApply) error {
	cfg := config.GetConfig()
	return s.eventProducer.SendMessage(cfg.MsgApplyPointsTopic, points)
}
