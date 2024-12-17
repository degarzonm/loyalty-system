package domain

import "time"

type Brand struct {
	ID               int
	Name             string
	Token            string
	PassHash         string
	RegistrationDate time.Time
}

type Branch struct {
	ID               int
	BrandID          int
	Name             string
	RegistrationDate time.Time
}

type Campaign struct {
	ID            int
	CampaignName  string
	BrandID       int
	MinValue      float64
	MaxValue      float64
	StartDate     time.Time
	EndDate       time.Time
	PointFactor   float64
	CoinFactor    float64
	CustomerCount int
	Status        string
	Branches      []int
}

type Reward struct {
	ID          int
	BrandId     int
	RewardName  string
	PricePoints int
	StartDate   time.Time
	EndDate     time.Time
}
type Purchase struct {
	ID           int
	CustomerID   int
	Amount       float64
	PurchaseDate time.Time
	BrandID      int
	BranchID     int
	CoinsUsed    int
}

type LealPointsApply struct {
	CustomerID int
	BrandID    int
	Points     int
	Coins      int
	Reason     string
}
