package http

import (
	"github.com/gin-gonic/gin"
)

// NewRouter returns a new gin Engine with all the routes needed for the
// brand service.
func NewRouter(h *Handler) *gin.Engine {
	r := gin.Default()
	//health check
	r.GET("/ping", h.Ping)

	// Brand endpoints
	r.POST("/new-brand", h.NewBrand)
	r.POST("/login-brand", h.LoginBrand)
	r.POST("/new-branch", h.NewBranch)
	r.GET("/my-branches", h.MyBranches)
	r.POST("/new-campaign", h.NewCampaign)
	r.POST("/modify-campaign", h.ModifyCampaign)
	r.GET("/my-campaigns", h.MyCampaigns)
	r.POST("/new-reward", h.NewReward)
	r.GET("/my-rewards", h.MyRewards)
	return r
}
