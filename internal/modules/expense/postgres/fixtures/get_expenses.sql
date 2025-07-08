-- Multiple expenses for get_expenses tests
-- These expenses have different dates for testing pagination and ordering

INSERT INTO expenses (id, name, amount_cents, description, group_id, category_id, payer_id, receiver_id, split_ratio, split_type, created_at, updated_at, version) VALUES
(20, 'Primeira Despesa', 1000, 'Primeira despesa para teste', 100, 100, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-01-01 10:00:00', '2024-01-01 10:00:00', 0),
(21, 'Segunda Despesa', 2000, 'Segunda despesa para teste', 100, 100, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-01-02 10:00:00', '2024-01-02 10:00:00', 0),
(22, 'Terceira Despesa', 3000, 'Terceira despesa para teste', 100, 101, 100, 101, '{"payer": 60, "receiver": 40}', 'proportional', '2024-01-03 10:00:00', '2024-01-03 10:00:00', 0),
(23, 'Quarta Despesa', 4000, 'Quarta despesa para teste', 100, 101, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-01-04 10:00:00', '2024-01-04 10:00:00', 0),
(24, 'Quinta Despesa', 5000, 'Quinta despesa para teste', 100, 102, 100, 101, '{"payer": 70, "receiver": 30}', 'proportional', '2024-01-05 10:00:00', '2024-01-05 10:00:00', 0),
-- Expenses for another group (should not appear in group 100 results)
(25, 'Despesa Grupo 101', 6000, 'Despesa do grupo 101', 101, 100, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-01-06 10:00:00', '2024-01-06 10:00:00', 0),
-- Expenses for testing search
(30, 'Almoço McDonald', 2500, 'Almoço no McDonald para teste de busca', 100, 100, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-01-10 12:00:00', '2024-01-10 12:00:00', 0),
(31, 'Jantar Pizza', 4500, 'Jantar com pizza no restaurante', 100, 100, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-01-10 19:00:00', '2024-01-10 19:00:00', 0),
(32, 'Supermercado Compras', 8000, 'Compras no supermercado para casa', 100, 101, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-01-11 10:00:00', '2024-01-11 10:00:00', 0),
(33, 'Cinema Filme', 3000, 'Assistir filme no cinema', 100, 103, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-01-11 20:00:00', '2024-01-11 20:00:00', 0),
(34, 'Uber Transporte', 1500, 'Transporte de uber para trabalho', 100, 102, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-01-12 08:00:00', '2024-01-12 08:00:00', 0),
-- Expenses with refund amounts for testing
(40, 'Compra com Reembolso Total', 10000, 'Compra que foi totalmente reembolsada', 100, 101, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-01-15 10:00:00', '2024-01-15 10:00:00', 0),
(41, 'Compra com Reembolso Parcial', 8000, 'Compra com reembolso parcial', 100, 101, 100, 101, '{"payer": 60, "receiver": 40}', 'proportional', '2024-01-16 10:00:00', '2024-01-16 10:00:00', 0),
(42, 'Jantar com Reembolso', 5000, 'Jantar que teve reembolso', 100, 100, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', '2024-01-17 19:00:00', '2024-01-17 19:00:00', 0);

-- Update expenses with refund amounts
UPDATE expenses SET refund_amount_cents = 10000 WHERE id = 40;
UPDATE expenses SET refund_amount_cents = 3000 WHERE id = 41;
UPDATE expenses SET refund_amount_cents = 2000 WHERE id = 42; 