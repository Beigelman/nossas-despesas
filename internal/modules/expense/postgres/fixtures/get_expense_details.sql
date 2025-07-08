-- Basic setup for get_expense_details tests
-- This will be used alongside common/basic_setup.sql

-- Basic expense without refund
INSERT INTO expenses (id, name, amount_cents, refund_amount_cents, description, group_id, category_id, payer_id, receiver_id, split_ratio, split_type, created_at, updated_at, version) VALUES
(1, 'Almoço Básico', 2500, NULL, 'Almoço no restaurante', 100, 100, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', NOW(), NOW(), 0),
(2, 'Jantar Completo', 7500, NULL, 'Jantar com todos os campos preenchidos', 100, 100, 100, 101, '{"payer": 60, "receiver": 40}', 'proportional', NOW() - INTERVAL '1 hour', NOW() - INTERVAL '1 hour', 0),
(3, 'Compra com Reembolso', 10000, 2000, 'Compra que teve reembolso parcial', 100, 101, 100, 101, '{"payer": 50, "receiver": 50}', 'equal', NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours', 0);