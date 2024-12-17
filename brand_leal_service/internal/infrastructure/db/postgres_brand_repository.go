package db

import (
	"database/sql"
	"github.com/degarzonm/brand_leal_service/internal/domain"
)

type postgresBrandRepo struct {
	db *sql.DB
}

func NewPostgresBrandRepo(db *sql.DB) domain.BrandsRepository {
	return &postgresBrandRepo{db: db}
}

// GetBrandByID obtains a brand by its id , if it does not exist it returns nil
func (r *postgresBrandRepo) GetBrandByID(id int) (*domain.Brand, error) {
	query := `SELECT id, brand_name, token, registration_date FROM brand WHERE id = $1`
	row := r.db.QueryRow(query, id)
	var b domain.Brand
	if err := row.Scan(&b.ID, &b.Name, &b.Token, &b.RegistrationDate); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

// GetBrandByName retrieves a brand by its name from the database.
// It returns the brand object if found, or nil if no brand with the given name exists.
// If an error occurs during the query execution or scanning, it returns the error.

func (r *postgresBrandRepo) GetBrandByName(brandName string) (*domain.Brand, error) {
	query := `SELECT id, brand_name,pass_hash, token, registration_date FROM brand WHERE brand_name = $1`
	row := r.db.QueryRow(query, brandName)
	var b domain.Brand
	if err := row.Scan(&b.ID, &b.Name, &b.PassHash, &b.Token, &b.RegistrationDate); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

// GetBrandByNameAndPass retrieves a brand from the database using the brand name and password hash.
// It returns the brand object if a matching brand is found, or nil if no such brand exists.
// If an error occurs during the query execution or scanning process, it returns the error.

func (r *postgresBrandRepo) GetBrandByNameAndPass(brandName, passHash string) (*domain.Brand, error) {
	query := `SELECT id, brand_name, token, registration_date FROM brand WHERE brand_name = $1 AND pass_hash = $2`
	row := r.db.QueryRow(query, brandName, passHash)
	var b domain.Brand
	if err := row.Scan(&b.ID, &b.Name, &b.Token, &b.RegistrationDate); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

// CreateBrand creates a new brand in the database and returns the brand object if it was created successfully,
// or an error if there was an issue during the creation process. The brandName, passHash, and token parameters
// are used to initialize the brand object. The returned brand object will have the id, name, registration_date, and token fields populated.
// If the creation fails due to an error or if the brand already exists, the function returns nil and the error.
func (r *postgresBrandRepo) CreateBrand(brandName, passHash, token string) (*domain.Brand, error) {
	query := `INSERT INTO brand (brand_name, pass_hash, token) VALUES ($1, $2, $3) RETURNING id, registration_date`
	var b domain.Brand
	b.Name = brandName
	b.Token = token
	row := r.db.QueryRow(query, brandName, passHash, token)
	if err := row.Scan(&b.ID, &b.RegistrationDate); err != nil {
		return nil, err
	}
	return &b, nil
}

// UpdateBrandToken updates the token of a brand in the database identified by the given brandID.
// It returns an error if the update operation fails.

func (r *postgresBrandRepo) UpdateBrandToken(brandID int, token string) error {
	query := `UPDATE brand SET token=$1 WHERE id=$2`
	_, err := r.db.Exec(query, token, brandID)
	return err
}
