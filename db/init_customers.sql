CREATE TABLE IF NOT EXISTS customer (
    id SERIAL PRIMARY KEY,
    customer_name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    phone VARCHAR(20) UNIQUE NOT NULL,
    pass_hash VARCHAR(255),
    token VARCHAR(255),
    leal_coins INT DEFAULT 0,
    registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS purchase (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL REFERENCES customer(id),
    amount DECIMAL(10, 2) NOT NULL,
    purchase_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    brand_id INT,
    branch_id INT,
    coins_used INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS leal_points (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL REFERENCES customer(id),
    brand_id INT NOT NULL,
    points INT DEFAULT 0,
    UNIQUE (customer_id, brand_id)
);

CREATE TABLE IF NOT EXISTS leal_points_transactions (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL REFERENCES customer(id),
    brand_id INT NOT NULL,
    change INT NOT NULL,
    reason VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS redeemed (
    id SERIAL PRIMARY KEY,
    customer_id INT,
    brand_id INT NOT NULL,
    reward_id INT NOT NULL,
    points_spend INT NOT NULL,
    date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_leal_points_customer_id ON leal_points(customer_id);

CREATE INDEX idx_leal_points_transactions_customer_id ON leal_points_transactions(customer_id);

CREATE INDEX idx_leal_points_transactions_brand_id ON leal_points_transactions(brand_id);

CREATE INDEX idx_purchase_customer_id ON purchase(customer_id);

CREATE INDEX idx_purchase_brand_id ON purchase(brand_id);

CREATE INDEX idx_purchase_branch_id ON purchase(branch_id);

CREATE INDEX idx_redeemed_customer_id_brand_id_reward_id ON redeemed(customer_id, brand_id, reward_id);

CREATE INDEX idx_redeemed_date ON redeemed(date);