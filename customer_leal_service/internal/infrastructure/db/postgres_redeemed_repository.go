package db

import (
	"database/sql"

	"github.com/degarzonm/customer_leal_service/internal/domain"
)

type postgresRedeemedRepo struct {
	db *sql.DB
}

func NewPostgresRedeemedRepo(db *sql.DB) domain.RedeemedRepository {
	return &postgresRedeemedRepo{db: db}
}

// RedeemReward records a redeem operation for a given customer and brand.
//
// It will record the operation in the database and return the
// Redeemed struct with the ID and Date fields populated. If any error
// occurs, it will return that error.
func (r *postgresRedeemedRepo) RedeemReward(redeemed *domain.Redeemed) (*domain.Redeemed, error) {
	query := `INSERT INTO redeemed (customer_id, brand_id, reward_id, points_spend) VALUES ($1, $2 , $3, $4) RETURNING id , date`

	row := r.db.QueryRow(query, redeemed.CustomerID, redeemed.BrandID, redeemed.RewardID, redeemed.PointsSpend)

	if err := row.Scan(&redeemed.ID, &redeemed.Date); err != nil {
		return nil, err
	}

	return redeemed, nil
}
