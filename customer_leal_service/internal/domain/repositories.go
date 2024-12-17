package domain

// Repositorios para acceder a los datos
type CustomerRepository interface {
	CreateCustomer(name string, email string, phone string, pass string, token string) (*Customer, error)
	GetCustomerByID(id int) (*Customer, error)
	GetCustomerByEmail(email string) (*Customer, error)
	UpdateCustomerToken(id int, token string) error
}

type PointsRepository interface {
	GetPointsByCustomerID(id int) ([]LealPoints, error)
	UpdatePoints(customerID int, brandID int, points int) error
	RecordPointsTransaction(transaction *LealPointsTransaction) error
	RecordCoins(customerID int, coins int) error
	GetPoinysByCustomerIDAndBrandID(customerID int, brandID int) (*LealPoints, error)
}

type CoinsRepository interface {
	GetCoinsByCustomerID(id int) (int, error)
	UpdateCustomerCoins(id int, coins int) error
}

type PurchasesRepository interface {
	RecordPurchase(purchase *Purchase) (*Purchase, error)
}

type RedeemedRepository interface {
	RedeemReward(redeemed *Redeemed) (*Redeemed, error)
}
