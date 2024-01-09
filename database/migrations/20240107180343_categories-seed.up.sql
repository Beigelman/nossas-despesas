BEGIN;

INSERT INTO category_groups (id, name, icon, created_at, updated_at, deleted_at, version)
VALUES  (1, 'Casa', '1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0),
        (2, 'Comidas e Bebidas', '2', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0),
        (3, 'Entretenimento', '3', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0),
        (4, 'Transporte', '4', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0),
        (5, 'Vida', '5', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0),
        (6, 'Geral', '6', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0),
        (7, 'Viagem', '6', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0);

INSERT INTO categories (id, name, icon, created_at, updated_at, deleted_at, version, category_group_id)
-- Casa
VALUES  (1, 'Aluguel', '1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        (2, 'Condomínio', '2', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        (3, 'Água', '3', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        (4, 'Luz', '4', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        (5, 'Internet', '5', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        (6, 'Manutenção', '6', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        (7, 'Animal de estimação', '7', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        (8, 'Eletrodomésticos', '8', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        (9, 'Móveis', '9', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        (10, 'Serviços', '10', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        (11, 'Produtos de limpeza', '11', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        (12, 'Outros', '12', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
        -- Comidas e Bebidas
        (13, 'Supermercado', '13', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 2),
        (14, 'Restaurantes', '14', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 2),
        (15, 'Delivery', '15', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 2),
        (16, 'Bebidas', '16', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 2),
        (17, 'Outros', '17', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
        -- Entretenimento
        (18, 'Cinema', '18', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
        (19, 'Shows', '19', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
        (20, 'Jogos', '20', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
        (21, 'Livros', '21', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
        (22, 'Música', '22', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
        (23, 'Streaming', '23', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
        (24, 'Festas', '24', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
        (25, 'Outros', '25', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
        -- Transporte
        (26, 'Combustível', '26', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
        (27, 'Estacionamento', '27', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
        (28, 'Manutenção', '28', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
        (29, 'Transporte público', '29', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
        (30, 'Aplicativos de transporte', '30', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
        (31, 'Pedágio', '31', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
        (32, 'Aluguel de veículos', '32', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
        (33, 'Seguro', '33', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
        (34, 'Outros', '34', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
        -- Vida
        (35, 'Academia', '35', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 5),
        (36, 'Cuidados pessoais', '36', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 5),
        (37, 'Saúde', '37', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 5),
        (38, 'Farmácia', '38', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 5),
        (39, 'Plano odontológico', '39', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 5),
        (40, 'Plano de saúde', '40', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 5),
        (41, 'Outros', '41', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
        -- Geral
        (42, 'Impostos', '42', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
        (43, 'Outros', '43', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
        -- Viagem
        (44, 'Passagens', '44', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 7),
        (45, 'Hospedagem', '45', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 7),
        (46, 'Alimentação', '46', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 7),
        (47, 'Transporte', '47', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 7),
        (48, 'Passeios', '48', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 7),
        (49, 'Compras', '49', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 7),
        (50, 'Outros', '50', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 7);

COMMIT;
