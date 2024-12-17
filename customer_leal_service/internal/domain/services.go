package domain

type CustomerService interface {
	CreateCustomer(name string, email string, phone string, pass string) (*Customer, error)
	LoginCustomer(email string, pass string) (*Customer, error)
	GetCustomerByID(id int) (*Customer, error)
	ValidateToken(customerID int, token string) error
}
type PointService interface {
	GetCustomerPoints(customerID int) ([]LealPoints, error)
	UpdatePoints(customerID int, brandID int, pointsDelta int, reason string) error
}

type CoinService interface {
	GetCustomerCoins(customerID int) (int, error)
	UpdateCustomerCoins(id int, coins int) error
}

type PurchaseService interface {
	ProcessPurchase(attempPurchase *Purchase) (*Purchase, error)
}

type RedeemService interface {
	RedeemReward(redeem *Redeemed) (*Redeemed, error)
}
