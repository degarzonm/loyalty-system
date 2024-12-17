package domain

type BrandsRepository interface {
	GetBrandByID(id int) (*Brand, error)
	GetBrandByName(brandName string) (*Brand, error)
	CreateBrand(brandName, passHash, token string) (*Brand, error)
	UpdateBrandToken(brandID int, token string) error
}

type BranchesRepository interface {
	CreateBranch(brandID int, branchName string) (*Branch, error)
	GetBranchesByBrandID(brandID int) ([]Branch, error)
	LinkBranchToBaseCampaign(branchID *Branch) error
}

type CampaignRepository interface {
	CreateCampaign(c *Campaign, branchIDs []int) (*Campaign, error)
	GetCampaignByID(id int) (*Campaign, error)
	UpdateCampaign(c *Campaign, branchIDs []int) error
	UpdateCustomerCountCampaign(c *Campaign) error
	GetCampaignsByBrandID(brandID int) ([]Campaign, error)
	GetBranchesForCampaign(campaignID int) ([]Branch, error)
	GetCampaignsForBranch(branchID int) ([]Campaign, error)
	GetBaseCampaignForBrand(brandID int) (*Campaign, error)
}

type RewardRepository interface {
	CreateReward(r *Reward) (*Reward, error)
	GetRewardsByBrand(id int) ([]Reward, error)
}
