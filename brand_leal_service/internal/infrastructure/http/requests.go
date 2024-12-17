package http

type NewBrandRequest struct {
	BrandName string `json:"brand_name"`
	Pass      string `json:"pass"`
}

type LoginBrandRequest struct {
	BrandName string `json:"brand_name"`
	Pass      string `json:"pass"`
}

type NewBranchRequest struct {
	BranchName string `json:"branch_name"`
}

type NewCampaignRequest struct {
	CampaignName string  `json:"campaign_name"`
	BrandID      int     `json:"brand_id"`
	BranchIDs    string  `json:"branch_ids"` // comma separated ids
	MinValue     float64 `json:"min_value"`
	MaxValue     float64 `json:"max_value"`
	StartDate    string  `json:"start_date"`
	EndDate      string  `json:"end_date"`
	Status       string  `json:"status"`
	PointFactor  float64 `json:"point_factor"`
	CoinFactor   float64 `json:"coin_factor"`
}

type ModifyCampaignRequest struct {
	CampaignID   int     `json:"campaign_id"`
	CampaignName string  `json:"campaign_name"`
	BrandID      int     `json:"brand_id"`
	BranchIDs    string  `json:"branch_id"`
	MinValue     float64 `json:"min_value"`
	MaxValue     float64 `json:"max_value"`
	StartDate    string  `json:"start_date"`
	EndDate      string  `json:"end_date"`
	Status       string  `json:"status"`
	PointFactor  float64 `json:"point_factor"`
	CoinFactor   float64 `json:"coin_factor"`
}

type NewRewardRequest struct {
	BrandID     int    `json:"brand_id"`
	RewardName  string `json:"reward_name"`
	PricePoints int    `json:"price_points"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}
