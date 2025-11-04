-- Expenses organized by categories for get_expenses_per_category tests

-- Expenses in Alimentação category group
INSERT INTO expenses (id, name, amount_cents, description, group_id, category_id, payer_id, receiver_id, split_ratio, split_type, created_at, updated_at, version) VALUES
(50, 'Almoço Restaurante 1', 2500, 'Almoço no restaurante', 100, 100, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-06-01 12:00:00', '2024-06-01 12:00:00', 0),
(51, 'Almoço Restaurante 2', 3000, 'Outro almoço no restaurante', 100, 100, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-06-02 12:00:00', '2024-06-02 12:00:00', 0),
(52, 'Compras Supermercado 1', 8000, 'Compras no supermercado', 100, 101, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-06-03 10:00:00', '2024-06-03 10:00:00', 0),
(53, 'Compras Supermercado 2', 12000, 'Mais compras no supermercado', 100, 101, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-06-04 10:00:00', '2024-06-04 10:00:00', 0);

-- Expenses in Transporte category group
INSERT INTO expenses (id, name, amount_cents, description, group_id, category_id, payer_id, receiver_id, split_ratio, split_type, created_at, updated_at, version) VALUES
(54, 'Uber 1', 1500, 'Transporte de uber', 100, 102, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-06-05 08:00:00', '2024-06-05 08:00:00', 0),
(55, 'Uber 2', 2000, 'Outro transporte de uber', 100, 102, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-06-06 18:00:00', '2024-06-06 18:00:00', 0);

-- Expenses in Lazer category group
INSERT INTO expenses (id, name, amount_cents, description, group_id, category_id, payer_id, receiver_id, split_ratio, split_type, created_at, updated_at, version) VALUES
(56, 'Cinema', 3000, 'Assistir filme no cinema', 100, 103, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-06-07 20:00:00', '2024-06-07 20:00:00', 0);

-- Expenses outside the date range (should not appear in results)
INSERT INTO expenses (id, name, amount_cents, description, group_id, category_id, payer_id, receiver_id, split_ratio, split_type, created_at, updated_at, version) VALUES
(57, 'Despesa Fora do Período', 5000, 'Despesa anterior ao período', 100, 100, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-05-01 12:00:00', '2024-05-01 12:00:00', 0);

-- Expenses for another group (should not appear in results)
INSERT INTO expenses (id, name, amount_cents, description, group_id, category_id, payer_id, receiver_id, split_ratio, split_type, created_at, updated_at, version) VALUES
(58, 'Despesa Outro Grupo', 4000, 'Despesa de outro grupo', 101, 100, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-06-08 12:00:00', '2024-06-08 12:00:00', 0); 