package domain

import "time"

// Entidades principales
type Customer struct {
	ID               int
	Name             string
	Email            string
	Phone            string
	PassHash         string
	Token            string
	LealCoins        int
	RegistrationDate time.Time
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

type LealPoints struct {
	ID         int
	CustomerID int
	BrandID    int
	Points     int
}

type LealPointsTransaction struct {
	ID         int
	CustomerID int
	BrandID    int
	Change     int
	Reason     string
	Date       time.Time
}

type Reward struct {
	ID          int
	BrandID     int
	RewardName  string
	PricePoints int
	StartDate   time.Time
	EndDate     time.Time
}

type Redeemed struct {
	ID          int
	CustomerID  int
	BrandID     int
	RewardID    int
	PointsSpend int
	Date        time.Time
}

type LealPointsApply struct {
	CustomerID int
	BrandID    int
	Points     int
	Coins      int
	Reason     string
}
