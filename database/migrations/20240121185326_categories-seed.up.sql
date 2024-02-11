BEGIN;

INSERT INTO category_groups (name, icon, created_at, updated_at, deleted_at, version)
VALUES  ('Casa', 'Home', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0),
        ('Comidas e Bebidas', 'Apple', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0),
        ('Entretenimento', 'RollerCoaster', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0),
        ('Transporte', 'CarFront', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0),
        ('Vida', 'Trees', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0),
        ('Saúde', 'HeartPulse', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0),
        ('Viagem', 'Luggage', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0),
        ('Geral', 'ScrollText', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0);

INSERT INTO categories (name, icon, created_at, updated_at, deleted_at, version, category_group_id)
VALUES
        -- Casa
        ('Aluguel', 'KeyRound', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        ('Condomínio', 'Building', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        ('Água', 'Droplet', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        ('Lúz', 'Lightbulb', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        ('Gás', 'Flame', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        ('Internet', 'Wifi', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        ('Celular', 'Smartphone', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        ('Manutenção', 'Hammer', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        ('Animal de estimação', 'PawPrint', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        ('Eletrodomésticos', 'Refrigerator', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        ('Móveis', 'Sofa', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        ('Serviços', 'HandPlatter', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        ('Limpeza', 'WashingMachine', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        ('Outros', 'Home', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 1),
        -- Comidas e Bebidas
        ('Supermercado', 'ShoppingCart', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 2),
        ('Restaurantes', 'Utensils', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 2),
        ('Delivery', 'Bike', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 2),
        ('Bebidas', 'Martini', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 2),
        ('Outros', 'Apple', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 2),
        -- Entretenimento
        ('Cinema', 'Popcorn', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
        ('Shows', 'Theater', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
        ('Jogos', 'Gamepad2', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
        ('Livros', 'Book', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
        ('Música', 'Music4', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
        ('Streaming', 'Tv', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
        ('Festas', 'PartyPopper', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
        ('Outros', 'RollerCoaster', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 3),
        -- Transporte
        ('Combustível', 'Fuel', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
        ('Estacionamento', 'ParkingMeter', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
        ('Manutenção', 'Wrench', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
        ('Transporte público', 'BusFront', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
        ('Aplicativos de transporte', 'CarTaxiFront', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
        ('Pedágio', 'TrainFrontTunnel', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
        ('Aluguel de veículos', 'Car', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
        ('Seguro', 'Shield', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
        ('Multas', 'TicketX', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
        ('Outros', 'CarFront', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 4),
        -- Vida
        ('Academia/Esportes', 'Dumbbell', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 5),
        ('Cuidados pessoais', 'PersonStanding', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 5),
        ('Educação', 'GraduationCap', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 5),
        ('Presentes', 'Gift', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 5),
        ('Roupas', 'Shirt', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 5),
        ('Outros', 'Trees', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 5),
        -- Saúde
        ('Farmácia', 'Pill', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
        ('Dentista', 'Smile', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
        ('Psicólogo', 'Brain', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
        ('Plano de saúde', 'ShieldPlus', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
        ('Médico', 'Stethoscope', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
        ('Nutricionista', 'Carrot', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
        ('Exames', 'Microscope', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
        ('Procedimentos médicos', 'ClipboardPlus', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
        ('Outros', 'HeartPulse', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 6),
        -- Viagem
        ('Passagens', 'Plane', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 7),
        ('Hospedagem', 'Hotel', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 7),
        ('Alimentação', 'ConciergeBell', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 7),
        ('Transporte', 'Sailboat', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 7),
        ('Passeios', 'Compass', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 7),
        ('Compras', 'ShoppingBag', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 7),
        ('Outros', 'Luggage', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 7),
        -- Geral
        ('Impostos', 'BadgeDollarSign', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 8),
        ('Outros', 'ScrollText', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 0, 8);

COMMIT;
