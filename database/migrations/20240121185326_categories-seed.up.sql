BEGIN;

INSERT INTO category_groups (name, icon, created_at, updated_at, deleted_at, version)
VALUES  ('Casa', 'Home', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0),
        ('Comidas e Bebidas', 'Apple', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0),
        ('Entretenimento', 'RollerCoaster', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0),
        ('Transporte', 'CarFront', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0),
        ('Vida', 'Trees', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0),
        ('Saúde', 'HeartPulse', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0),
        ('Viagem', 'Plane', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0),
        ('Geral', 'ScrollText', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0);

INSERT INTO categories (name, icon, created_at, updated_at, deleted_at, version, category_group_id)
VALUES
        -- Casa
        ('Aluguel', 'Home', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        ('Condomínio', 'Home', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        ('Água', 'Home', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        ('Eletricidade', 'Home', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        ('Aquecimento/Gás', 'Home', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        ('Internet', 'Home', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        ('Manutenção', 'Home', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        ('Animal de estimação', 'Home', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        ('Eletrodomésticos', 'Home', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        ('Móveis', 'Home', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        ('Serviços', 'Home', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        ('Produtos de limpeza', 'Home', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        ('Outros', 'Home', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        -- Comidas e Bebidas
        ('Supermercado', 'Apple', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 2),
        ('Restaurantes', 'Apple', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 2),
        ('Delivery', 'Apple', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 2),
        ('Bebidas', 'Apple', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 2),
        ('Outros', 'Apple', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 2),
        -- Entretenimento
        ('Cinema', 'RollerCoaster', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
        ('Shows', 'RollerCoaster', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
        ('Jogos', 'RollerCoaster', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
        ('Livros', 'RollerCoaster', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
        ('Música', 'RollerCoaster', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
        ('Streaming', 'RollerCoaster', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
        ('Festas', 'RollerCoaster', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
        ('Outros', 'RollerCoaster', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
        -- Transporte
        ('Combustível', 'CarFront', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
        ('Estacionamento', 'CarFront', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
        ('Manutenção', 'CarFront', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
        ('Transporte público', 'CarFront', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
        ('Aplicativos de transporte', 'CarFront', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
        ('Pedágio', 'CarFront', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
        ('Aluguel de veículos', 'CarFront', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
        ('Seguro', 'CarFront', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
        ('Outros', 'CarFront', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
        -- Vida
        ('Academia/Esportes', 'Trees', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 5),
        ('Cuidados pessoais', 'Trees', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 5),
        ('Educação', 'Trees', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 5),
        ('Presentes', 'Trees', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 5),
        ('Roupas', 'Trees', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 5),
        ('Outros', 'Trees', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 5),
        -- Saúde
        ('Farmácia', 'HeartPulse', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
        ('Plano odontológico', 'HeartPulse', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
        ('Plano de saúde', 'HeartPulse', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
        ('Consulta médica', 'HeartPulse', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
        ('Exames', 'HeartPulse', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
        ('Procedimentos médicos', 'HeartPulse', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
        ('Outros', 'HeartPulse', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
        -- Viagem
        ('Passagens', 'Plane', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 7),
        ('Hospedagem', 'Plane', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 7),
        ('Alimentação', 'Plane', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 7),
        ('Transporte', 'Plane', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 7),
        ('Passeios', 'Plane', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 7),
        ('Compras', 'Plane', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 7),
        ('Outros', 'Plane', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 7),
        -- Geral
        ('Impostos', 'ScrollText', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 8),
        ('Outros', 'ScrollText', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 8);

COMMIT;
