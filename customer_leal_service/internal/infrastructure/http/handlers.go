package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/degarzonm/customer_leal_service/internal/domain"
	"github.com/degarzonm/customer_leal_service/internal/infrastructure/util"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	customerService domain.CustomerService
	pointsService   domain.PointService
	coinService     domain.CoinService
	purchaseService domain.PurchaseService
	redeemService   domain.RedeemService
}

func NewHandler(cs domain.CustomerService, ps domain.PointService, ccs domain.CoinService, pcs domain.PurchaseService, rs domain.RedeemService) *Handler {
	return &Handler{customerService: cs, pointsService: ps, coinService: ccs, purchaseService: pcs, redeemService: rs}
}

func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong_customers"})
}

// NewCustomer handles the request to register a new customer. The customer's
// name, email, phone and password should be sent in the request body in JSON
// format. The response will contain the customer's ID and a token which can be
// used to authenticate the customer in future requests. If the request is
// invalid, the response will contain an error message. If the customer service
// fails to create the customer, the response will contain an error message with
// status code 500.
func (h *Handler) NewCustomer(c *gin.Context) {
	var req NewCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, err := h.customerService.CreateCustomer(req.CustomerName, req.Email, req.Phone, req.Pass)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"customer_id": customer.ID, "token": customer.Token})
}

// LoginCustomer authenticates a customer using their email and password.
// The request should contain a JSON object with customer_email and pass fields.
// If authentication is successful, it responds with the customer's ID and token.
// If the request is malformed, it responds with a 400 status code and an error message.
// If authentication fails, it responds with a 403 status code and an error message.

func (h *Handler) LoginCustomer(c *gin.Context) {
	var req LoginCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.CustomerEmail = util.Sanitize(req.CustomerEmail)
	b, err := h.customerService.LoginCustomer(req.CustomerEmail, req.Pass)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"customer_id": b.ID, "token": b.Token})
}

// GetCustomerPoints retrieves the loyalty points for an authorized customer.
// The customer ID is extracted from the request headers after authorization.
// If the authorization fails, a 500 status code and an error message are returned.
// On success, it returns the points in a JSON response with a 200 status code.
// If any error occurs while fetching the points, a 500 status code and an error message are returned.

func (h *Handler) GetCustomerPoints(c *gin.Context) {
	customerID, err := h.authorizeCustomer(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	points, err := h.pointsService.GetCustomerPoints(customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"points": points})
}

// GetCustomerCoins retrieves the coins for an authorized customer.
// The customer ID is extracted from the request headers after authorization.
// If the authorization fails, a 500 status code and an error message are returned.
// On success, it returns the coins in a JSON response with a 200 status code.
// If any error occurs while fetching the coins, a 500 status code and an error message are returned.
func (h *Handler) GetCustomerCoins(c *gin.Context) {
	customerID, err := h.authorizeCustomer(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	coins, err := h.coinService.GetCustomerCoins(customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"coins": coins})

}

// Redeem exchanges points for a reward.
// The request should contain a JSON object with brand_id, reward_id and points_spend fields.
// If the request is malformed, it responds with a 400 status code and an error message.
// If the authorization fails, a 500 status code and an error message are returned.
// On success, it returns the redeem ID in a JSON response with a 200 status code.
// If any error occurs while redeeming the points, a 500 status code and an error message are returned.
func (h *Handler) Redeem(c *gin.Context) {
	customerID, err := h.authorizeCustomer(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var req RedeemRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	redeem, err := h.redeemService.RedeemReward(
		&domain.Redeemed{
			CustomerID:  customerID,
			BrandID:     req.BrandID,
			RewardID:    req.RewardID,
			PointsSpend: req.PointsSpend,
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"redeem_id": redeem.ID})

}

// Purchase processes a purchase of a customer.
// The request should contain a JSON object with amount, brand_id, branch_id and coins_used fields.
// If the request is malformed, it responds with a 400 status code and an error message.
// If the authorization fails, a 500 status code and an error message are returned.
// On success, it returns the purchase ID in a JSON response with a 200 status code.
// If any error occurs while processing the purchase, a 500 status code and an error message are returned.
func (h *Handler) Purchase(c *gin.Context) {
	customerID, err := h.authorizeCustomer(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var req PurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	purchase, err := h.purchaseService.ProcessPurchase(&domain.Purchase{CustomerID: customerID,
		Amount:    req.Amount,
		BrandID:   req.BrandID,
		BranchID:  req.BranchID,
		CoinsUsed: req.CoinsUsed})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"purchase_id": purchase.ID})
}

// authorizeCustomer validates the customer based on headers in the gin.Context.
// It requires "Leal-Customer-Id" and "Leal-Customer-Token" headers to be present.
// Returns the customer ID if successful, otherwise returns an error if headers are missing,
// if the customer ID is invalid, or if token validation fails.

func (h *Handler) authorizeCustomer(c *gin.Context) (int, error) {
	customerIDStr := c.GetHeader("Leal-Customer-Id")
	if customerIDStr == "" {
		return 0, errors.New("customer_id header required")
	}

	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil {
		return 0, errors.New("invalid customer_id")
	}

	tokenReq := c.GetHeader("Leal-Customer-Token")
	if tokenReq == "" {
		return 0, errors.New("token header required")
	}

	err = h.customerService.ValidateToken(customerID, tokenReq)
	if err != nil {
		return 0, err
	}

	return customerID, nil
}
