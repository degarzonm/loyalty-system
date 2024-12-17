package db

import (
	"database/sql"
	"log"

	"github.com/degarzonm/customer_leal_service/internal/domain"
)

type postgresPointsRepo struct {
	db *sql.DB
}

func NewPostgresPointsRepo(db *sql.DB) domain.PointsRepository {
	return &postgresPointsRepo{db: db}
}

// GetPointsByCustomerID retrieves all loyalty points associated with a given customer ID.
// It returns a slice of LealPoints entities, each containing the customer ID, brand ID, and points.
// If the query encounters an error, it returns the error. If no points are found, it returns an empty slice.

func (r *postgresPointsRepo) GetPointsByCustomerID(customer_id int) ([]domain.LealPoints, error) {
	query := `SELECT  customer_id, brand_id, points FROM leal_points WHERE customer_id = $1`
	rows, err := r.db.Query(query, customer_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.LealPoints
	for rows.Next() {
		var lp domain.LealPoints
		if err := rows.Scan(&lp.CustomerID, &lp.BrandID, &lp.Points); err != nil {
			return nil, err
		}
		results = append(results, lp)
	}

	return results, nil
}

// UpdatePoints updates the points balance for a given customer ID and brand ID by the given delta.
// The delta can be positive or negative. If the resulting points balance would be negative, it is
// set to 0. If there is no current points balance for the customer ID and brand ID, it creates
// one. The update is done in a single atomic operation. If the operation encounters an error,
// it returns the error.
func (r *postgresPointsRepo) UpdatePoints(customerID int, brandID int, delta int) error {
	query := `
		INSERT INTO leal_points (customer_id, brand_id, points)
		VALUES ($1, $2, GREATEST($3, 0))
		ON CONFLICT (customer_id, brand_id)
		DO UPDATE
		SET points = GREATEST(leal_points.points + $3, 0)
	`
	_, err := r.db.Exec(query, customerID, brandID, delta)
	return err
}

// RecordPointsTransaction records a points transaction in the leal_points_transactions table.
// It logs the transaction for debugging purposes and returns an error if the operation fails.
func (r *postgresPointsRepo) RecordPointsTransaction(transaction *domain.LealPointsTransaction) error {
	log.Println("Db received transaction: ", transaction)
	query := `INSERT INTO leal_points_transactions (customer_id, brand_id, change, reason) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Query(query, transaction.CustomerID, transaction.BrandID, transaction.Change, transaction.Reason)
	return err
}

// RecordCoins updates the coin balance for a given customer by adding the specified number of coins.
// It logs the transaction details for debugging purposes. The function returns an error if the update
// operation fails.

func (r *postgresPointsRepo) RecordCoins(customerID int, coins int) error {
	log.Println("Db received coins: ", coins, " for customer: ", customerID)
	query := `UPDATE customer SET coins = coins + $1 WHERE id = $2`
	_, err := r.db.Query(query, coins, customerID)
	if err != nil {
		return err
	}
	return nil
}

// GetPoinysByCustomerIDAndBrandID retrieves the points balance for a given customer ID and brand ID.
// If the query encounters an error, it returns the error. If no points are found, it returns nil.
func (r *postgresPointsRepo) GetPoinysByCustomerIDAndBrandID(customerID int, brandID int) (*domain.LealPoints, error) {
	query := `SELECT brand_id, points FROM leal_points WHERE customer_id = $1 AND brand_id = $2`
	row := r.db.QueryRow(query, customerID, brandID)
	var points domain.LealPoints
	if err := row.Scan(&points.BrandID, &points.Points); err != nil {
		return nil, err
	}
	return &points, nil
}
