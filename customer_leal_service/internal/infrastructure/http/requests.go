package http

type NewCustomerRequest struct {
	CustomerName string `json:"customer_name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Pass         string `json:"pass"`
}

type LoginCustomerRequest struct {
	CustomerEmail string `json:"email"`
	Pass          string `json:"pass"`
}

type PurchaseRequest struct {
	CustomerID int     `json:"customer_id"`
	Amount     float64 `json:"amount"`
	BrandID    int     `json:"brand_id"`
	BranchID   int     `json:"branch_id"`
	CoinsUsed  int     `json:"coins_used"`
}

type RedeemRewardRequest struct {
	CustomerID  int `json:"customer_id"`
	BrandID     int `json:"brand_id"`
	RewardID    int `json:"reward_id"`
	PointsSpend int `json:"points_spend"`
}
