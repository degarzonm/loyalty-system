package db

import (
	"database/sql"

	"github.com/degarzonm/customer_leal_service/internal/domain"
)

type postgresCustomerRepo struct {
	db *sql.DB
}

func NewPostgresCustomerRepo(db *sql.DB) domain.CustomerRepository {
	return &postgresCustomerRepo{db: db}
}

// CreateCustomer creates a new customer and returns the customer if successful.
// The customer name, email, phone, password and token are required.
// An error is returned if any of the required fields are not provided.
// The customer's password is hashed before being stored.
func (r *postgresCustomerRepo) CreateCustomer(name string, email string, phone string, pass string, token string) (*domain.Customer, error) {
	query := `INSERT INTO customer (customer_name, email, phone, pass_hash, token, leal_coins) VALUES ($1, $2, $3, $4 , $5 , 0) RETURNING id`
	row := r.db.QueryRow(query, name, email, phone, pass, token)
	var c domain.Customer
	c.Name = name
	c.Email = email
	c.Phone = phone
	c.Token = token
	if err := row.Scan(&c.ID); err != nil {
		return nil, err
	}
	return &c, nil

}

// GetCustomerByID retrieves a customer by ID.
// It returns the customer if found, or an error if the customer is not found.
func (r *postgresCustomerRepo) GetCustomerByID(id int) (*domain.Customer, error) {
	query := `SELECT id, customer_name, email, phone ,pass_hash, token, leal_coins FROM customer WHERE id = $1`
	row := r.db.QueryRow(query, id)
	var c domain.Customer
	if err := row.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.PassHash, &c.Token, &c.LealCoins); err != nil {
		return nil, err
	}
	return &c, nil
}

// GetCustomerByEmail retrieves a customer by email.
// It returns the customer if found, or an error if the customer is not found.
func (r *postgresCustomerRepo) GetCustomerByEmail(email string) (*domain.Customer, error) {
	query := `SELECT id, customer_name, email, phone ,pass_hash, token, leal_coins FROM customer WHERE email = $1`
	row := r.db.QueryRow(query, email)
	var c domain.Customer
	if err := row.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.PassHash, &c.Token, &c.LealCoins); err != nil {
		return nil, err
	}
	return &c, nil

}

// UpdateCustomerToken updates the token for the customer with the given ID.
// It returns an error if the query fails.
func (r *postgresCustomerRepo) UpdateCustomerToken(id int, token string) error {
	query := `UPDATE customer SET token = $1 WHERE id = $2`
	_, err := r.db.Exec(query, token, id)
	return err
}
