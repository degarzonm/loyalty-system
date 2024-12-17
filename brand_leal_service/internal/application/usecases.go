package application

import (
	"errors"
	"log"
	"time"

	"github.com/degarzonm/brand_leal_service/internal/domain"
	"github.com/degarzonm/brand_leal_service/internal/infrastructure/util"
)

type brandService struct {
	brandRepo    domain.BrandsRepository
	campaignRepo domain.CampaignRepository
}

type branchService struct {
	branchRepo domain.BranchesRepository
}

type campaignService struct {
	campaignRepo domain.CampaignRepository
}

type rewardService struct {
	rewardRepo domain.RewardRepository
}

func NewBrandService(br domain.BrandsRepository, cr domain.CampaignRepository) domain.BrandService {
	return &brandService{brandRepo: br, campaignRepo: cr}
}

func NewBranchService(br domain.BranchesRepository) domain.BranchService {
	return &branchService{branchRepo: br}
}

func NewCampaignService(cr domain.CampaignRepository) domain.CampaignService {
	return &campaignService{campaignRepo: cr}
}

func NewRewardService(r domain.RewardRepository) domain.RewardService {
	return &rewardService{rewardRepo: r}
}

// CreateBrand creates a new brand in the database. It requires a name and password, and returns
// a new brand object and an error. If the name or password are empty, or if there is an error
// generating the token or creating the brand, the function returns an error. The function also
// creates a base campaign for the brand, with a name of "base", a start date of January 1, 2000,
// and an end date of January 1, 2100. The campaign is created with a point factor and coin factor
// of 0.001, and a customer count of 0. The campaign status is set to "active".
func (s *brandService) CreateBrand(name, pass string) (*domain.Brand, error) {

	if name == "" || pass == "" {
		return nil, errors.New("name or password are empty")
	}

	passHash := util.HashPassword(pass)

	token, err := util.GenerateToken()
	if err != nil {
		return nil, errors.New("error generating token")
	}
	newBrand, err := s.brandRepo.CreateBrand(name, passHash, token)
	if err != nil {
		return nil, err
	}
	//generate base campaign
	baseCampaign := &domain.Campaign{
		CampaignName:  "base",
		BrandID:       newBrand.ID,
		MinValue:      0.0,
		MaxValue:      1000000000.0,
		StartDate:     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:       time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC),
		PointFactor:   0.001,
		CoinFactor:    0.001,
		CustomerCount: 0,
		Status:        "active",
	}

	_, err = s.campaignRepo.CreateCampaign(baseCampaign, []int{})
	if err != nil {
		return nil, err
	}

	return newBrand, err
}

// LoginBrand validates a brand's credentials and returns the brand object and a token for authentication, or an error. If the name or password are empty, or if there is an error retrieving the brand or generating the token, the function returns an error. The function also updates the brand's token in the database. If the brand does not exist, or if the password is invalid, the function returns an error.
func (s *brandService) LoginBrand(name, pass string) (*domain.Brand, error) {

	if name == "" || pass == "" {
		return nil, errors.New("name or password are empty")
	}

	b, err := s.brandRepo.GetBrandByName(name)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, errors.New("brand not found")
	}

	log.Println("brand found:", b)

	if !util.CheckPassHash(pass, b.PassHash) {
		return nil, errors.New("invalid password")
	}
	// update token
	newToken, err := util.GenerateToken()
	if err != nil {
		return nil, errors.New("failed generating token")
	}

	err = s.brandRepo.UpdateBrandToken(b.ID, newToken)
	if err != nil {
		return nil, err
	}
	b.Token = newToken
	return b, nil
}

// ValidateToken checks if the provided token matches the stored token for the brand with the given brandID.
// It returns an error if the brand is not found or if the token is invalid.

func (s *brandService) ValidateToken(brandID int, token string) error {
	b, err := s.brandRepo.GetBrandByID(brandID)
	if err != nil {
		return err
	}
	if b == nil {
		return errors.New("brand not found")
	}
	if b.Token != token {
		return errors.New("invalid token")
	}
	return nil
}

// CreateBranch adds a new branch for the specified brand. It takes the brand ID and the branch name as inputs,
// and returns the newly created branch object or an error. If the branch creation in the repository fails,
// or if linking the branch to the base campaign fails, it returns an error.

func (s *branchService) CreateBranch(brandID int, branchName string) (*domain.Branch, error) {

	newBranch, err := s.branchRepo.CreateBranch(brandID, branchName)
	if err != nil {
		return nil, err
	}

	// Step 2: Link the branch to the campaign
	err = s.branchRepo.LinkBranchToBaseCampaign(newBranch)
	if err != nil {
		return nil, err
	}

	return newBranch, nil
}

func (s *branchService) GetBranches(brandID int) ([]domain.Branch, error) {
	return s.branchRepo.GetBranchesByBrandID(brandID)
}

// CreateCampaign creates a new campaign in the database. It takes a campaign object and a list of
// branch IDs as inputs, and returns the newly created campaign object or an error. If the campaign
// start date is after its end date, it returns an error. The function also creates a new campaign
// branch for each of the provided branch IDs.
func (s *campaignService) CreateCampaign(campaign *domain.Campaign, branches []int) (*domain.Campaign, error) {

	if campaign.StartDate.After(campaign.EndDate) {
		return nil, errors.New("start_date cannot be after end_date")
	}
	return s.campaignRepo.CreateCampaign(campaign, branches)
}

// UpdateCampaign updates a campaign in the database. It takes a campaign object and a list of
// branch IDs as inputs, and returns the updated campaign object or an error. If the campaign
// start date is after its end date, it returns an error. The function also updates the campaign
// branches by deleting the existing ones and inserting the new ones provided. If the campaign
// is not found, it returns an error.
func (s *campaignService) UpdateCampaign(campaign *domain.Campaign, branches []int) (*domain.Campaign, error) {
	existing, err := s.campaignRepo.GetCampaignByID(campaign.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("campaign not found")
	}
	if campaign.StartDate.After(campaign.EndDate) {
		return nil, errors.New("start_date cannot be after end_date")
	}
	err = s.campaignRepo.UpdateCampaign(campaign, branches)
	return campaign, err
}

// GetCampaigns retrieves all campaigns for a given brand ID from the database.
//
// The function returns a slice of Campaign objects, or an error if there is a problem
// communicating with the database. The campaigns are sorted in descending order of their start
// dates.
func (s *campaignService) GetCampaigns(brandID int) ([]domain.Campaign, error) {
	return s.campaignRepo.GetCampaignsByBrandID(brandID)
}

// CreateReward creates a new reward in the database. It takes a reward object as input, and returns
// the newly created reward object or an error. The function also sets the reward ID of the provided
// reward object to the newly created reward ID.
func (r *rewardService) CreateReward(reward *domain.Reward) (*domain.Reward, error) {
	return r.rewardRepo.CreateReward(reward)
}

// GetRewardsByBrand retrieves all rewards for a given brand ID from the database.
//
// The function returns a slice of Reward objects, or an error if there is a problem
// communicating with the database. The rewards are sorted in descending order of their start
// dates.
func (r *rewardService) GetRewardsByBrand(brandID int) ([]domain.Reward, error) {
	return r.rewardRepo.GetRewardsByBrand(brandID)
}
