package db

import (
	"database/sql"

	"github.com/degarzonm/customer_leal_service/internal/domain"
)

type postgresCoinsRepo struct {
	db *sql.DB
}

func NewPostgresCoinsRepo(db *sql.DB) domain.CoinsRepository {
	return &postgresCoinsRepo{db: db}
}

// GetCoinsByCustomerID retrieves the number of leal coins for a given customer
// from the database by customer ID. It returns the number of coins or an error
// if the retrieval fails.

func (r *postgresCoinsRepo) GetCoinsByCustomerID(customerID int) (int, error) {
	query := `SELECT leal_coins FROM customer WHERE id = $1`
	row := r.db.QueryRow(query, customerID)
	var coins int
	if err := row.Scan(&coins); err != nil {
		return 0, err
	}
	return coins, nil
}

// UpdateCustomerCoins updates the number of coins a customer has. The number of coins
// is modified by the given delta. If the resulting number of coins is negative, the
// number of coins is reset to 0. The function returns an error if the update query
// fails.
func (r *postgresCoinsRepo) UpdateCustomerCoins(customerID int, coins_delta int) error {
	query := `UPDATE customer SET leal_coins = GREATEST(leal_coins + $1, 0) WHERE id = $2`
	_, err := r.db.Exec(query, coins_delta, customerID)
	return err
}
