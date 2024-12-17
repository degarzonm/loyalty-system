package db

import (
	"database/sql"
	"errors"

	"github.com/degarzonm/brand_leal_service/internal/domain"
)

type postgresBranchRepo struct {
	db *sql.DB
}

func NewPostgresBranchRepo(db *sql.DB) domain.BranchesRepository {
	return &postgresBranchRepo{db: db}
}

// CreateBranch creates a new branch for a given brand_id and returns the newly created branch if successful.
// The returned branch includes the generated id and registration_date.
func (r *postgresBranchRepo) CreateBranch(brandID int, branchName string) (*domain.Branch, error) {
	query := `INSERT INTO branch (brand_id, branch_name) VALUES ($1, $2) RETURNING id, registration_date`
	row := r.db.QueryRow(query, brandID, branchName)
	var br domain.Branch
	br.BrandID = brandID
	br.Name = branchName
	if err := row.Scan(&br.ID, &br.RegistrationDate); err != nil {
		return nil, err
	}
	return &br, nil
}

// LinkBranchToBaseCampaign links a branch to the base campaign for a given brand_id by inserting
// a record into campaign_branches. It first queries the campaign table to retrieve the id of the
// base campaign for the given brand_id. If the base campaign is not found, it returns an error.
// If the base campaign is found, it inserts a record into campaign_branches with the campaign_id
// and branch_id. If the insert fails, it returns an error. Otherwise, it returns nil.
func (r *postgresBranchRepo) LinkBranchToBaseCampaign(branch *domain.Branch) error {
	// Query to find the base campaign for the brand
	campaignQuery := `
		SELECT id
		FROM campaign
		WHERE brand_id = $1 AND campaign_name = 'base'
		LIMIT 1`

	var campaignID int

	// Step 1: Retrieve the base campaign ID
	err := r.db.QueryRow(campaignQuery, branch.BrandID).Scan(&campaignID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("base campaign not found for the given brand_id")
		}
		return err
	}

	// Step 2: Insert into campaign_branches
	insertQuery := `
		INSERT INTO campaign_branches (campaign_id, branch_id)
		VALUES ($1, $2)`

	_, err = r.db.Exec(insertQuery, campaignID, branch.ID)
	if err != nil {
		return err
	}

	return nil
}

// GetBranchesByBrandID retrieves a list of branches associated with a given brand_id.
// It executes a SQL query that filters branches by the given brand_id and returns the
// results as a slice of domain.Branch objects. If the query fails, it returns an error.
// Otherwise, it returns the list of branches.
func (r *postgresBranchRepo) GetBranchesByBrandID(brandID int) ([]domain.Branch, error) {
	query := `SELECT id, branch_name, registration_date FROM branch WHERE brand_id = $1`
	rows, err := r.db.Query(query, brandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var branches []domain.Branch
	for rows.Next() {
		var br domain.Branch
		br.BrandID = brandID
		if err := rows.Scan(&br.ID, &br.Name, &br.RegistrationDate); err != nil {
			return nil, err
		}
		branches = append(branches, br)
	}
	return branches, nil
}
