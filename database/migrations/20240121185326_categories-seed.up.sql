BEGIN;

INSERT INTO category_groups (name, icon, created_at, updated_at, deleted_at, version)
VALUES  ('Casa', '1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0),
        ('Comidas e Bebidas', '2', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0),
        ('Entretenimento', '3', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0),
        ('Transporte', '4', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0),
        ('Vida', '5', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0),
        ('Saúde', '6', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0),
        ('Geral', '7', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0),
        ('Viagem', '8', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0);

INSERT INTO categories (name, icon, created_at, updated_at, deleted_at, version, category_group_id)
VALUES
-- Casa
('Aluguel', '1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
('Condomínio', '2', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
('Água', '3', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
('Eletricidade', '4', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
('Aquecimento/Gás', '4', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
('Internet', '5', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
('Manutenção', '6', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
('Animal de estimação', '7', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
('Eletrodomésticos', '8', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
('Móveis', '9', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
('Serviços', '10', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
('Produtos de limpeza', '11', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
('Outros', '12', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
-- Comidas e Bebidas
('Supermercado', '13', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 2),
('Restaurantes', '14', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 2),
('Delivery', '15', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 2),
('Bebidas', '16', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 2),
('Outros', '17', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 2),
-- Entretenimento
('Cinema', '18', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
('Shows', '19', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
('Jogos', '20', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
('Livros', '21', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
('Música', '22', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
('Streaming', '23', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
('Festas', '24', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
('Outros', '25', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
-- Transporte
('Combustível', '26', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
('Estacionamento', '27', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
('Manutenção', '28', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
('Transporte público', '29', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
('Aplicativos de transporte', '30', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
('Pedágio', '31', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
('Aluguel de veículos', '32', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
('Seguro', '33', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
('Outros', '34', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
-- Vida
('Academia/Esportes', '35', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 5),
('Cuidados pessoais', '36', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 5),
('Educação', '36', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 5),
('Presentes', '41', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 5),
('Roupas', '41', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 5),
('Outros', '43', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 5),
-- Saúde
('Farmácia', '38', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
('Plano odontológico', '39', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
('Plano de saúde', '40', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
('Consulta médica', '41', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
('Exames', '41', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
('Procedimentos médicos', '41', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
('Outros', '45', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
-- Geral
('Impostos', '44', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 7),
('Outros', '45', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 7),
-- Viagem
('Passagens', '46', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 8),
('Hospedagem', '47', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 8),
('Alimentação', '48', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 8),
('Transporte', '49', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 8),
('Passeios', '50', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 8),
('Compras', '51', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 8),
('Outros', '52', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 8);

COMMIT;
