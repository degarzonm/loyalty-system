package db

import (
	"database/sql"

	"github.com/degarzonm/brand_leal_service/internal/domain"
)

type postgresRewardRepo struct {
	db *sql.DB
}

func NewPostgresRewardRepo(db *sql.DB) domain.RewardRepository {
	return &postgresRewardRepo{db: db}
}

// CreateReward creates a new reward in the database. It takes a reward object as input, and returns
// the newly created reward object or an error. The function also sets the reward ID of the provided
// reward object to the newly created reward ID.
func (r *postgresRewardRepo) CreateReward(reward *domain.Reward) (*domain.Reward, error) {
	query := `INSERT INTO reward (brand_id, reward_name, price_points, start_date,end_date) 
	VALUES ($1, $2, $3, $4, $5) 
	RETURNING id, reward_name, price_points`
	row := r.db.QueryRow(query, reward.BrandId, reward.RewardName, reward.PricePoints, reward.StartDate, reward.EndDate)
	var br domain.Reward
	br.BrandId = reward.BrandId
	br.RewardName = reward.RewardName
	br.PricePoints = reward.PricePoints

	if err := row.Scan(&br.ID, &br.RewardName, &br.PricePoints); err != nil {
		return nil, err
	}

	return &br, nil
}

// GetRewardsByBrand retrieves all rewards for a given brand ID from the database.
//
// The function returns a slice of Reward objects, or an error if there is a problem
// communicating with the database. The rewards are sorted in descending order of their start
// dates.
func (r *postgresRewardRepo) GetRewardsByBrand(brandID int) ([]domain.Reward, error) {
	query := ` SELECT id, brand_id, reward_name, price_points, start_date, end_date FROM reward WHERE brand_id = $1`
	rows, err := r.db.Query(query, brandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rewards []domain.Reward
	for rows.Next() {
		var r domain.Reward
		if err := rows.Scan(&r.ID, &r.BrandId, &r.RewardName, &r.PricePoints, &r.StartDate, &r.EndDate); err != nil {
			return nil, err
		}
		rewards = append(rewards, r)
	}

	return rewards, nil
}
