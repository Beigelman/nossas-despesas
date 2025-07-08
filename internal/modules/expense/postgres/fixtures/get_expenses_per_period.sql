-- Expenses organized by periods for get_expenses_per_period tests

-- June 2024 expenses - different days
INSERT INTO expenses (id, name, amount_cents, description, group_id, category_id, payer_id, receiver_id, split_ratio, split_type, created_at, updated_at, version) VALUES
(60, 'Despesa 01 Junho', 1000, 'Primeira despesa de junho', 100, 100, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-06-01 10:00:00', '2024-06-01 10:00:00', 0),
(61, 'Despesa 01 Junho B', 1500, 'Segunda despesa do dia 01', 100, 100, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-06-01 15:00:00', '2024-06-01 15:00:00', 0),
(62, 'Despesa 02 Junho', 2000, 'Despesa do dia 02', 100, 101, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-06-02 12:00:00', '2024-06-02 12:00:00', 0),
(63, 'Despesa 15 Junho', 3000, 'Despesa do meio do mês', 100, 102, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-06-15 14:00:00', '2024-06-15 14:00:00', 0),
(64, 'Despesa 30 Junho', 2500, 'Última despesa de junho', 100, 103, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-06-30 18:00:00', '2024-06-30 18:00:00', 0);

-- July 2024 expenses - different days
INSERT INTO expenses (id, name, amount_cents, description, group_id, category_id, payer_id, receiver_id, split_ratio, split_type, created_at, updated_at, version) VALUES
(65, 'Despesa 01 Julho', 4000, 'Primeira despesa de julho', 100, 100, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-07-01 09:00:00', '2024-07-01 09:00:00', 0),
(66, 'Despesa 15 Julho', 3500, 'Despesa do meio de julho', 100, 101, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-07-15 16:00:00', '2024-07-15 16:00:00', 0),
(67, 'Despesa 31 Julho', 5000, 'Última despesa de julho', 100, 102, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-07-31 20:00:00', '2024-07-31 20:00:00', 0);

-- August 2024 expenses
INSERT INTO expenses (id, name, amount_cents, description, group_id, category_id, payer_id, receiver_id, split_ratio, split_type, created_at, updated_at, version) VALUES
(68, 'Despesa 10 Agosto', 6000, 'Despesa de agosto', 100, 103, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-08-10 11:00:00', '2024-08-10 11:00:00', 0);

-- Expenses for another group (should not appear in results)
INSERT INTO expenses (id, name, amount_cents, description, group_id, category_id, payer_id, receiver_id, split_ratio, split_type, created_at, updated_at, version) VALUES
(69, 'Despesa Outro Grupo', 7000, 'Despesa de outro grupo', 101, 100, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-06-15 12:00:00', '2024-06-15 12:00:00', 0); 