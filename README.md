# Loyalty Program Platform

## Architecture

The project follows a microservices architecture with two primary services:
- **Customer Service**: Manages customer-related operations
- **Brand Service**: Handles brand and campaign management

Key architectural components:
- PostgreSQL databases for data persistence
- Kafka for asynchronous inter-service communication
- Nginx as an API gateway
- Docker and Docker Compose for containerization and local development

## Prerequisites

Before getting started, ensure you have the following installed:
- Docker (version 20.10 or higher)
- Docker Compose (version 1.29 or higher)
- Git

## Project Setup

### 1. Clone the Repository

```bash
git clone https://github.com/degarzonm/loyalty-system.git
cd leal-co
```

### 2. Environment Configuration

1. Create a `.env` file in the project root directory
2. Add the following environment variables:

```env
POSTGRES_USER=your_postgres_user
POSTGRES_PASSWORD=your_postgres_password
POSTGRES_DB_CUSTOMERS=leal_customers
POSTGRES_DB_BRANDS=leal_brands
MSG_PURCHASE=purchase-topic
MSG_APPLY_POINTS=apply-points-topic
CUSTOMER_GROUP_NAME=customer-group
BRAND_GROUP_NAME=brand-group
```

### 3. Build and Run the Project

```bash
docker-compose up --build
```

This command will:
- Build the services
- Start PostgreSQL databases
- Initialize Kafka and Zookeeper
- Launch the Customer and Brand services
- Set up the Nginx API gateway

### 4. Verify Services

After startup, the following services will be available:

- **API Gateway**: `http://localhost`
- **Customer Service**: `http://localhost:8081`
- **Brand Service**: `http://localhost:8080`
- **Postgres (Customers)**: `localhost:5432`
- **Postgres (Brands)**: `localhost:5433`
- **Kafka**: `localhost:9092`

## Endpoints

### Brand Service Endpoints

#### Authentication
- `POST /new-brand`: Register a new brand
- `POST /login-brand`: Brand login

#### Branch Management
- `POST /new-branch`: Add a new branch
- `GET /my-branches`: Retrieve brand's branches

#### Campaign Management
- `POST /new-campaign`: Create a new campaign
- `POST /modify-campaign`: Update an existing campaign
- `GET /my-campaigns`: Retrieve brand's campaigns

#### Reward Management
- `POST /new-reward`: Create a new reward
- `GET /my-rewards`: Retrieve brand's rewards

### Customer Service Endpoints

#### Authentication
- `POST /new-customer`: Register a new customer
- `POST /login-customer`: Customer login

#### Points and Coins
- `GET /my-points`: View customer's points
- `GET /my-coins`: View customer's coins

#### Transactions
- `POST /purchase`: Record a purchase
- `POST /redeem`: Redeem rewards

## Example API Calls
### Brand Service Endpoints

#### 1. Ping Brand Service
```bash
curl -X GET http://localhost/ping_brands
```

#### 2. Create a New Brand
```bash
curl -X POST http://localhost/new-brand \
     -H "Content-Type: application/json" \
     -d '{
         "brand_name": "Colanta",
         "pass": "hola"
     }'
```

#### 3. Brand Login
```bash
curl -X POST http://localhost/login-brand \
     -H "Content-Type: application/json" \
     -d '{
         "brand_name": "Texaco",
         "pass": "hola"
     }'
```

#### 4. Add a New Branch
```bash
curl -X POST http://localhost/new-branch \
     -H "Leal-Brand-Id: 1" \
     -H "Leal-Brand-Token: {{brand-token}}" \
     -H "Content-Type: application/json" \
     -d '{
         "branch_name": "sucursal 5"
     }'
```

#### 5. Retrieve Brand's Branches
```bash
curl -X GET http://localhost/my-branches \
     -H "Leal-Brand-Id: 2" \
     -H "Leal-Brand-Token: {{brand-token}}"
```

#### 6. Create a New Campaign
```bash
curl -X POST http://localhost/new-campaign \
     -H "Leal-Brand-Id: 1" \
     -H "Leal-Brand-Token: {{brand-token}}" \
     -H "Content-Type: application/json" \
     -d '{
         "campaign_name": "camp_2",
         "brand_id": 1,
         "branch_ids": "2,4",
         "min_value": 10000.0,
         "max_value": 1000000.0,
         "start_date": "2024-12-15",
         "end_date": "2025-12-12",
         "status": "active",
         "point_factor": 0.5,
         "coin_factor": 1
     }'
```

#### 7. Modify an Existing Campaign
```bash
curl -X POST http://localhost/modify-campaign \
     -H "Leal-Brand-Id: 1" \
     -H "Leal-Brand-Token: {{brand-token}}" \
     -H "Content-Type: application/json" \
     -d '{
         "campaign_id": 2,
         "campaign_name": "camp_awesome",
         "brand_id": 1,
         "branch_id": "1",
         "min_value": 10000.0,
         "max_value": 100000.0,
         "start_date": "2024-12-15",
         "end_date": "2026-12-12",
         "status": "active",
         "point_factor": 0.3,
         "coin_factor": 0.7
     }'
```

#### 8. Create a New Reward
```bash
curl -X POST http://localhost/new-reward \
     -H "Leal-Brand-Id: 1" \
     -H "Leal-Brand-Token: {{brand-token}}" \
     -H "Content-Type: application/json" \
     -d '{
         "brand_id": 1,
         "reward_name": "free 2 gallon",
         "price_points": 140,
         "start_date": "2024-12-12",
         "end_date": "2025-12-12"
     }'
```

### Customer Service Endpoints

#### 1. Ping Customer Service
```bash
curl -X GET http://localhost/ping_customers
```

#### 2. Create a New Customer
```bash
curl -X POST http://localhost/new-customer \
     -H "Content-Type: application/json" \
     -d '{
         "customer_name": "juan",
         "email": "juan@leal.com",
         "phone": "3001234567",
         "pass": "hola"
     }'
```

#### 3. Customer Login
```bash
curl -X POST http://localhost/login-customer \
     -H "Content-Type: application/json" \
     -d '{
         "email": "juan@leal.com",
         "pass": "hola"
     }'
```

#### 4. Retrieve Customer Points
```bash
curl -X GET http://localhost/my-points \
     -H "Leal-Customer-Id: 1" \
     -H "Leal-Customer-Token: {{customer-token}}"
```

#### 5. Retrieve Customer Coins
```bash
curl -X GET http://localhost/my-coins \
     -H "Leal-Customer-Id: 1" \
     -H "Leal-Customer-Token: {{customer-token}}"
```

#### 6. Record a Purchase
```bash
curl -X POST http://localhost/purchase \
     -H "Leal-Customer-Id: 1" \
     -H "Leal-Customer-Token: {{customer-token}}" \
     -H "Content-Type: application/json" \
     -d '{
         "customer_id": 1,
         "amount": 500000.0,
         "brand_id": 1,
         "branch_id": 3,
         "coins_used": 50
     }'
```

#### 7. Redeem a Reward
```bash
curl -X POST http://localhost/redeem \
     -H "Leal-Customer-Id: 1" \
     -H "Leal-Customer-Token: {{customer-token}}" \
     -H "Content-Type: application/json" \
     -d '{
         "customer_id": 1,
         "brand_id": 1,
         "reward_id": 1,
         "points_spend": 100
     }'
```

## Notes on API Calls

- Replace `{{brand-token}}` and `{{customer-token}}` with actual tokens received during login
- Ensure the `Leal-Brand-Id` or `Leal-Customer-Id` matches the ID of the authenticated user
- All endpoints require appropriate authentication headers
- Timestamps should be in ISO 8601 format (YYYY-MM-DD)
 
## Project Considerations

- Points and coins are managed separately
- Campaigns support flexible configurations (date ranges, branch selection, purchase value thresholds)
- Base campaigns can be modified
- Purchases trigger point and coin calculations based on brand-specific rules

 
## License

Distributed under the MIT License. See `LICENSE` for more information.

## Contact

Project Link: [https://github.com/degarzonm/loyalty-system](https://github.com/degarzonm/loyalty_system)
