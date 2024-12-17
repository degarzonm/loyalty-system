package http

import (
	"github.com/gin-gonic/gin"
)

// NewRouter returns a new gin Engine with all the routes needed for the
// customer service.
func NewRouter(h *Handler) *gin.Engine {
	r := gin.Default()

	// Health check endpoint
	r.GET("/ping", h.Ping)

	// Customer endpoints
	r.POST("/new-customer", h.NewCustomer)
	r.POST("/login-customer", h.LoginCustomer)
	r.GET("/my-points/", h.GetCustomerPoints)
	r.GET("/my-coins/", h.GetCustomerCoins)
	r.POST("/redeem", h.Redeem)
	r.POST("/purchase", h.Purchase)

	return r
}
