CREATE TABLE IF NOT EXISTS brand (
    id SERIAL PRIMARY KEY,
    brand_name VARCHAR(100) NOT NULL UNIQUE,
    pass_hash VARCHAR(255),
    token VARCHAR(255),
    registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS branch (
    id SERIAL PRIMARY KEY,
    brand_id INT NOT NULL REFERENCES brand(id),
    branch_name VARCHAR(100) NOT NULL,
    registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_branch_per_brand UNIQUE (brand_id, branch_name)
);

CREATE TABLE IF NOT EXISTS campaign (
    id SERIAL PRIMARY KEY,
    campaign_name VARCHAR(100) NOT NULL,
    brand_id INT REFERENCES brand(id),
    min_value DECIMAL(20, 2),
    max_value DECIMAL(20, 2),
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    status VARCHAR(30),
    point_factor DECIMAL(10, 4),
    coin_factor DECIMAL(10, 4),
    customer_count INT DEFAULT 0,
    CONSTRAINT unique_campaign_per_brand UNIQUE (brand_id, campaign_name)
);

CREATE TABLE IF NOT EXISTS campaign_branches (
    campaign_id INT NOT NULL REFERENCES campaign(id),
    branch_id INT NOT NULL REFERENCES branch(id),
    PRIMARY KEY (campaign_id, branch_id)
);

CREATE TABLE IF NOT EXISTS reward (
    id SERIAL PRIMARY KEY,
    brand_id INT NOT NULL REFERENCES brand(id),
    reward_name VARCHAR(100) NOT NULL,
    price_points INT NOT NULL,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    CONSTRAINT unique_reward_per_brand UNIQUE (brand_id, reward_name)
);

CREATE INDEX idx_brand_name ON brand(brand_name);

CREATE INDEX idx_branch_brand_id ON branch(brand_id);

CREATE INDEX idx_branch_name ON branch(branch_name);

CREATE INDEX idx_campaign_brand_id_start_date ON campaign(brand_id, start_date);

CREATE INDEX idx_campaign_status ON campaign(status);

CREATE INDEX idx_campaign_start_end_date ON campaign(start_date, end_date);

CREATE INDEX idx_campaign_brand_id_min_max_value ON campaign(brand_id, min_value, max_value);

CREATE INDEX idx_campaign_branches_campaign_id_branch_id ON campaign_branches(campaign_id, branch_id);

CREATE INDEX idx_reward_brand_id ON reward(brand_id);

CREATE INDEX idx_reward_start_end_date ON reward(start_date, end_date);