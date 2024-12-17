package db

import (
	"database/sql"

	"github.com/degarzonm/customer_leal_service/internal/domain"
)

type postgresPurchasesRepo struct {
	db *sql.DB
}

func NewPostgresPurchasesRepo(db *sql.DB) domain.PurchasesRepository {
	return &postgresPurchasesRepo{db: db}
}

// RecordPurchase records a purchase in the database, returning the purchase with the ID and PurchaseDate populated
// or an error if something went wrong
func (r *postgresPurchasesRepo) RecordPurchase(purchase *domain.Purchase) (*domain.Purchase, error) {
	query := `INSERT INTO purchase (customer_id, amount, brand_id, branch_id, coins_used) VALUES ($1, $2, $3 , $4, $5) RETURNING id, purchase_date`

	row := r.db.QueryRow(query, purchase.CustomerID, purchase.Amount, purchase.BrandID, purchase.BranchID, purchase.CoinsUsed)

	if err := row.Scan(&purchase.ID, &purchase.PurchaseDate); err != nil {
		return nil, err
	}
	return purchase, nil

}
