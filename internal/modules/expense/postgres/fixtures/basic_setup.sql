-- Basic setup for expense tests
-- Users
INSERT INTO users (id, name, email, created_at, updated_at, version) VALUES
(100, 'Payer User', 'payer@test.com', NOW(), NOW(), 0),
(101, 'Receiver User', 'receiver@test.com', NOW(), NOW(), 0);

-- Category Groups
INSERT INTO category_groups (id, name, icon, created_at, updated_at, version) VALUES
(100, 'Alimentação', 'food', NOW(), NOW(), 0),
(101, 'Transporte', 'transport', NOW(), NOW(), 0),
(102, 'Lazer', 'entertainment', NOW(), NOW(), 0);

-- Categories
INSERT INTO categories (id, name, icon, category_group_id, created_at, updated_at, version) VALUES
(100, 'Restaurante', 'restaurant', 100, NOW(), NOW(), 0),
(101, 'Supermercado', 'grocery', 100, NOW(), NOW(), 0),
(102, 'Uber', 'car', 101, NOW(), NOW(), 0),
(103, 'Cinema', 'movie', 102, NOW(), NOW(), 0);

-- Groups
INSERT INTO groups (id, name, created_at, updated_at, version) VALUES
(100, 'Test Group', NOW(), NOW(), 0),
(101, 'Another Group', NOW(), NOW(), 0); 