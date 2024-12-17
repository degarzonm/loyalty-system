package domain

type BrandService interface {
	CreateBrand(name, pass string) (*Brand, error)
	LoginBrand(name, pass string) (*Brand, error)
	ValidateToken(brandID int, token string) error
}

type BranchService interface {
	CreateBranch(brandID int, branchName string) (*Branch, error)
	GetBranches(brandID int) ([]Branch, error)
}

type CampaignService interface {
	CreateCampaign(campaign *Campaign, branchIds []int) (*Campaign, error)
	UpdateCampaign(campaign *Campaign, branchIds []int) (*Campaign, error)
	GetCampaigns(brandID int) ([]Campaign, error)
}

type RewardService interface {
	CreateReward(reward *Reward) (*Reward, error)
	GetRewardsByBrand(brandID int) ([]Reward, error)
}
