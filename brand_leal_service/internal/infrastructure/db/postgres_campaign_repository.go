package db

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	"github.com/degarzonm/brand_leal_service/internal/domain"
)

type postgresCampaignRepo struct {
	db *sql.DB
}

func NewPostgresCampaignRepo(db *sql.DB) domain.CampaignRepository {
	return &postgresCampaignRepo{db: db}
}

// CreateCampaign inserts a new campaign into the database along with its associated branches.
// It first inserts the campaign details into the campaign table, and retrieves the generated campaign ID.
// If the campaign name is "base", it skips the insertion into the campaign_branches table.
// Otherwise, it inserts the given branch IDs into the campaign_branches table, linking them with the campaign ID.
// Returns the created campaign with its ID filled or an error if something goes wrong.

func (r *postgresCampaignRepo) CreateCampaign(c *domain.Campaign, branchIDs []int) (*domain.Campaign, error) {
	log.Println("Creating campaign:", c)
	query := `INSERT INTO campaign (campaign_name, brand_id, min_value, max_value, start_date, end_date, status, point_factor, coin_factor)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	row := r.db.QueryRow(query, c.CampaignName, c.BrandID, c.MinValue, c.MaxValue, c.StartDate, c.EndDate, c.Status, c.PointFactor, c.CoinFactor)

	if err := row.Scan(&c.ID); err != nil {
		return nil, err
	}

	// If campaign_name is "base", skip inserting into campaign_branches
	if c.CampaignName == "base" {
		return c, nil
	}

	// Insert into campaign_branches for all branch IDs
	for _, bid := range branchIDs {
		_, err := r.db.Exec(`INSERT INTO campaign_branches (campaign_id, branch_id) VALUES ($1, $2)`, c.ID, bid)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

// GetCampaignByID returns a campaign by its ID or nil if the campaign does not exist.
func (r *postgresCampaignRepo) GetCampaignByID(id int) (*domain.Campaign, error) {
	query := `SELECT campaign_name, brand_id, min_value, max_value, start_date, end_date, status, point_factor, coin_factor, customer_count FROM campaign WHERE id=$1`
	row := r.db.QueryRow(query, id)
	var c domain.Campaign
	c.ID = id
	if err := row.Scan(&c.CampaignName, &c.BrandID, &c.MinValue, &c.MaxValue, &c.StartDate, &c.EndDate, &c.Status, &c.PointFactor, &c.CoinFactor, &c.CustomerCount); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

// UpdateCampaign updates an existing campaign in the database with the given campaign and branch IDs.
// It first updates the campaign details in the campaign table.
// Then it deletes the existing branch IDs associated with the campaign from the campaign_branches table.
// Finally it inserts the given branch IDs into the campaign_branches table, linking them with the campaign ID.
// Returns an error if something goes wrong.
func (r *postgresCampaignRepo) UpdateCampaign(c *domain.Campaign, branchIDs []int) error {
	query := `UPDATE campaign SET campaign_name=$1, min_value=$2, max_value=$3, start_date=$4, end_date=$5, point_factor=$6, coin_factor=$7 WHERE id=$8 AND brand_id=$9`
	_, err := r.db.Exec(query, c.CampaignName, c.MinValue, c.MaxValue, c.StartDate, c.EndDate, c.PointFactor, c.CoinFactor, c.ID, c.BrandID)
	if err != nil {
		return err
	}

	// first delete the existing branch IDs
	_, err = r.db.Exec(`DELETE FROM campaign_branches WHERE campaign_id=$1`, c.ID)
	if err != nil {
		return err
	}

	// Insert the new branch IDs
	for _, bid := range branchIDs {
		_, err := r.db.Exec(`INSERT INTO campaign_branches (campaign_id, branch_id) VALUES ($1, $2)`, c.ID, bid)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetCampaignsByBrandID returns a list of all campaigns for the given brand ID and their associated branch IDs.
// The branch IDs are returned as a comma-separated string in the "branch_ids" field.
// If a campaign does not have any branch IDs, the "branch_ids" field will be an empty string.
// The campaigns are returned in descending order of start date.
func (r *postgresCampaignRepo) GetCampaignsByBrandID(brandID int) ([]domain.Campaign, error) {
	query := `
		SELECT 
			c.id, c.campaign_name, c.min_value, c.max_value, c.start_date, c.end_date, 
			c.status, c.point_factor, c.coin_factor, c.customer_count,
			COALESCE(STRING_AGG(cb.branch_id::TEXT, ','), '') AS branch_ids
		FROM 
			campaign c
		LEFT JOIN 
			campaign_branches cb ON c.id = cb.campaign_id
		WHERE 
			c.brand_id = $1
		GROUP BY 
			c.id
		ORDER BY 
			c.start_date DESC`

	rows, err := r.db.Query(query, brandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var campaigns []domain.Campaign
	for rows.Next() {
		var c domain.Campaign
		var branchIDs string

		if err := rows.Scan(
			&c.ID, &c.CampaignName, &c.MinValue, &c.MaxValue,
			&c.StartDate, &c.EndDate, &c.Status, &c.PointFactor,
			&c.CoinFactor, &c.CustomerCount, &branchIDs,
		); err != nil {
			return nil, err
		}

		// Convert branchIDs string into a slice of integers
		if branchIDs != "" {
			for _, id := range strings.Split(branchIDs, ",") {
				branchID, err := strconv.Atoi(id)
				if err != nil {
					return nil, err
				}
				c.Branches = append(c.Branches, branchID)
			}
		}

		c.BrandID = brandID
		campaigns = append(campaigns, c)
	}
	return campaigns, nil
}

// GetBranchesForCampaign returns a list of all branches associated with the given campaign ID.
func (r *postgresCampaignRepo) GetBranchesForCampaign(campaignID int) ([]domain.Branch, error) {
	query := `SELECT b.id, b.brand_id, b.branch_name, b.registration_date 
		FROM campaign_branches cb
		INNER JOIN branch b ON cb.branch_id = b.id
		WHERE cb.campaign_id=$1`
	rows, err := r.db.Query(query, campaignID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var branches []domain.Branch
	for rows.Next() {
		var br domain.Branch
		if err := rows.Scan(&br.ID, &br.BrandID, &br.Name, &br.RegistrationDate); err != nil {
			return nil, err
		}
		branches = append(branches, br)
	}
	return branches, nil
}

// GetCampaignsForBranch retrieves all active campaigns associated with a specific branch ID.
// It queries the campaign_branches table to find campaigns linked to the given branch ID
// and joins with the campaign table to obtain campaign details. Only campaigns with an "active"
// status are retrieved. It returns a slice of Campaign objects or an error if the query fails.

func (r *postgresCampaignRepo) GetCampaignsForBranch(branchID int) ([]domain.Campaign, error) {
	query := `SELECT c.id, c.campaign_name, c.brand_id, c.min_value, c.max_value, 
	                 c.start_date, c.end_date, c.status, c.point_factor, c.coin_factor, c.customer_count
	          FROM campaign_branches cb
	          INNER JOIN campaign c ON cb.campaign_id = c.id
	          WHERE cb.branch_id = $1
			  AND c.status = $2`
	rows, err := r.db.Query(query, branchID, "active")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var campaigns []domain.Campaign
	for rows.Next() {
		var camp domain.Campaign
		if err := rows.Scan(&camp.ID, &camp.CampaignName, &camp.BrandID, &camp.MinValue, &camp.MaxValue,
			&camp.StartDate, &camp.EndDate, &camp.Status, &camp.PointFactor, &camp.CoinFactor, &camp.CustomerCount); err != nil {
			return nil, err
		}
		campaigns = append(campaigns, camp)
	}
	return campaigns, nil
}

// GetBaseCampaignForBrand retrieves the base campaign for a specific brand ID from the database.
// It executes a query to select a campaign with the name 'base' and the specified brand ID.
// If a base campaign is found, it returns a pointer to the Campaign object.
// If no base campaign is found, it returns nil.
// If any error occurs during the query execution, it returns the error.

func (r *postgresCampaignRepo) GetBaseCampaignForBrand(brandID int) (*domain.Campaign, error) {
	query := `SELECT c.id, c.campaign_name, c.brand_id, c.min_value, c.max_value, 
	                 c.start_date, c.end_date, c.status, c.point_factor, c.coin_factor, c.customer_count
	          FROM campaign c
	          WHERE c.campaign_name = 'base' AND c.brand_id = $1`

	row := r.db.QueryRow(query, brandID)

	var camp domain.Campaign
	err := row.Scan(&camp.ID, &camp.CampaignName, &camp.BrandID, &camp.MinValue, &camp.MaxValue,
		&camp.StartDate, &camp.EndDate, &camp.Status, &camp.PointFactor, &camp.CoinFactor, &camp.CustomerCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No base campaign found
		}
		return nil, err
	}
	return &camp, nil
}

// UpdateCustomerCountCampaign increments the customer count for a specified campaign by 1 in the database.
// It executes an update query on the campaign table using the campaign ID provided.
// Returns an error if the update query fails, otherwise returns nil.

func (r *postgresCampaignRepo) UpdateCustomerCountCampaign(c *domain.Campaign) error {
	query := `UPDATE campaign SET customer_count = customer_count+1  WHERE id = $1`
	_, err := r.db.Exec(query, c.ID)
	return err
}
