BEGIN;

INSERT INTO category_groups (name, icon, created_at, updated_at, deleted_at, version)
VALUES  ('Balanço', 'Landmark', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0);

INSERT INTO categories (name, icon, created_at, updated_at, deleted_at, version, category_group_id)
VALUES
    -- Casa
    ('Utensílios', 'PocketKnife', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
    -- Comidas e Bebidas
    ('Lanches', 'Sandwich', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 2),
    -- Vida
    ('Seguro', 'Shield', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 5),
    -- Balanço
    ('Reembolso', 'Receipt', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 9),
    ('Cashback', 'IterationCw', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 9),
    ('Outros', 'Landmark', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 9);

COMMIT;