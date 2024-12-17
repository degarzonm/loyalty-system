package http

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/degarzonm/brand_leal_service/internal/domain"
	"github.com/degarzonm/brand_leal_service/internal/infrastructure/util"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	brandService    domain.BrandService
	branchService   domain.BranchService
	campaignService domain.CampaignService
	rewardService   domain.RewardService
}

func NewHandler(bs domain.BrandService, bss domain.BranchService, cs domain.CampaignService, r domain.RewardService) *Handler {
	return &Handler{brandService: bs, branchService: bss, campaignService: cs, rewardService: r}
}

// Ping checks if the service is up and running.
// It returns a simple JSON message "pong_brands".
func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong_brands"})
}

// NewBrand creates a new brand and returns its brand_id and token.
// It requires a JSON object with a brand_name and a pass field.
// If the request is correct, it returns a JSON with the brand_id and the token.
// If the request is incorrect, it returns a JSON with an error message.
// If the service has a problem, it returns a JSON with an error message.
func (h *Handler) NewBrand(c *gin.Context) {
	log.Println("NewBrand Request")
	var req NewBrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("object received: ", req)
	req.BrandName = util.Sanitize(req.BrandName)
	b, err := h.brandService.CreateBrand(req.BrandName, req.Pass)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"brand_id": b.ID, "token": b.Token})
}

// LoginBrand logs in a brand and returns its brand_id and token.
// It requires a JSON object with a brand_name and a pass field.
// If the request is correct, it returns a JSON with the brand_id and the token.
// If the request is incorrect, it returns a JSON with an error message.
// If the service has a problem, it returns a JSON with an error message.
func (h *Handler) LoginBrand(c *gin.Context) {
	var req LoginBrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.BrandName = util.Sanitize(req.BrandName)
	b, err := h.brandService.LoginBrand(req.BrandName, req.Pass)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"brand_id": b.ID, "token": b.Token})
}

// NewBranch creates a new branch for the authorized brand.
// It requires a JSON object with a branch_name field.
// If the brand is not authorized, it returns a 401 Unauthorized error.
// If the JSON binding fails, it returns a 400 Bad Request error.
// If the branch creation fails, it returns a 500 Internal Server Error.
// On success, it returns a 200 OK status with the branch_id, brand_id, and branch_name.

func (h *Handler) NewBranch(c *gin.Context) {
	brandID, err := h.auhorizeBrand(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req NewBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.BranchName = util.Sanitize(req.BranchName)
	br, err := h.branchService.CreateBranch(brandID, req.BranchName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"branch_id": br.ID, "brand_id": br.BrandID, "branch_name": br.Name})
}

// MyBranches returns all the branches of the authorized brand.
// It requires a valid JWT token in the Authorization header.
// If the brand is not authorized, it returns a 401 Unauthorized error.
// If the service has a problem, it returns a 500 Internal Server Error.
// On success, it returns a 200 OK status with a JSON object
// with a list of branch objects.
func (h *Handler) MyBranches(c *gin.Context) {

	brandID, err := h.auhorizeBrand(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	branches, err := h.branchService.GetBranches(brandID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"branches": branches})
}

// NewCampaign creates a new campaign for the authorized brand.
// It requires a valid JWT token in the Authorization header.
// If the brand is not authorized, it returns a 401 Unauthorized error.
// If the request is invalid, it returns a 400 Bad Request error.
// If the service has a problem, it returns a 500 Internal Server Error.
// On success, it returns a 200 OK status with a JSON object
// with the campaign ID.
func (h *Handler) NewCampaign(c *gin.Context) {

	brandID, err := h.auhorizeBrand(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req NewCampaignRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar y convertir BranchIDs
	branchIDs, err := util.ParseBranchIDs(req.BranchIDs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid branch_ids"})
		return
	}

	start, err := util.ParseDate(req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date"})
		return
	}
	end, err := util.ParseDate(req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date"})
		return
	}

	campaign := &domain.Campaign{
		CampaignName: req.CampaignName,
		BrandID:      brandID,
		MinValue:     req.MinValue,
		MaxValue:     req.MaxValue,
		StartDate:    start,
		EndDate:      end,
		PointFactor:  req.PointFactor,
		CoinFactor:   req.CoinFactor,
		Status:       req.Status,
	}

	camp, err := h.campaignService.CreateCampaign(campaign, branchIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"campaign_id": camp.ID})
}

// ModifyCampaign modify a campaign, given a valid brand, campaign id
// and new values for the campaign. The campaign is modified in the
// given branches. If the campaign is not found, or the brand is not
// authorized, an error is returned. If the dates are invalid, an
// error is returned. If there is an error in the database, an error
// is returned.
func (h *Handler) ModifyCampaign(c *gin.Context) {
	brandID, err := h.auhorizeBrand(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	var req ModifyCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	branchIDs, err := util.ParseBranchIDs(req.BranchIDs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid branch_ids"})
		return
	}

	start, err := util.ParseDate(req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date"})
		return
	}
	end, err := util.ParseDate(req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date"})
		return
	}

	campaign := &domain.Campaign{
		ID:           req.CampaignID,
		CampaignName: req.CampaignName,
		BrandID:      brandID,
		MinValue:     req.MinValue,
		MaxValue:     req.MaxValue,
		StartDate:    start,
		EndDate:      end,
		PointFactor:  req.PointFactor,
		CoinFactor:   req.CoinFactor,
		Status:       req.Status,
	}

	campaign, err = h.campaignService.UpdateCampaign(campaign, branchIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"campaign_id": campaign.ID})
}

// MyCampaigns returns all the campaigns of the brand of the given token.
// If the token is invalid, an error is returned. If there is an error
// in the database, an error is returned.
func (h *Handler) MyCampaigns(c *gin.Context) {
	brandID, err := h.auhorizeBrand(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	campaigns, err := h.campaignService.GetCampaigns(brandID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"campaigns": campaigns})
}

// NewReward creates a new reward for the authorized brand.
// It requires a valid JWT token in the Authorization header.
// If the brand is not authorized, it returns a 401 Unauthorized error.
// If the request is invalid, it returns a 400 Bad Request error.
// If the service has a problem, it returns a 500 Internal Server Error.
// On success, it returns a 200 OK status with a JSON object
// with the reward ID.
func (h *Handler) NewReward(c *gin.Context) {
	brandID, err := h.auhorizeBrand(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	var req NewRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	start, err := util.ParseDate(req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date"})
		return
	}
	end, err := util.ParseDate(req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date"})
		return
	}
	reward := &domain.Reward{
		BrandId:     brandID,
		RewardName:  req.RewardName,
		PricePoints: req.PricePoints,
		StartDate:   start,
		EndDate:     end,
	}
	reward, err = h.rewardService.CreateReward(reward)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"reward_id": reward.ID})
}

// MyRewards returns all the rewards of the authorized brand.
// It requires a valid brand ID and token in the request headers.
// If the brand is not authorized, it returns a 401 Unauthorized error.
// If there is an error retrieving rewards from the service, it returns a 500 Internal Server Error.
// On success, it returns a 200 OK status with a JSON object containing a list of rewards.

func (h *Handler) MyRewards(c *gin.Context) {
	brandID, err := h.auhorizeBrand(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	rewards, err := h.rewardService.GetRewardsByBrand(brandID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"rewards": rewards})
}

// auhorizeBrand checks if the given brand ID and token are valid.
// It reads the required headers from the gin.Context object.
// If any of the required headers are missing, or if the validation fails, it returns an error.
// On success, it returns the brand ID and a nil error.
func (h *Handler) auhorizeBrand(c *gin.Context) (int, error) {
	brandIDStr := c.GetHeader("Leal-Brand-Id")
	if brandIDStr == "" {
		return 0, errors.New("brand_id header required")
	}

	brandID, err := strconv.Atoi(brandIDStr)
	if err != nil {
		return 0, errors.New("invalid brand_id")
	}

	tokenReq := c.GetHeader("Leal-Brand-Token")
	if tokenReq == "" {
		return 0, errors.New("token header required")
	}

	err = h.brandService.ValidateToken(brandID, tokenReq)
	if err != nil {
		return 0, err
	}

	return brandID, nil

}
