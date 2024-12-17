package application

import (
	"errors"

	"github.com/degarzonm/customer_leal_service/internal/domain"
	"github.com/degarzonm/customer_leal_service/internal/infrastructure/util"
)

type customerService struct {
	customerRepo domain.CustomerRepository
}

func NewCustomerService(p domain.CustomerRepository) domain.CustomerService {
	return &customerService{customerRepo: p}
}

// CreateCustomer creates a new customer and returns the customer if successful.
// The customer name, email, phone and password are required.
// An error is returned if any of the required fields are not provided.
// The customer's password is hashed before being stored.
func (c *customerService) CreateCustomer(name string, email string, phone string, pass string) (*domain.Customer, error) {

	if name == "" || pass == "" || email == "" || phone == "" {
		return nil, errors.New("name, email, phone and password are required")
	}
	passHash := util.HashPassword(pass)
	token, err := util.GenerateToken()

	if err != nil {
		return nil, errors.New("error generating token")
	}
	return c.customerRepo.CreateCustomer(name, email, phone, passHash, token)
}

// LoginCustomer authenticates a customer using their email and password.
// It returns the customer object with a refreshed authentication token if successful.
// An error is returned if the customer is not found, the password is invalid,
// or if there is an error generating or updating the token.

func (c *customerService) LoginCustomer(email string, pass string) (*domain.Customer, error) {
	b, err := c.customerRepo.GetCustomerByEmail(email)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, errors.New("brand not found")
	}

	if !util.CheckPassHash(pass, b.PassHash) {
		return nil, errors.New("invalid password")
	}
	// update token
	newToken, err := util.GenerateToken()
	if err != nil {
		return nil, errors.New("failed generating token")
	}

	err = c.customerRepo.UpdateCustomerToken(b.ID, newToken)
	if err != nil {
		return nil, err
	}
	b.Token = newToken
	return b, nil
}

// GetCustomerByID retrieves a customer by ID.
// It returns the customer if found, or an error if the customer is not found.
func (c *customerService) GetCustomerByID(id int) (*domain.Customer, error) {
	return c.customerRepo.GetCustomerByID(id)
}

// ValidateToken checks if the provided token matches the stored token for the customer with the given customerID.
// It returns an error if the customer is not found or if the token is invalid.
func (c *customerService) ValidateToken(customerID int, token string) error {
	b, err := c.customerRepo.GetCustomerByID(customerID)
	if err != nil {
		return err
	}
	if b == nil {
		return errors.New("customer not found")
	}
	if b.Token != token {
		return errors.New("invalid token")
	}
	return nil
}
